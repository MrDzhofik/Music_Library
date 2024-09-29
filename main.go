package main

import (
	"context"
	"log"
	"net/http"

	"effectiveMobile/internal/handlers"
	"effectiveMobile/internal/storage"
	"effectiveMobile/internal/usecase"

	"github.com/gorilla/mux"
)

func main() {
	conn, err := storage.GetPostgres()
	if err != nil {
		log.Panic("Ошибка соединения с базой данных")
	}

	defer conn.Close(context.Background())

	// Инициализация слоёв
	songStorage := storage.NewSongStorage(conn)
	songUsecase := usecase.NewSongUsecase(songStorage)
	songHandler := handlers.NewSongHandler(songUsecase)

	// Настройка роутера
	router := mux.NewRouter()
	router.HandleFunc("/api/songs", songHandler.GetAllSongs).Methods("GET")
	router.HandleFunc("/api/song/{id:[0-9]+}", songHandler.GetSongByID).Methods("GET")
	router.HandleFunc("/api/song/add", songHandler.AddSong).Methods("POST")
	router.HandleFunc("/api/song/delete", songHandler.DeleteSong).Methods("DELETE")
	router.HandleFunc("/api/song/update", songHandler.AddSong).Methods("PUT")

	log.Println("Server started at http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
