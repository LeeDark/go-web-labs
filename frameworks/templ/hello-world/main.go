package main

import (
	"fmt"
	"net/http"

	"github.com/a-h/templ"
)

func main() {
	component := hello("John")

	http.Handle("/", templ.Handler(component))

	fmt.Println("Listening on :3001")
	err := http.ListenAndServe(":3001", nil)
	if err != nil {
		fmt.Println(err)
	}
}
