package usecase

import (
	"context"
	"effectiveMobile/internal/storage"
	"effectiveMobile/models"
	"log"
)

type SongUsecase interface {
	GetAllSongs(ctx context.Context, filter string) ([]models.Song, error)
	GetSongByID(ctx context.Context, id int) (models.Song, error)
	AddSong(ctx context.Context, song models.Song) error
	UpdateSong(ctx context.Context, id int, song models.Song) error
	DeleteSong(ctx context.Context, id int) error
	AddGroup(ctx context.Context, group models.Group) (int, error)
}

type songUsecase struct {
	songStorage storage.SongStorage
	infoLog     *log.Logger
	errorLog    *log.Logger
}

func NewSongUsecase(s storage.SongStorage, infoLog, errorLog *log.Logger) SongUsecase {
	return &songUsecase{
		songStorage: s,
		infoLog:     infoLog,
		errorLog:    errorLog,
	}
}

func (uc *songUsecase) GetAllSongs(ctx context.Context, filter string) ([]models.Song, error) {
	return uc.songStorage.GetAllSongs(ctx, filter)
}

func (uc *songUsecase) GetSongByID(ctx context.Context, id int) (models.Song, error) {
	return uc.songStorage.GetSongByID(ctx, id)
}

func (uc *songUsecase) AddSong(ctx context.Context, song models.Song) error {
	return uc.songStorage.AddSong(ctx, song)
}

func (uc *songUsecase) UpdateSong(ctx context.Context, id int, newSong models.Song) error {
	return uc.songStorage.UpdateSong(ctx, id, newSong)
}

func (uc *songUsecase) DeleteSong(ctx context.Context, id int) error {
	return uc.songStorage.DeleteSong(ctx, id)
}

func (uc *songUsecase) AddGroup(ctx context.Context, group models.Group) (int, error) {
	return uc.songStorage.AddGroup(ctx, group)
}
