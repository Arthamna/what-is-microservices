package handlers

import (
	// "encoding/json"
	"context"
	"fmt"
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

func (p *Products) GetAllProduct(w http.ResponseWriter, r *http.Request) {
	lp := data.GetAllProduct()

	// serializes json
	err := lp.ToJson(w)
	if err != nil{
		http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	}
}

func (p *Products) AddProduct(w http.ResponseWriter, r *http.Request) {

	// lp := &data.Product{}

	// err := lp.FromJson(r.Body)
	// if err != nil{
	// 	http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	// }

	// create context copy with passed data struct
	prod := r.Context().Value(KeyProduct{}).(data.Product)

	data.AddProduct(&prod)
}

func (p *Products) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r) // like ctx package
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Unable to convert id", http.StatusBadRequest)
		return
	} 

	// product  := &data.Product{}
	// // decode from r.body to struct (data)
	// err = product.FromJson(r.Body)
	// if err != nil {
	// 	http.Error(w, "Unable to marshal json", http.StatusInternalServerError)
	// 	return
	// }

	product := r.Context().Value(KeyProduct{}).(data.Product)

	err = data.UpdateProduct(id, &product)
	if err != nil {
		http.Error(w, "Unable to update product", http.StatusBadRequest)
		return
	}

}

func (p *Products) MiddlewareProduct(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		product := &data.Product{}
		err := product.FromJson(r.Body)
		if err != nil {
			p.l.Println("[ERROR] deserializing product", err)
			http.Error(w, "Error reading product", http.StatusBadRequest)
			return
		}

		// validation
		err = product.Validate()
		if err != nil {
			p.l.Println("[ERROR] validating product", err)
			http.Error(
				w, 
				fmt.Sprintf("Error validating product %s", err),
				http.StatusBadRequest,
			)
			return
		}

		//add product to ctx
		ctx := context.WithValue(r.Context(), KeyProduct{}, *product)
		r = r.WithContext(ctx)

		// call next handler, which can be another middleware or final handler
		next.ServeHTTP(w, r)
	})

}
