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
	if r.Method == http.MethodGet {
		p.GetProduct(w, r)
	}

	if r.Method == http.MethodPost {
		p.AddProduct(w, r)
	}

	// currently return error if not specified
	w.WriteHeader(http.StatusMethodNotAllowed)
}

func (p *Products) GetProduct(w http.ResponseWriter, r *http.Request) {
	lp := data.GetProduct()

	// serializes json
	err := lp.ToJson(w)
	if err != nil{
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {
	// take a template from slices (because we implement the method on slices)
	lp := &data.Product{}

	err := lp.FromJson(r.Body)
	if err != nil{
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}

	data.AddProduct(*lp)
}

