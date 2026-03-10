package handlers

import (
	"encoding/json"
	"log"
	"net/http"

	"example.com/m/data"
)

type Products struct{
	l *log.Logger
}

// constructor
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// method
func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	lp := data.GetProduct()
	// how to print the response, on slice data, with simpliest method ?
	d, err := json.Marshal(lp)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
	w.Write(d)
}

