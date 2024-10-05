package handlers

import (
	"context"
	"effectiveMobile/internal/usecase"
	"effectiveMobile/models"
	"encoding/json"
	"log"
	"net/http"
)

type SongHandler struct {
	songUsecase usecase.SongUsecase
	infoLog     *log.Logger
	errorLog    *log.Logger
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

type BadRequest struct {
	Message string `json:"message"`
}

func NewSongHandler(songUsecase usecase.SongUsecase, infoLog, errorLog *log.Logger) *SongHandler {
	return &SongHandler{
		songUsecase: songUsecase,
		infoLog:     infoLog,
		errorLog:    errorLog}
}

// Get all songs godoc
// @Summary      List songs
// @Description  get songs
// @Tags         song
// @Produce      json
// @Param        name   query      string  true  "Song name"
// @Param        group   query      string  true  "Group name"
// @Success      200  {object}  []models.Song
// @Failure      400  {object}  BadRequest
// @Router       /api/songs [get]
func (h *SongHandler) GetAllSongs(w http.ResponseWriter, r *http.Request) {
	h.infoLog.Println("Получаем все песни")
	filterName := r.URL.Query().Get("name")
	filterGroup := r.URL.Query().Get("group")
	ctx := context.Background()
	songs, err := h.songUsecase.GetAllSongs(ctx, filterName, filterGroup)

	if err != nil {
		h.errorLog.Printf("Неправильный запрос: %v", err)
		http.Error(w, "Неправильный запрос", http.StatusBadRequest)
		return
	}

	if len(songs) == 0 {
		http.Error(w, "Нет подходящих песен!", http.StatusOK)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(songs)
}

// Get one song by ID godoc
// @Summary      Give Song with certain ID
// @Description  get song by ID
// @Tags         song
// @Produce      json
// @Param        id   path      int  true  "Song ID"
// @Success      200  {object}  models.Song
// @Failure      400  {object}  BadRequest
// @Router       /api/song/{id} [get]
func (h *SongHandler) GetSongByID(w http.ResponseWriter, r *http.Request) {
	h.infoLog.Println("Получаем песню по ID")
	url := r.URL
	str := url.String()

	id := str[len(str)-1] - 48

	defer r.Body.Close()

	song, err := h.songUsecase.GetSongByID(r.Context(), int(id))

	if err != nil {
		h.errorLog.Printf("Неправильный запрос: %v", err)
		http.Error(w, "Неправильный запрос", http.StatusBadRequest)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(song)
}

// Добавление песни запросом
// {
//  "group": "Muse",
//  "song": "Supermassive Black Hole"
// }

// Add new song
// @Summary      Add new song
// @Description  Add song to library
// @Tags         song
// @Accept       json
// @Produce      json
// @Success      200  {string}  message
// @Failure      400  {string}  http.BadRequest
// @Failure      500  {string}  http.InternalServerError
// @Router       /api/song/create [post]
func (h *SongHandler) AddSong(w http.ResponseWriter, r *http.Request) {
	h.infoLog.Println("Добавляем песню")
	var song AddSongRequest

	err := json.NewDecoder(r.Body).Decode(&song)
	if err != nil {
		h.errorLog.Printf("Неправильный запрос: %v", err)
		http.Error(w, "Неправильный запрос", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	id, err := h.songUsecase.AddGroup(r.Context(), models.Group{
		Name: &song.Group,
	})

	if err != nil {
		h.errorLog.Printf("Ошибка добавления группы: %v", err)
		http.Error(w, "Ошибка добавления группы", http.StatusInternalServerError)
		return
	}

	err = h.songUsecase.AddSong(r.Context(), models.Song{
		Group: &id,
		Name:  &song.Song,
	})

	if err != nil {
		h.errorLog.Printf("Ошибка добавления песни: %v", err)
		http.Error(w, "Ошибка добавления песни", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Песня добавлена успешно"})
}

// Update new song
// @Summary      Update song
// @Description  Update song in library by ID
// @Tags         song
// @Accept       json
// @Produce      json
// @Success      200  {string}  message
// @Failure      400  {string}  http.BadRequest
// @Failure      500  {string}  http.InternalServerError
// @Router       /api/song/update [put]
func (h *SongHandler) UpdateSong(w http.ResponseWriter, r *http.Request) {
	h.infoLog.Println("Обновляем песню по ID")
	var newSong SongRequest

	err := json.NewDecoder(r.Body).Decode(&newSong)
	if err != nil {
		h.errorLog.Printf("Неправильный запрос: %v", err)
		http.Error(w, "Неправильный запрос", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	err = h.songUsecase.UpdateSong(r.Context(), newSong.ID, newSong.NewSong)

	if err != nil {
		h.errorLog.Printf("Ошибка изменения песни: %v", err)
		http.Error(w, "Ошибка изменения песни", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Песня успешно изменена"})
}

// @Summary      Delete song
// @Description  Delete song from library by ID
// @Tags         song
// @Accept       json
// @Produce      json
// @Success      200  {string}  message
// @Failure      400  {string}  http.BadRequest
// @Failure      500  {string}  http.InternalServerError
// @Router       /api/song/delete [delete]
func (h *SongHandler) DeleteSong(w http.ResponseWriter, r *http.Request) {
	h.infoLog.Println("Удаляем песню по ID")
	var songID SongIDRequest

	err := json.NewDecoder(r.Body).Decode(&songID)
	if err != nil {
		h.errorLog.Printf("Неправильный запрос: %v", err)
		http.Error(w, "Неправильный запрос", http.StatusBadRequest)
		return
	}

	defer r.Body.Close()

	err = h.songUsecase.DeleteSong(r.Context(), songID.ID)

	if err != nil {
		h.errorLog.Printf("Ошибка удаления песни: %v", err)
		http.Error(w, "Ошибка удаления песни", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Песня успешно удалена"})
}
