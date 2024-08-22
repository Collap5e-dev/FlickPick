package repo

import (
	"context"
	"log"

	"github.com/Collap5e-dev/FlickPick/internal/handler"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewMovieList (db *sqlx.DB) *MovieRepo {
	return &MovieRepo{
		db: db,
	}
}

type MovieRepo struct {
	db *sqlx.DB
}

func (r *MovieRepo) GetMovieList(ctx context.Context) ([]handler.MovieStruct, error) {
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
		log.Fatalf("Ошибка при выполнении запроса: %v\n", err)
	}
	defer movieTable.Close()
	movieList := make([]handler.MovieStruct, 0)
	for movieTable.Next() {
		var movie_id, kinopoisk_id, year int
		var name, genre string
		var rating_kp, rating_imdb, rating_avg float64
		err := movieTable.Scan(&movie_id, &name, &rating_kp, &rating_imdb, &kinopoisk_id, &rating_avg, &genre, &year)
		if err != nil {
			log.Fatalf("Ошибка при сканировании строки: %v\n", err)
			continue
		}
		mStruct := handler.MovieStruct{
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
