CREATE TABLE content
(
    movie_id     SERIAL PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    rating_kp    NUMERIC(3, 1),
    rating_imdb  NUMERIC(3, 1),
    kinopoisk_id VARCHAR(50),
    rating_avg   NUMERIC(3, 1),
    preview      TEXT,
    trailer      TEXT,
    genre        VARCHAR(100),
    year         INT,
    type         VARCHAR(50)
);