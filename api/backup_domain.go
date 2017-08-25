package main

import (
	"fmt"
	"net/http"

	mux "github.com/gorilla/mux"
)

func ListBackups(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "%v", vars["id"])
}

func CreateBackup(w http.ResponseWriter, r *http.Request) {

}
