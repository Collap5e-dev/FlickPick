package repo

import (
	"context"
	"fmt"
	"log"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"

	"github.com/Collap5e-dev/FlickPick/internal/model"
)

func NewMovieList(db *sqlx.DB) *MovieRepo {
	return &MovieRepo{
		db: db,
	}
}

type MovieRepo struct {
	db *sqlx.DB
}

func (r *MovieRepo) CreateUser(ctx context.Context, user model.User) error {
	// TODO: implement me
	return nil
}

func (r *MovieRepo) GetMovieList(ctx context.Context) ([]model.Movie, error) {
	movieTable, err := r.db.Query(`
		SELECT
			movie_id,
			name,
			rating_kp,
			rating_imdb,
			kinopoisk_id,
			rating_avg,
			genre,
			year
		FROM
			content
		ORDER BY
			name
	`)
	if err != nil {
		return nil, fmt.Errorf("ошибка при выполнении запроса: %w", err)
	}
	defer movieTable.Close()
	movieList := make([]model.Movie, 0)
	for movieTable.Next() {
		var movie_id, kinopoisk_id, year int
		var name, genre string
		var rating_kp, rating_imdb, rating_avg float64
		err := movieTable.Scan(&movie_id, &name, &rating_kp, &rating_imdb, &kinopoisk_id, &rating_avg, &genre, &year)
		if err != nil {
			log.Fatalf("Ошибка при сканировании строки: %v\n", err)
			continue
		}
		mStruct := model.Movie{
			Movie_id:     movie_id,
			Name:         name,
			Rating_kp:    rating_kp,
			Rating_imdb:  rating_imdb,
			Kinopoisk_id: kinopoisk_id,
			Rating_avg:   rating_avg,
			Genre:        genre,
			Year:         year}
		movieList = append(movieList, mStruct)
	}
	return movieList, nil
}
