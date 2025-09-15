package repos

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Darari17/be-tickitz/internal/models"
	"github.com/Darari17/be-tickitz/internal/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
)

type MovieRepo struct {
	db    *pgxpool.Pool
	redis *redis.Client
}

func NewMovieRepo(db *pgxpool.Pool, redis *redis.Client) *MovieRepo {
	return &MovieRepo{db: db, redis: redis}
}

func (mr *MovieRepo) GetUpcomingMovies(ctx context.Context, page int) ([]models.Movie, error) {
	const pageSize = 4
	offset := (page - 1) * pageSize

	redisKey := fmt.Sprintf("movies:upcoming:page:%d", page)
	var cached []models.Movie

	ok, err := utils.GetCacheRedis(ctx, mr.redis, redisKey, &cached)
	if err != nil {
		fmt.Printf("redis error: %v\n", err)
	} else if ok {
		return cached, nil
	}

	sql := `
		SELECT m.id, m.backdrop_path, m.overview, m.popularity, m.poster_path,
		       m.release_date, m.duration, m.title, m.director_name,
		       m.created_at, m.updated_at, m.deleted_at,
		       COALESCE(JSON_AGG(DISTINCT jsonb_build_object('id', g.id, 'name', g.name))
		                FILTER (WHERE g.id IS NOT NULL), '[]') AS genres,
		       COALESCE(JSON_AGG(DISTINCT jsonb_build_object('id', c.id, 'name', c.name))
		                FILTER (WHERE c.id IS NOT NULL), '[]') AS casts
		FROM movies m
		LEFT JOIN movies_genres mg ON m.id = mg.movies_id
		LEFT JOIN genres g ON g.id = mg.genres_id
		LEFT JOIN movies_casts mc ON m.id = mc.movies_id
		LEFT JOIN casts c ON c.id = mc.casts_id
		WHERE m.release_date > NOW()
		GROUP BY m.id
		ORDER BY m.release_date ASC
		LIMIT $1 OFFSET $2
	`

	rows, err := mr.db.Query(ctx, sql, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		var genresJSON, castsJSON []byte

		if err := rows.Scan(
			&m.ID, &m.Backdrop, &m.Overview, &m.Popularity, &m.Poster,
			&m.ReleaseDate, &m.Duration, &m.Title, &m.Director,
			&m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
			&genresJSON, &castsJSON,
		); err != nil {
			return nil, err
		}

		_ = json.Unmarshal(genresJSON, &m.Genres)
		_ = json.Unmarshal(castsJSON, &m.Casts)

		movies = append(movies, m)
	}

	if err := utils.SetCacheRedis(ctx, mr.redis, redisKey, movies, 5*time.Minute); err != nil {
		fmt.Printf("failed to set redis cache: %v\n", err)
	}

	return movies, nil
}

func (mr *MovieRepo) GetPopularMovies(ctx context.Context, page int) ([]models.Movie, error) {
	const pageSize = 4
	offset := (page - 1) * pageSize

	redisKey := fmt.Sprintf("movies:popular:page:%d", page)
	var cached []models.Movie
	ok, err := utils.GetCacheRedis(ctx, mr.redis, redisKey, &cached)
	if err != nil {
		fmt.Printf("redis error: %v\n", err)
	} else if ok {
		return cached, nil
	}

	sql := `
		SELECT m.id, m.backdrop_path, m.overview, m.popularity, m.poster_path,
		       m.release_date, m.duration, m.title, m.director_name,
		       m.created_at, m.updated_at, m.deleted_at,
		       COALESCE(JSON_AGG(DISTINCT jsonb_build_object('id', g.id, 'name', g.name))
		                FILTER (WHERE g.id IS NOT NULL), '[]') AS genres,
		       COALESCE(JSON_AGG(DISTINCT jsonb_build_object('id', c.id, 'name', c.name))
		                FILTER (WHERE c.id IS NOT NULL), '[]') AS casts
		FROM movies m
		LEFT JOIN movies_genres mg ON m.id = mg.movies_id
		LEFT JOIN genres g ON g.id = mg.genres_id
		LEFT JOIN movies_casts mc ON m.id = mc.movies_id
		LEFT JOIN casts c ON c.id = mc.casts_id
		GROUP BY m.id
		ORDER BY m.popularity DESC
		LIMIT $1 OFFSET $2
	`

	rows, err := mr.db.Query(ctx, sql, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		var genresJSON, castsJSON []byte

		if err := rows.Scan(
			&m.ID, &m.Backdrop, &m.Overview, &m.Popularity, &m.Poster,
			&m.ReleaseDate, &m.Duration, &m.Title, &m.Director,
			&m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
			&genresJSON, &castsJSON,
		); err != nil {
			return nil, err
		}

		_ = json.Unmarshal(genresJSON, &m.Genres)
		_ = json.Unmarshal(castsJSON, &m.Casts)

		movies = append(movies, m)
	}

	if err := utils.SetCacheRedis(ctx, mr.redis, redisKey, movies, 5*time.Minute); err != nil {
		fmt.Printf("failed to set redis cache: %v\n", err)
	}

	return movies, nil
}

