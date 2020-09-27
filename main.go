package main

import (
	"log"
	"net/http"
)

// создаём структуру для нашего сервера
type server struct {
	data     *database
}

var s = &server{}

func main(){
	// открываем БД-сессию
	s.setupDB()
	// настраиваем маршрутизацию
	s.setupRoutes()
	// поднимаем сервер
	log.Fatal(http.ListenAndServe(":8080", nil))
}