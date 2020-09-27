package main

import (
	"database/sql"
)

type database struct {
	db *sql.DB
}

// функция создания нового адреса в БД
func (d *database) newAddress (address string, short string) error {
	_, err := d.db.Exec("INSERT INTO url (address, short) VALUES($1, $2)", address, short)
	if err != nil {
		return err
	}
	return err
}

// функция извлечения из БД нужного адреса
func (d *database) getNeeded (address string) string {
	row := d.db.QueryRow("SELECT * FROM url WHERE short = ($1)", address)
	// здесь работаем с адресами, потому что QueryRow вовзращает указатель
	url := &URL{}
	err := row.Scan(&url.ID, &url.Address, &url.Short)
	if err != nil{
		panic(err)
	}
	return url.Address
}

// функция получения максимального (последнего ID) из БД
func (d *database) getMax () int {
	row := d.db.QueryRow("SELECT max(id) FROM url")
	var maxId int
	err := row.Scan(&maxId)
	if err != nil {
		panic(err)
	}
	return maxId
}

