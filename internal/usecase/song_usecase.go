package usecase

import (
	"context"
	"effectiveMobile/internal/storage"
	"effectiveMobile/models"
)

type SongUsecase interface {
	GetAllSongs(ctx context.Context) ([]models.Song, error)
	GetSongByID(ctx context.Context, id int) (models.Song, error)
	AddSong(ctx context.Context, song models.Song) error
	UpdateSong(ctx context.Context, id int, song models.Song) error
	DeleteSong(ctx context.Context, id int) error
	AddGroup(ctx context.Context, group models.Group) (int, error)
}

type songUsecase struct {
	songStorage storage.SongStorage
}

func NewSongUsecase(s storage.SongStorage) SongUsecase {
	return &songUsecase{songStorage: s}
}

func (uc *songUsecase) GetAllSongs(ctx context.Context) ([]models.Song, error) {
	return uc.songStorage.GetAllSongs(ctx)
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
