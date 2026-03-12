package handlers

import (
	// "encoding/json"
	"log"
	"net/http"
	"strconv"

	"example.com/m/data"
	"github.com/gorilla/mux"
)

type Products struct{
	l *log.Logger
}

type KeyProduct struct{}

// constructor
func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// // method
// func (p *Products) ServeHTTP(w http.ResponseWriter, r *http.Request) {
// 	if r.Method == http.MethodGet {
// 		p.GetAllProduct(w, r)
// 	}

// 	if r.Method == http.MethodPost {
// 		p.AddProduct(w, r)
// 	}

// 	if r.Method == http.MethodPut {
// 		// no need to matching regex because gorilla already support that
// 		p.UpdateProduct(w, r)
// 		return
// 	}

// 	// currently return error if not specified
// 	w.WriteHeader(http.StatusMethodNotAllowed)
// }

func (p *Products) GetAllProduct(w http.ResponseWriter, r *http.Request) {
	lp := data.GetAllProduct()

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

func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // like ctx package
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
		return
	} 

	product  := &data.Product{}
	// decode from r.body to struct (data)
	err = product.FromJson(r.Body)
	if err != nil {
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
		return
	}

	// product := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, product)
	if err != nil {
		http.Error(w, "Unable to update product", http.StatusBadRequest)
		return
	}


	
}

