package main

import (
	"encoding/json"
	"net/http"

	uuid "github.com/google/uuid"
)

type ImagePreset struct {
	ID          uuid.UUID
	DisplayName string
	Image       string
}

var imagePresets = map[string]ImagePreset{
	"a8ef3baa-6fbd-4a9b-ad94-247c9273c57d": ImagePreset{
		ID:          uuid.Must(uuid.Parse("a8ef3baa-6fbd-4a9b-ad94-247c9273c57d")),
		DisplayName: "Wordpress 4.8.1 running PHP 5.6 on Apache",
		Image:       "wordpress:4.8.1-php5.6-apache",
	},
	"1c98846e-b60e-47bc-a1ca-707057ee70b8": ImagePreset{
		ID:          uuid.Must(uuid.Parse("1c98846e-b60e-47bc-a1ca-707057ee70b8")),
		DisplayName: "Wordpress 4.8.1 running PHP 7.1 on Apache",
		Image:       "wordpress:4.8.1-php7.1-apache",
	},
}

func ListContainerPresets(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(imagePresets)
}
