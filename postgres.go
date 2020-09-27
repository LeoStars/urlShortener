package main

import (
	"database/sql"
	"log"
	_ "github.com/lib/pq"
)

// настройка БД для работы
func (s *server) setupDB() {
	log.Println("Setting up db....")
	var err error
	var db *sql.DB
	// вводим все данные нашей базы данные
	connStr := "user=postgres password=postgres dbname=url_data sslmode=disable"
	// открываем БД-сессию
	db, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
	// открыли сессию - создали ссылку на нужную нам базу данных
	s.data = &database{
		db: db,
	}
	// вывели успешное открытие сессии
	log.Println("Success")
}
