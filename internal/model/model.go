package model

type User struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

type Movie struct {
	Movie_id     int     `json:"movie_id"`
	Kinopoisk_id int     `json:"kinopoisk_id"`
	Year         int     `json:"year"`
	Name         string  `json:"name"`
	Genre        string  `json:"genre"`
	Rating_kp    float64 `json:"rating_kp"`
	Rating_imdb  float64 `json:"rating_imdb"`
	Rating_avg   float64 `json:"rating_avg"`
}

type Playlist struct {
}
