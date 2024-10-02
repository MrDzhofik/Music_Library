package storage

import (
	"context"
	"effectiveMobile/models"
	"log"

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
	db       *pgx.Conn
	infoLog  *log.Logger
	errorLog *log.Logger
}

func NewSongStorage(db *pgx.Conn, infoLog, errorLog *log.Logger) SongStorage {
	return &songStorage{
		db:       db,
		infoLog:  infoLog,
		errorLog: errorLog,
	}
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

	return song[0], err
}

func (s *songStorage) DeleteSong(ctx context.Context, id int) error {
	query := `DELETE FROM songs WHERE id = $1`

	_, err := s.db.Exec(context.Background(), query, id)

	return err
}

func (s *songStorage) AddSong(ctx context.Context, song models.Song) error {
	query := "INSERT INTO songs (group_id, song_name) VALUES ($1, $2)"
	_, err := s.db.Exec(
		context.Background(),
		query,
		song.Group, song.Name,
	)
	return err
}

func (s *songStorage) AddGroup(ctx context.Context, group models.Group) (int, error) {
	var groupID int

	// Сначала проверим, существует ли уже такая группа (артист)
	query := `SELECT id FROM groups WHERE group_name = $1`
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
	query := "UPDATE songs SET song_name = $1, text = $2, link = $3 WHERE id = $4"
	_, err := s.db.Exec(
		context.Background(),
		query,
		newSong.Name, newSong.Text, newSong.Link, id,
	)
	return err
}
