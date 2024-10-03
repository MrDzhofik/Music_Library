package main

import (
	"context"
	"log"
	"net/http"
	"os"

	"effectiveMobile/internal/handlers"
	"effectiveMobile/internal/storage"
	"effectiveMobile/internal/usecase"

	"github.com/gorilla/mux"
)

func setupRouter(songHandler *handlers.SongHandler) *mux.Router {
	router := mux.NewRouter()
	router.HandleFunc("/api/songs", songHandler.GetAllSongs).Methods("GET")
	router.HandleFunc("/api/song/{id:[0-9]+}", songHandler.GetSongByID).Methods("GET")
	router.HandleFunc("/api/song/add", songHandler.AddSong).Methods("POST")
	router.HandleFunc("/api/song/delete", songHandler.DeleteSong).Methods("DELETE")
	router.HandleFunc("/api/song/update", songHandler.UpdateSong).Methods("PUT")

	return router
}

func setupLoggers() (*log.Logger, *log.Logger, *os.File, *os.File) {
	// Открытие файлов логов
	infoFile, err := os.OpenFile("info.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	errorFile, err := os.OpenFile("error.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		log.Fatal(err)
	}

	infoLog := log.New(infoFile, "INFO\t", log.Ldate|log.Ltime)

	errorLog := log.New(errorFile, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	if _, err := infoFile.WriteString("------------------------------\n"); err != nil {
		errorLog.Print(err)
	}

	if _, err := errorFile.WriteString("------------------------------\n"); err != nil {
		errorLog.Print(err)
	}

	return infoLog, errorLog, infoFile, errorFile
}

func main() {
	// Настройка логгеров
	infoLog, errorLog, infoFile, errorFile := setupLoggers()
	defer infoFile.Close()
	defer errorFile.Close()

	// Подключаем базу
	conn, err := storage.GetPostgres()
	if err != nil {
		log.Panic("Ошибка соединения с базой данных")
	}

	defer conn.Close(context.Background())

	// Инициализация слоёв
	songStorage := storage.NewSongStorage(conn, infoLog, errorLog)
	songUsecase := usecase.NewSongUsecase(songStorage, infoLog, errorLog)
	songHandler := handlers.NewSongHandler(songUsecase, infoLog, errorLog)

	// Настройка роутера
	router := setupRouter(songHandler)

	// Создаем новую структуру http.Server, оставляем тот же адрес и роутер, а для ошибок используем наш логгер
	srv := &http.Server{
		Addr:     "localhost:8080",
		ErrorLog: errorLog,
		Handler:  router,
	}

	infoLog.Printf("Запуск сервера на %s", srv.Addr)
	// Вызываем метод ListenAndServe() от нашей новой структуры http.Server
	err = srv.ListenAndServe()
	errorLog.Fatal(err)
	errorLog.Fatal(http.ListenAndServe(":8080", router))
}
