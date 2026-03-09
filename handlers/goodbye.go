package handlers

import (
	"fmt"
	"io"
	"log"
	"net/http"
)


type Goodbye struct{
	l *log.Logger
}

func NewGoodbye(l *log.Logger)  *Goodbye{
	return &Goodbye{}
}

func (h *Goodbye) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	d, err := io.ReadAll(r.Body)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Fprintf(w, "Goodbye %s", d)
}
