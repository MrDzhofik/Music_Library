package storage

import (
	"context"
	"fmt"

	"effectiveMobile/models"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type SongStorage interface {
	GetAllSongs(ctx context.Context) ([]models.Song, error)
	GetSongByID(ctx context.Context, id int) (models.Song, error)
	AddSong(ctx context.Context, song models.Song) error
	UpdateSong(ctx context.Context, id int, song models.Song) error
	DeleteSong(ctx context.Context, id int) error
	AddGroup(ctx context.Context, group models.Group) (int, error)
}

type songStorage struct {
	db *pgx.Conn
}

func NewSongStorage(db *pgx.Conn) SongStorage {
	return &songStorage{db: db}
}

func (s *songStorage) GetAllSongs(ctx context.Context) ([]models.Song, error) {
	query := `SELECT song_name name, group_name, release_date, link 
	FROM songs s
	INNER JOIN groups g ON g.id = s.group_id`

	var songs []models.Song
	err := pgxscan.Select(context.Background(), s.db, &songs, query)

	return songs, err
}

func (s *songStorage) GetSongByID(ctx context.Context, id int) (models.Song, error) {
	query := `SELECT song_name name, group_name, release_date, text, link 
	FROM songs s
	INNER JOIN groups g ON g.id = s.group_id
	WHERE s.id = $1`

	var song []models.Song
	err := pgxscan.Select(context.Background(), s.db, &song, query, id)

	fmt.Println(err)
	return song[0], err
}

func (s *songStorage) DeleteSong(ctx context.Context, id int) error {
	query := `DELETE FROM songs WHERE id = $1`

	_, err := s.db.Exec(context.Background(), query, id)

	return err
}

func (s *songStorage) AddSong(ctx context.Context, song models.Song) error {
	query := "INSERT INTO Song (artist_id, song_name, release_date, lyrics, link) VALUES ($1, $2, $3, $4, $5)"
	_, err := s.db.Exec(
		context.Background(),
		query,
		song.Group, song.Name, song.ReleaseDate, song.Text, song.Link,
	)
	return err
}

func (s *songStorage) AddGroup(ctx context.Context, group models.Group) (int, error) {
	var groupID int

	// Сначала проверим, существует ли уже такая группа (артист)
	query := `SELECT id FROM Artist WHERE artist_name = $1`
	err := s.db.QueryRow(ctx, query, group.Name).Scan(&groupID)

	if err == nil {
		// Артист уже существует, возвращаем его ID
		return groupID, nil
	}

	if err != pgx.ErrNoRows {
		// Произошла ошибка, отличная от "запись не найдена"
		return 0, err
	}

	query = "INSERT INTO groups(group_name) VALUES ($1) RETURNING id"
	err = s.db.QueryRow(
		context.Background(),
		query,
		group.Name,
	).Scan(&groupID)
	return groupID, err
}

func (s *songStorage) UpdateSong(ctx context.Context, id int, newSong models.Song) error {
	query := "UPDATE Song SET artist_id = $1, song_name = $2, release_date = $3, lyrics = $4, link = $5 WHERE id = $6"
	_, err := s.db.Exec(
		context.Background(),
		query,
		newSong.Group, newSong.Name, newSong.ReleaseDate, newSong.Text, newSong.Link, id,
	)
	return err
}
