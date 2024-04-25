package main

import (
	"crud/servidor"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func main() {

	router := mux.NewRouter()
	router.HandleFunc("/usuarios", servidor.CriarUsuario).Methods(http.MethodPost)

	fmt.Println("Escutando na porta 8080")
	log.Fatal(http.ListenAndServe(":8080", router))
}
