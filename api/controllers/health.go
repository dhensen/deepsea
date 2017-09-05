package controllers

import (
	"encoding/json"
	"local/deepsea/api/db"
	"net/http"
	"strconv"
)

func Health(w http.ResponseWriter, r *http.Request) {
	var dbStatus string
	if dbErr := db.DB.DB().Ping(); dbErr != nil {
		dbStatus = dbErr.Error()
	} else {
		dbStatus = "ok"
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"services": map[string]string{"database": dbStatus},
		"status":   strconv.FormatBool(dbStatus == "ok"),
	})
}
