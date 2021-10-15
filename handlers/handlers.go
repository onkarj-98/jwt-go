package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key")

var users = map[string]string{
	"user1": "password1",
	"user2": "password2",
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
	var creds Credentials
	log.Println("Started singin process...")

	// this line decodes the incoming request json body into the our structure
	err := json.NewDecoder(r.Body).Decode(&creds)
	if err != nil {
		// If the structure of the body is wrong, return an HTTP error
		w.WriteHeader(http.StatusBadRequest)
		log.Println("Structure of the JSON body is wrong")
		return
	}
	// get the expected password from out in memory map
	// here ok gives an error , if its not in memory map
	// cred we are taking from request
	expectedPassword, ok := users[creds.Username]

	// check passwoerd now
	if !ok || expectedPassword != creds.Password {
		log.Println("Wrong username or password")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}
	log.Println("Adding expiration time")
	expirationTime := time.Now().Add(5 * time.Minute)

	// Create the JWT claims, which includes the username and expiry time
	claims := &Claims{

		Username: creds.Username,
		StandardClaims: jwt.StandardClaims{

			ExpiresAt: expirationTime.Unix(),
		},
	}

	// Declare the token with the algorithm used for signing and the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the jwt string
	log.Println("Issueing new token")
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	// Finally, we set the clinet cookie for token as jwt generated
	// also setting an expiry time which is the same as the token iteself
	log.Println("Setting up cookies")
	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

}

func Welcome(w http.ResponseWriter, r *http.Request) {
	log.Println("Welcome started")
	// we can obtain the session token from the cookies, which come with every request
	c, err := r.Cookie("token")
	if err != nil {
		if err == http.ErrNoCookie {
			// If cookie is not set, return an unauthorized status
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("User is unauthorized")
			return
		}
		// for any other type of error, return a bad request status
		w.WriteHeader(http.StatusBadRequest)
		log.Println("http bad request")
		return
	}
	// get the JWT string from the cookie
	tokenstrng := c.Value

	// Initialize a new instance of `Claims`
	claims := &Claims{}

	tokken, err := jwt.ParseWithClaims(tokenstrng, claims, func(token *jwt.Token) (interface{}, error) {
		return jwtKey, nil
	})

	if err != nil {
		if err == jwt.ErrSignatureInvalid {
			w.WriteHeader(http.StatusUnauthorized)
			log.Println("Signature is invalid")
			return
		}
		log.Println("http Bad request")
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if !tokken.Valid {
		log.Println("Token is invalid")
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// Finally return the welcome message to user, along with their usename given in token
	log.Println("Successfully Welcome")
	w.Write([]byte(fmt.Sprintf("Welcome %s!", claims.Username)))
}

// Refreshing the token

func Refresh(w http.ResponseWriter, r *http.Request) {
	// some code of this function is same as welcome
	log.Println("Refresh started")

	c, err := r.Cookie("token")

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

	tknstr := c.Value
	claims := &Claims{}

	tkn, err := jwt.ParseWithClaims(tknstr, claims, func(token *jwt.Token) (interface{}, error) {
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

	if !tkn.Valid {
		w.WriteHeader(http.StatusUnauthorized)
		return
	}

	// we ensure that new token will be gnererated only , if it has life of less than 30 sec
	if time.Unix(claims.ExpiresAt, 0).Sub(time.Now()) > 30*time.Second {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	log.Println("Issueing new token")
	// |Now create a new token for the current user, with renewed expiration time
	expirationTime := time.Now().Add(5 * time.Minute)
	claims.ExpiresAt = expirationTime.Unix()
	token := jwt.NewWithClaims(jwt.SigningMethodES256, claims)
	tokenString, err := token.SignedString(jwtKey)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	http.SetCookie(w, &http.Cookie{
		Name:    "token",
		Value:   tokenString,
		Expires: expirationTime,
	})

}
