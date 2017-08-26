package main

import (
	"encoding/json"
	"log"
	"net/http"
	"time"

	jwtmiddleware "github.com/auth0/go-jwt-middleware"
	jwt "github.com/dgrijalva/jwt-go"
)

const jwtSecret = "superfoobar9000"
const jwtCookieKey = "access_token"

func CookieExtractor(jwtKey string) jwtmiddleware.TokenExtractor {
	return func(r *http.Request) (string, error) {
		cookie, err := r.Cookie(jwtKey)
		if err != nil {
			return "", nil
		}
		return cookie.Value, nil
	}
}

func Authenticated(f func(w http.ResponseWriter, r *http.Request)) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		jwtMiddleware := GetJWTMiddleware()
		err := jwtMiddleware.CheckJWT(w, r)

		// If there was an error, do not continue.
		if err != nil {
			return
		}

		f(w, r)
	}
}

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
			hmacSecret := []byte(jwtSecret)
			tokenString, err := token.SignedString(hmacSecret)
			if err != nil {
				log.Println(err)
				// TODO: handle error
			}

			tokenCookie := &http.Cookie{
				Name:     jwtCookieKey,
				Value:    tokenString,
				Expires:  time.Now().Add(7 * 24 * time.Hour),
				Secure:   false,
				HttpOnly: true,
			}
			http.SetCookie(w, tokenCookie)

			// logged in
			json.NewEncoder(w).Encode(map[string]string{
				"message":    "you are logged in dude",
				jwtCookieKey: tokenString,
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

func Logout(w http.ResponseWriter, r *http.Request) {
	tokenCookie := &http.Cookie{
		Name:     jwtCookieKey,
		Value:    "",
		Expires:  time.Unix(0, 0),
		Secure:   false,
		HttpOnly: true,
	}
	http.SetCookie(w, tokenCookie)
}
