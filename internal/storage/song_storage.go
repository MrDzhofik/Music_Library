package storage

import (
	"context"
	"effectiveMobile/models"
	"log"

	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
)

type SongStorage interface {
	GetAllSongs(ctx context.Context, filterName, ilterGroup string) ([]models.Song, error)
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

func (s *songStorage) GetAllSongs(ctx context.Context, filterName, filterGroup string) ([]models.Song, error) {
	s.infoLog.Print("Запускаем SQL запрос по получению всех песен")
	var songs []models.Song
	var err error

	query := `SELECT song_name name, group_name, release_date, link 
					FROM songs s
					INNER JOIN groups g ON g.id = s.group_id
					WHERE s.song_name LIKE '%' || $1  || '%' AND group_name LIKE '%' || $2  || '%'`

	err = pgxscan.Select(context.Background(), s.db, &songs, query, filterName, filterGroup)
	if err != nil {
		s.errorLog.Println(err)
	}

	return songs, err
}

func (s *songStorage) GetSongByID(ctx context.Context, id int) (models.Song, error) {
	s.infoLog.Print("Запускаем SQL запрос по получению песни по ID")
	query := `SELECT song_name name, group_name, release_date, text, link 
	FROM songs s
	INNER JOIN groups g ON g.id = s.group_id
	WHERE s.id = $1`

	var song []models.Song
	err := pgxscan.Select(context.Background(), s.db, &song, query, id)
	if err != nil {
		s.errorLog.Println(err)
	}

	return song[0], err
}

func (s *songStorage) DeleteSong(ctx context.Context, id int) error {
	s.infoLog.Print("Запускаем SQL запрос по удалению песни по ID")
	query := `DELETE FROM songs WHERE id = $1`

	_, err := s.db.Exec(context.Background(), query, id)
	if err != nil {
		s.errorLog.Println(err)
	}

	return err
}

func (s *songStorage) AddSong(ctx context.Context, song models.Song) error {
	s.infoLog.Print("Запускаем SQL запрос по добавлению песни")
	query := "INSERT INTO songs (group_id, song_name) VALUES ($1, $2)"
	_, err := s.db.Exec(
		context.Background(),
		query,
		song.Group, song.Name,
	)
	if err != nil {
		s.errorLog.Println(err)
	}
	return err
}

func (s *songStorage) AddGroup(ctx context.Context, group models.Group) (int, error) {
	var groupID int

	s.infoLog.Print("Запускаем SQL запрос по добавлению группы")

	// Сначала проверим, существует ли уже такая группа (артист)
	query := `SELECT id FROM groups WHERE group_name = $1`
	err := s.db.QueryRow(ctx, query, group.Name).Scan(&groupID)

	if err == nil {
		// Артист уже существует, возвращаем его ID
		return groupID, nil
	}

	if err != pgx.ErrNoRows {
		s.errorLog.Println(err)
		// Произошла ошибка, отличная от "запись не найдена"
		return 0, err
	}

	query = "INSERT INTO groups(group_name) VALUES ($1) RETURNING id"
	err = s.db.QueryRow(
		context.Background(),
		query,
		group.Name,
	).Scan(&groupID)
	if err != nil {
		s.errorLog.Println(err)
	}
	return groupID, err
}

func (s *songStorage) UpdateSong(ctx context.Context, id int, newSong models.Song) error {
	s.infoLog.Print("Запускаем SQL запрос по обновлению песни по ID")

	query := "UPDATE songs SET song_name = $1, text = $2, link = $3 WHERE id = $4"
	_, err := s.db.Exec(
		context.Background(),
		query,
		newSong.Name, newSong.Text, newSong.Link, id,
	)
	if err != nil {
		s.errorLog.Println(err)
	}
	return err
}
