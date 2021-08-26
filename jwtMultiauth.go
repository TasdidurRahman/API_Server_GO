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
		{"a","b"},
		{"c","d"},
		{"e","f"},
	}
	admins = []User{
		{"g","h"},
		{"i","j"},
	}
	tokenAuth *jwtauth.JWTAuth
	tokenAuth2 *jwtauth.JWTAuth
)

func checkUser(u User)bool{
	for _,v := range users{
		if u==v {
			return true
		}
	}
	return false
}
func checkAdmin(a User)bool{
	for _,v := range admins{
		if a==v {
			return true
		}
	}
	return false
}

func init(){
	tokenAuth = jwtauth.New("HS256",[]byte("mykey"),nil)
	tokenAuth2 = jwtauth.New("HS256",[]byte("okidoki"),nil)
}

func main(){
	r := chi.NewRouter()

	r.Post("/adminLogin",logAdmin)
	r.Post("/userLogin",logUser)

	r.Group(us)
	r.Group(adm)

	log.Fatal(http.ListenAndServe(":8081",r))
}

func us(r chi.Router){

	r.Use(jwtauth.Verifier(tokenAuth))
	r.Use(jwtauth.Authenticator)

	r.Get("/readUser", readUser)
}
func adm(r chi.Router){

	r.Use(jwtauth.Verifier(tokenAuth2))
	r.Use(jwtauth.Authenticator)

	r.Get("/readAdmin", readAdmin)

}

func readUser(w http.ResponseWriter,r *http.Request){
	w.Write([]byte("user can read..."))
}
func readAdmin(w http.ResponseWriter,r *http.Request){
	w.Write([]byte("admin can read..."))
}

func logAdmin(w http.ResponseWriter,r *http.Request){
	var a User
	json.NewDecoder(r.Body).Decode(&a)
	for _,v := range admins{
		if a == v {
			//gen cookie
			_,tokenString,e := tokenAuth2.Encode(map[string]interface{}{"aud":a.Name})
			if e != nil {
				panic(e)
			}
			// set cookie
			http.SetCookie(w,&http.Cookie{
				Name: "jwt",
				Value: tokenString,
				Expires: time.Now().Add(time.Second * 50000),
			})
			w.Write([]byte(tokenString))
			return
		}
	}
	w.Write([]byte("admin not found"))
}

func logUser(w http.ResponseWriter,r *http.Request){
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
