package main

import (
	"net/http"
)

func (s *server) setupRoutes() {
	// здесь настраиваем роутинг на три направления:
	// декодирование без кастомнрой ссылки
	http.HandleFunc("/decode", decode)
	// кастомизация ссылки
	http.HandleFunc("/custom", custom)
	// перенаправление по ссылке из БД
	http.HandleFunc("/", redirect)
}
