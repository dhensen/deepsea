package controllers

import (
	"encoding/json"
	"local/deepsea/api/db"
	"local/deepsea/api/k8s"
	"log"
	"net/http"
)

const (
	HEALTH_OK     = "ok"
	HEALTH_NOT_OK = "nok"
)

func Health(w http.ResponseWriter, r *http.Request) {
	// database health check
	dbStatus := HEALTH_OK
	if dbErr := db.DB.DB().Ping(); dbErr != nil {
		dbStatus = HEALTH_NOT_OK
		log.Println(dbErr.Error())
	}

	// kubernetes canary check: ask version and if you get an error something is not ok
	k8sStatus := HEALTH_OK
	if _, err := k8s.K8SDiscoveryClient.ServerVersion(); err != nil {
		k8sStatus = HEALTH_NOT_OK
		log.Println(err.Error())
	}

	isHealthy := dbStatus == HEALTH_OK && k8sStatus == HEALTH_OK
	var status string
	if isHealthy {
		status = HEALTH_OK
	} else {
		status = HEALTH_NOT_OK
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]interface{}{
		"services": map[string]string{
			"database":   dbStatus,
			"kubernetes": k8sStatus,
		},
		"status": status,
	})
}
