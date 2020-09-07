package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	. "strings"
)


type URLs struct {
	URLs []URL `json:"URLs"`
}

type URL struct {
	ID int `json:"id"`
	Address string `json:"address"`
	Short string `json:"short"`
}


func reverse (arr[]int) []int {
	var reverseString []int
	for i := len(arr) - 1; i >= 0; i--  {
		reverseString = append(reverseString, arr[i])
	}
	return reverseString
}

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

func jsonRead () URLs{
	byteValue, _ := ioutil.ReadFile("URLs.json")
	var url URLs
	err := json.Unmarshal(byteValue, &url)
	if err != nil {
		panic("Unmarshall was bad")
	}
	return url
}

func jsonWrite (url URLs) {
	decodeJson, err := json.MarshalIndent(url, "", "    ")
	if err != nil {
		fmt.Println(err)
		panic("There is an error!")
	}
	err = ioutil.WriteFile("URLs.json", decodeJson, 0644)
	return
}

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

func jsonAppend(url URLs, m map[int]string) URLs{
	var id int = url.URLs[len(url.URLs) - 1].ID + 1
	var add string
	fmt.Fscan(os.Stdin, &add)
	add = validateURL(add)
	newStruct := URL{ID: id, Address: add, Short: base62(m, id)}
	url.URLs = append(url.URLs, newStruct)
	return url
}

func jsonAppendCustom(url URLs) URLs{
	var id int = url.URLs[len(url.URLs) - 1].ID + 1
	var add string
	var custom string
	fmt.Fscan(os.Stdin, &add)
	add = validateURL(add)
	fmt.Fscan(os.Stdin, &custom)
	redirectingUrl := custom // "http://avi.to/" + custom
	newStruct := URL{ID: id, Address: add, Short: redirectingUrl}
	url.URLs = append(url.URLs, newStruct)
	return url
}

func validateURL (check string) string{
	u, err := url.Parse(check)
	if u.Scheme == "" {
		check = "http://" + check
	}
	u, err = url.ParseRequestURI(check)
	if err != nil || u.Host == "" || ContainsAny(check, ".") == false || HasSuffix(u.Host, ".") == true || HasPrefix(u.Host, ".") == true {
		panic("URL does not exist!")
	}
	return check + "/"
}

func redirect(w http.ResponseWriter, r *http.Request) {
	redirecting := r.URL.Path[1:]
	needed := findURL(jsonRead(), redirecting)
	http.Redirect(w, r, needed, 301)
}

func findURL (UrlSlice URLs, redir string) string {
	for _, i := range(UrlSlice.URLs){
		if Compare(i.Short, redir) == 0 {
			fmt.Println(i.Address)
			return i.Address
		}
	}
	return ""
}

func main(){
	fmt.Println("What do you want to do?")
	fmt.Println("1. Decode your URL")
	fmt.Println("2. Create custom URL")
	fmt.Println("3. Redirect")
	var choose byte
	fmt.Scanf("%d", &choose)
	switch choose {
	case 1:
		m := makingBaseMap()
		url := jsonRead()
		url = jsonAppend(url, m)
		jsonWrite(url)
	case 2:
		url := jsonRead()
		url = jsonAppendCustom(url)
		jsonWrite(url)
	case 3:
		http.HandleFunc("/", redirect)
		err := http.ListenAndServe(":9090", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}
	}
}