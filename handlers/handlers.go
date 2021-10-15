package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"log"
	"net/http"
	"time"
)




var jwtKey = []byte("my_secret_key")

var users = map[string]string {
	"user1":"password1",
	"user2":"password2",
}

type Credentials struct {
	Password string `json:"password"`
	Username string `json:"username"`
}
type Claims struct {
	Username string `json:"username"`
	jwt.StandardClaims
}

func Signin(w http.ResponseWriter, r *http.Request) {
	var cred Credentials
	log.Println("started signing process")
	err := json.NewDecoder(r.Body).Decode(&cred)
	if err != nil {
		log.Println("structure of json body is wrong")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	expectedPassoword , ok := users[cred.Username]

	if !ok || expectedPassoword != cred.Password {
		log.Println("Username or password is wrong")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Println("Signin success")
	expirationTime := time.Now().Add(5 * time.Minute)

	// create the new jwt token
	claims := Claims{
		Username: cred.Username,
		StandardClaims : jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)

	// create jwt string
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	// setup cookies, with expiration time
	log.Println("Setting up cookies")
	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: tokenString,
		Expires: expirationTime,
	})


}
func Welcome(w http.ResponseWriter, r *http.Request) {
	log.Println("Welcome started")
	c , err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie{
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("User is unauthorized")
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		log.Println("http bad request")
		return
	}
	tokenString := c.Value
	// initialize the new instance of claims
	claims := &Claims{}
	tokken, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tokken.Valid {
		log.Println("Token is invalid")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Println("Successfully welcome")
	w.Write([]byte (fmt.Sprintf("Welcome %s!", claims.Username)))
}

func Refresh(w http.ResponseWriter, r *http.Request) {
log.Println("Refresh started")

	c , err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("request is unauthorized")
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		log.Println("http bad request")
		return
	}
	tokenstring := c.Value
	claims := &Claims{}

	token , err := jwt.ParseWithClaims(tokenstring, claims, func(token *jwt.Token)(interface{}, error){
		return jwtKey, nil
	})
	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			return
		}
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !token.Valid {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("Issuing new token")

	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	new_token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	newTokenString , err := new_token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	http.SetCookie(w, &http.Cookie{
		Name: "token",
		Value: newTokenString,
		Expires: expirationTime,
	})
}