func (mr *MovieRepo) GetAllMovies(ctx context.Context, page int, search string, genreName string) ([]models.Movie, error) {
	const pageSize = 12
	offset := (page - 1) * pageSize

	redisKey := fmt.Sprintf("movies:all:page:%d:search:%s:genre:%s", page, search, genreName)
	var cached []models.Movie
	ok, err := utils.GetCacheRedis(ctx, mr.redis, redisKey, &cached)
	if err != nil {
		fmt.Printf("redis error: %v\n", err)
	} else if ok {
		return cached, nil
	}

	sql := `
		WITH filtered AS (
			SELECT m.id, m.backdrop_path, m.overview, m.popularity, m.poster_path,
			       m.release_date, m.duration, m.title, m.director_name,
			       m.created_at, m.updated_at, m.deleted_at,
			       COALESCE(JSON_AGG(DISTINCT jsonb_build_object('id', g.id, 'name', g.name))
			                FILTER (WHERE g.id IS NOT NULL), '[]') AS genres,
			       COALESCE(JSON_AGG(DISTINCT jsonb_build_object('id', c.id, 'name', c.name))
			                FILTER (WHERE c.id IS NOT NULL), '[]') AS casts
			FROM movies m
			LEFT JOIN movies_genres mg ON m.id = mg.movies_id
			LEFT JOIN genres g ON g.id = mg.genres_id
			LEFT JOIN movies_casts mc ON m.id = mc.movies_id
			LEFT JOIN casts c ON c.id = mc.casts_id
			WHERE ($1 = '' OR LOWER(m.title) LIKE LOWER('%' || $1 || '%'))
			  AND ($2 = '' OR LOWER(g.name) = LOWER($2))
			GROUP BY m.id
		)
		SELECT *
		FROM filtered
		ORDER BY release_date DESC
		LIMIT $3 OFFSET $4;
	`

	rows, err := mr.db.Query(ctx, sql, search, genreName, pageSize, offset)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []models.Movie
	for rows.Next() {
		var m models.Movie
		var genresJSON, castsJSON []byte

		if err := rows.Scan(
			&m.ID, &m.Backdrop, &m.Overview, &m.Popularity, &m.Poster,
			&m.ReleaseDate, &m.Duration, &m.Title, &m.Director,
			&m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
			&genresJSON, &castsJSON,
		); err != nil {
			return nil, err
		}

		_ = json.Unmarshal(genresJSON, &m.Genres)
		_ = json.Unmarshal(castsJSON, &m.Casts)

		movies = append(movies, m)
	}

	if err := utils.SetCacheRedis(ctx, mr.redis, redisKey, movies, 5*time.Minute); err != nil {
		fmt.Printf("failed to set redis cache: %v\n", err)
	}
	return movies, nil
}

func (mr *MovieRepo) GetMovieDetail(ctx context.Context, id int) (*models.Movie, error) {
	redisKey := fmt.Sprintf("movies:detail:%d", id)
	var cached models.Movie
	ok, err := utils.GetCacheRedis(ctx, mr.redis, redisKey, &cached)
	if err != nil {
		fmt.Printf("redis error: %v\n", err)
	} else if ok {
		return &cached, nil
	}

	sql := `
		SELECT m.id, m.backdrop_path, m.overview, m.popularity, m.poster_path,
		       m.release_date, m.duration, m.title, m.director_name,
		       m.created_at, m.updated_at, m.deleted_at,
		       COALESCE(
		           JSON_AGG(DISTINCT jsonb_build_object('id', g.id, 'name', g.name))
		           FILTER (WHERE g.id IS NOT NULL), '[]'
		       ) AS genres,
		       COALESCE(
		           JSON_AGG(DISTINCT jsonb_build_object('id', c.id, 'name', c.name))
		           FILTER (WHERE c.id IS NOT NULL), '[]'
		       ) AS casts
		FROM movies m
		LEFT JOIN movies_genres mg ON m.id = mg.movies_id
		LEFT JOIN genres g ON g.id = mg.genres_id
		LEFT JOIN movies_casts mc ON m.id = mc.movies_id
		LEFT JOIN casts c ON c.id = mc.casts_id
		WHERE m.id = $1
		GROUP BY m.id
	`

	var m models.Movie
	var genresJSON, castsJSON []byte

	err = mr.db.QueryRow(ctx, sql, id).Scan(
		&m.ID, &m.Backdrop, &m.Overview, &m.Popularity, &m.Poster,
		&m.ReleaseDate, &m.Duration, &m.Title, &m.Director,
		&m.CreatedAt, &m.UpdatedAt, &m.DeletedAt,
		&genresJSON, &castsJSON,
	)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(genresJSON, &m.Genres); err != nil {
		return nil, err
	}
	if err := json.Unmarshal(castsJSON, &m.Casts); err != nil {
		return nil, err
	}

	if err := utils.SetCacheRedis(ctx, mr.redis, redisKey, m, 5*time.Minute); err != nil {
		fmt.Printf("failed to set redis cache: %v\n", err)
	}
	return &m, nil
}
