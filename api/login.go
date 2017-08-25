package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

var users_storage = map[string]User{
	"dhensen": User{
		Name:     "dhensen",
		Password: "1234",
	},
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.PostFormValue("username")
	password := r.PostFormValue("password")

	if username == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "username is required",
		})
		return
	}

	if password == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(map[string]string{
			"message": "password is required",
		})
		return
	}

	if user, ok := users_storage[username]; ok {
		//do something here
		if user.Password == password {

			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"username":    username,
				"access_type": "admin",
			})

			// TODO: put this secret in config
			hmacSecret := []byte("superfoobar9000")
			tokenString, err := token.SignedString(hmacSecret)
			if err != nil {
				log.Println(err)
				// TODO: handle error
			}

			tokenCookie := &http.Cookie{
				Name:     "access_token",
				Value:    tokenString,
				Expires:  time.Now().Add(7 * 24 * time.Hour),
				Secure:   false,
				HttpOnly: true,
			}
			http.SetCookie(w, tokenCookie)

			// logged in
			json.NewEncoder(w).Encode(map[string]string{
				"message":      "you are logged in dude",
				"access_token": tokenString,
			})
			return
		}

		json.NewEncoder(w).Encode(map[string]string{
			"message": "wrong username or password",
		})
		return
	}

	json.NewEncoder(w).Encode(map[string]string{
		"message": "wrong username or password",
	})
	return
}
