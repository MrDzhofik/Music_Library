package handlers

import (
	"context"
	"effectiveMobile/internal/usecase"
	"effectiveMobile/models"
	"encoding/json"
	"net/http"
)

type SongHandler struct {
	songUsecase usecase.SongUsecase
}

type SongIDRequest struct {
	ID int `json:"id"`
}

type AddSongRequest struct {
	Group string `json:"group"`
	Song  string `json:"song"`
}

type SongRequest struct {
	ID      int         `json:"id"`
	NewSong models.Song `json:"song"`
}

func NewSongHandler(songUsecase usecase.SongUsecase) *SongHandler {
	return &SongHandler{songUsecase: songUsecase}
}

// Просмотр всех песен
func (h *SongHandler) GetAllSongs(w http.ResponseWriter, r *http.Request) {
	ctx := context.Background()
	songs, err := h.songUsecase.GetAllSongs(ctx)

	if err != nil {
		http.Error(w, "Неправильный запрос", http.StatusBadRequest)
		return
	}

	if len(songs) == 0 {
		http.Error(w, "Библиотека пуста!", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(songs)
}

// Просмотр одной песни по id
func (h *SongHandler) GetSongByID(w http.ResponseWriter, r *http.Request) {
	url := r.URL
	str := url.String()

	id := str[len(str)-1] - 48

	defer r.Body.Close()

	song, err := h.songUsecase.GetSongByID(r.Context(), int(id))

	if err != nil {
		http.Error(w, "Неправильный запрос", http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(song)
}

// Добавление песни запросом
// {
//  "group": "Muse",
//  "song": "Supermassive Black Hole"
// }

func (h *SongHandler) AddSong(w http.ResponseWriter, r *http.Request) {
	var song AddSongRequest

	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		http.Error(w, "Неправильный запрос", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	id, err := h.songUsecase.AddGroup(r.Context(), models.Group{
		Name: &song.Group,
	})

	if err != nil {
		http.Error(w, "Ошибка добавления группы", http.StatusInternalServerError)
		return
	}

	err = h.songUsecase.AddSong(r.Context(), models.Song{
		Group: &id,
		Name:  &song.Song,
	})

	if err != nil {
		http.Error(w, "Ошибка добавления песни", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Song added successfully"})
}

// Обновление песни через id и названием песни
func (h *SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	var newSong SongRequest

	err := json.NewDecoder(r.Body).Decode(&newSong)
	if err != nil {
		http.Error(w, "Неправильный запрос", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	err = h.songUsecase.UpdateSong(r.Context(), newSong.ID, newSong.NewSong)

	if err != nil {
		http.Error(w, "Ошибка изменения песни", http.StatusInternalServerError)
		return
	}
}

// Удаление песни по id
func (h *SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	var songID SongIDRequest

	err := json.NewDecoder(r.Body).Decode(&songID)
	if err != nil {
		http.Error(w, "Неправильный запрос", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	err = h.songUsecase.DeleteSong(r.Context(), songID.ID)

	if err != nil {
		http.Error(w, "Ошибка удаления песни", http.StatusInternalServerError)
		return
	}
}
