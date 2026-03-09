package main

import (
	"log"
	"net/http"
	"os"

	"example.com/m/handlers"
)

func main(){

	// http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request){
	// 	// log.Println("Hello World")
	// 	d, err := io.ReadAll(r.Body)
	// 	if err != nil {
	// 		log.Fatal(err)
	// 	}
	// 	fmt.Fprintf(w, "%s", d)
	// })
	
	l := log.New(os.Stdout, "product-api", log.LstdFlags)

	hh := handlers.NewHello(l)
	gh := handlers.NewGoodbye(l)

	sm := http.NewServeMux()
	sm.Handle("/", hh)
	sm.Handle("/goodbye", gh)

	http.ListenAndServe(":8080", sm)
}