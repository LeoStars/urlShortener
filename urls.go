package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	. "strings"
)

type URL struct { // JSON-структура, содержащая поля:
	ID int `json:"id"` // идентификатора
	Address string `json:"address"` // адреса исходного
	Short string `json:"short"` // адреса зашифрованного
}

// функция реверса массив для алгоритма base62
func reverse (arr[]int) []int {
	var reverseString []int
	for i := len(arr) - 1; i >= 0; i--  {
		reverseString = append(reverseString, arr[i])
	}
	return reverseString
}

// Алгоритм base62
// имеем словарь,  в котором представлено 62 ключа и, значит, 62 элемента
// каждый элемент - это символ из последовательности [a..z, A..Z, 0..9]
// каждый ключ - это число от 0 до 62
// при занесении мы берём следующий идентификатор из JSON
// после последнего в предыдущей БД
// Например, последний в БД сейчас - 54
// значит, добавится 55
// После всего этого переводим чсило в 62-ичную систему счисления
// путём получения остатков от деления на 62 каждый раз
// реверсим полученный массив остатков
// и по ключам, соответствующим цислам в ревёрснутом массиве, находим буквы для кода
func base62 (code_table map[int]string, id int) string {
	var code []int
	for ; id > 0; {
		code = append(code, id%62)
		id /= 62
	}
	code = reverse(code)
	result := ""
	for _, i := range code {
		result += code_table[i]
	}
	return result // "http://avi.to/" + result
}

// создаём словарь для алгоритма base62
func makingBaseMap ()map[int]string {
	m := make(map[int]string)
	for i := 0; i < 51; i++ {
		if i < 26 {
			m[i] = string(i + 97)
		} else {
			m[i] = string(i + 39)
		}
	}
	for i:= 51; i <= 61; i++ {
		m[i] = string(i - 3)
	}
	return m
}

// проверка валидности введённой ссылки (которая исходная)
func validateURL (check string) (string, error){
	u, err := url.Parse(check) // тут ссылку парсим
	if u.Scheme == "" { // тут приводим её к виду с http://
		check = "http://" + check
	}
	u, err = url.ParseRequestURI(check) // снова парсим, чтобы учесть http://
	// нижу мы проверяем наличие точки в ссылке как обязательно части
	// точка не в конце и не в начале, а в середине
	// хост тоже пустовать не должен. http:// или /src/myProject в ссылках не пройдут!
	if err != nil || u.Host == "" || ContainsAny(check, ".") == false || HasSuffix(u.Host, ".") == true || HasPrefix(u.Host, ".") == true {
		err = errors.New("Wrong Address!")
	}
	return check + "/", err
}

// редирект-функция, достающая закодированную часть (Path) из адреса страницы
// и находящая его в БД, а затем перенаправлящая по найденной ссылке
func redirect(w http.ResponseWriter, r *http.Request) {
	redirecting := r.URL.Path[1:]
	needed := s.data.getNeeded(redirecting)
	http.Redirect(w, r, needed, 301)
	//ниже функция, открывающая наш URL сразу в окне браузера (в случае ввода из командной строки через curl
	exec.Command("explorer", needed).Run()
}

// функция декодирования URL по base62
func decode(w http.ResponseWriter, r *http.Request) {
	// структура запроса, который мы отправляем на сервер
	type Req struct {
		Address string `json:"address"`
	}
	var req Req
	// создаём декодер для парсинга json из запроса
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		panic(err)
	}
	// создаём словарь, проверяем валидность
	m := makingBaseMap()
	add, err := validateURL(req.Address)
	if err != nil {
		panic(err)
	}
	// так как у нас два запроса к базе данных, используем транзакции
	// здесь транзакцию начинаем
	tx, err := s.data.db.Begin()
	if err != nil {
		panic(err)
	}
	// в случае ошибки откатывает транзакцию
	defer tx.Rollback()
	// получаем максимальный (последний ID)
	id := s.data.getMax()
	// создаём структуру для нового адреса
	fmt.Println(id)
	newStruct := URL{ID: id, Address: add, Short: base62(m, id)}
	// добавляем новый адрес на сервер в БД
	err = s.data.newAddress(newStruct.Address, newStruct.Short)
	// закрываем транзакции
	tx.Commit()
	if err != nil {
		panic(err)
	}
	log.Println("Successfully created an URL")
}

//в этой функции всё то же самое, что и в decode, но base62 не нужен
func custom(w http.ResponseWriter, r *http.Request) {
	type Req struct {
		Address string `json:"address"`
		Short string `json:"short"`
	}
	var req Req
	decoder := json.NewDecoder(r.Body)
	err := decoder.Decode(&req)
	if err != nil {
		panic(err)
	}
	add, err := validateURL(req.Address)
	if err != nil {
		panic(err)
	}
	tx, err := s.data.db.Begin()
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()
	id := s.data.getMax()
	newStruct := URL{ID: id, Address: add, Short: req.Short}
	err = s.data.newAddress(newStruct.Address, newStruct.Short)
	if err != nil {
		panic(err)
	}
	tx.Commit()
	log.Println("Successfully created an URL")
}