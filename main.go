package main

import (
	"io/ioutil"
	"log"
	"net/http"

	"github.com/k0kubun/pp"
)

func main() {

	res, err := http.Get("https://api.github.com/users/ryym")
	if err != nil {
		log.Fatal(err)
	}
	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal(err)
	}
	json := string(body)

	lexer := NewLexer(json)
	result, err := Parse(lexer)
	if err != nil {
		log.Fatal(err)
	}

	pp.Print(result)
}
