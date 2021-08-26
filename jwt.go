package main

import (
	"encoding/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
	"log"
	"net/http"
	"time"
)

type User struct{
	Name string `json:"name"`
	Pass string `json:"pass"`
}

var(
	users = []User{
		{"ab","ab"},
		{"c","d"},
		{"e","f"},
	}
	tokenAuth *jwtauth.JWTAuth
)

func check(u User)bool{
	for _,v := range users{
		if u==v {
			return true
		}
	}
	return false
}

func init(){
	tokenAuth = jwtauth.New("HS256",[]byte("mykey"),nil)
}

func main(){
	r := chi.NewRouter()

	r.Post("/login",makeLogin)

	r.Route("/",au)
	log.Fatal(http.ListenAndServe(":8081",r))
}

func au(r chi.Router){

	r.Use(jwtauth.Verifier(tokenAuth))
	r.Use(jwtauth.Authenticator)

	r.Get("/", readData)
	r.Post("/createUser", createNewUser)
}

func makeLogin(w http.ResponseWriter,r *http.Request){
	var u User
	json.NewDecoder(r.Body).Decode(&u)
	for _,v := range users{
		if u == v {
			//gen cookie
			_,tokenString,e := tokenAuth.Encode(map[string]interface{}{"aud":u.Name})
			if e != nil {
				panic(e)
			}
			// set cookie
			http.SetCookie(w,&http.Cookie{
				Name: "jwt",
				Value: tokenString,
				Expires: time.Now().Add(time.Second * 500),
			})
			w.Write([]byte(tokenString))
			return
		}
	}
	w.Write([]byte("username not found"))
}

func createNewUser(w http.ResponseWriter,r *http.Request){
	var u User
	json.NewDecoder(r.Body).Decode(&u)
	users = append(users,u)
	w.Write([]byte("successfully created new user"))
}

func readData(w http.ResponseWriter,r *http.Request){
	w.Header().Set("Content-Type","application/json")
	json.NewEncoder(w).Encode(users)
}