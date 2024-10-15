package model

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Movie struct {
	Movie_id     int     `json:"movie_id"`
	Name         string  `json:"name"`
	Rating_kp    float64 `json:"rating_kp"`
	Rating_imdb  float64 `json:"rating_imdb"`
	Kinopoisk_id int     `json:"kinopoisk_id"`
	Rating_avg   float64 `json:"rating_avg"`
	Preview      string  `json:"preview"`
	Trailer      string  `json:"trailer"`
	Genre        string  `json:"genre"`
	Year         int     `json:"year"`
	Type         string  `json:"type"`
}

type Playlist struct {
}
