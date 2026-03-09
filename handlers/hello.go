package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)


type Hello struct{
	l *log.Logger
}

func NewHello(l *log.Logger)  *Hello{
	return &Hello{}
}

func (h *Hello) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Hello %s", d)
}
