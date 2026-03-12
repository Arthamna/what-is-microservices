package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"example.com/m/handlers"
	"github.com/gorilla/mux"
)

func main(){
	
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	ph := handlers.NewProducts(l)

	// sm := http.NewServeMux()
	sm := mux.NewRouter()

	getRouter := sm.Methods(http.MethodGet).Subrouter()
	getRouter.HandleFunc("/", ph.GetAllProduct)
	
	putRouter := sm.Methods(http.MethodPut).Subrouter()
	putRouter.HandleFunc("/{id:[0-9]+}", ph.UpdateProduct)

	addRouter := sm.Methods(http.MethodPost).Subrouter()
	addRouter.HandleFunc("/add/", ph.AddProduct)
	
	s := http.Server{
		Addr: ":8080",
		Handler: sm,
		IdleTimeout: 120*time.Second,
		ReadTimeout: 1*time.Second,
		WriteTimeout: 1*time.Second,
	}

	
	go func() {
		if err := s.ListenAndServe(); err != nil {
			l.Fatal(err)
		}
	}()
		
	sigChan := make(chan os.Signal)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, os.Kill)

	sig := <- sigChan
	l.Println("Received terminate, graceful shutdown", sig)
	
	// default context to graceful shutdown
	tc, _ := context.WithTimeout(context.Background(), 30*time.Second)

	s.Shutdown(tc)
}