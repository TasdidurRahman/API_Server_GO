package main

import (
"encoding/json"
"fmt"
"github.com/go-chi/chi/v5"
"log"
"net/http"
"strconv"
"time"
)

type User struct{
	Name string	`json:"name"`
	Pass string	`json:"pass,omitempty"`
	Type string	`json:"type,omitempty"`
}

type Product struct{
	Category string	`json:"category,omitempty"`
	Brand string	`json:"brand,omitempty"`
	Model string	`json:"model,omitempty"`
	Price int		`json:"price,omitempty"`
	Count int		`json:"count,omitempty"`
}

type Invoice struct{
	User User	`json:"user"`
	Product Product	`json:"product"`
	Quantity int	`json:"quantity,omitempty"`
	PurchaseDate time.Time	`json:"purchase_date,omitempty"`
}

var (
	users = []User{
		{"a","b","admin"},
		{"c","d","general"},
		{"e","f","general"},
		{"g","h","general"},
	}

	products = []Product{
		{"monitor","DELL","d22",15000,5},
		{"monitor","LG","g94",16000,1},
		{"mouse","A4tech","super",350,10},
		{"keyboard","Delux","soft",400,3},
		{"headphone","bits","mux35",22000,2},
		{"Router","TpLink","Archer10",3000,3},
	}

	invoices []Invoice
)

func makeInvoice(u User,p Product,q int){
	var I Invoice
	I.User = u
	I.Product = p
	I.Quantity = q
	I.PurchaseDate = time.Now()
	invoices = append(invoices,I)
}

func (u *User)searchInvoice() (invList []Invoice){
	for _,j := range invoices{
		if *u == j.User{
			invList = append(invList,j)
		}
	}
	return
}

func matchIncompleteCompleteProducts(p Product,q Product) bool{ // p : incomplete , q : complete
	if p.Category==""&&p.Brand==""&&p.Model=="" {
		return false
	}
	if (p.Category==q.Category || p.Category=="")&&
		(p.Model==q.Model || p.Model=="")&&
		(p.Brand==q.Brand || p.Brand==""){
		return true
	}
	return false
}

func (p *Product)search() (ps []Product){
	for _,j := range products{
		if matchIncompleteCompleteProducts(*p,j) {
			ps = append(ps,j)
		}
	}
	return
}

func (p *Product)changeCount(c int) {
	for i,j := range products{
		if matchIncompleteCompleteProducts(*p,j){
			//for ;c>0;c-- { // check
			products[i].Count += c
			//}
			return
		}
	}
}

func (p *Product)deleteProduct() {
	for i,j := range products{
		if matchIncompleteCompleteProducts(*p,j){
			products = append(products[:i],products[i+1:]...)
			return
		}
	}
}


func main() {
	r := chi.NewRouter()

	r.Post("/createUser",createUser)
	r.Post("/login",login)
	r.Put("/update/{c}",updateProduct)
	r.Post("/add",addNewProduct)
	r.Delete("/",deleteProduct)
	r.Get("/invoice",getInvoice)
	r.Get("/search",searchProduct)
	r.Post("/buy",buyProduct)

	log.Fatal(http.ListenAndServe(":8081",r))
}

func createUser(w http.ResponseWriter,r *http.Request)  {
	var u User
	json.NewDecoder(r.Body).Decode(&u)
	users = append(users,u)
	w.Write([]byte("successfully added user"))
	//for _,i := range users{
	//	fmt.Println(i)
	//}
	return
}
func login(w http.ResponseWriter,r *http.Request)  {
	var u User
	json.NewDecoder(r.Body).Decode(&u)
	for _,i := range users{
		if i.Name == u.Name && i.Pass == u.Pass {
			//gen token
			//gen cookie
			w.Write([]byte("login successfull"))
			return
		}
	}
	w.Write([]byte("no such credential"))
}
func updateProduct(w http.ResponseWriter,r *http.Request)  {
	var p Product
	c := chi.URLParam(r,"c")
	cc,_ := strconv.Atoi(c)
	json.NewDecoder(r.Body).Decode(&p)
	for i,j := range products{
		fmt.Println(j,p)
		if matchIncompleteCompleteProducts(p,j){
			products[i].changeCount(cc)
			json.NewEncoder(w).Encode(products)
			return
		}
	}

}
func addNewProduct(w http.ResponseWriter,r *http.Request)  {
	var p Product
	json.NewDecoder(r.Body).Decode(&p)
	products = append(products,p)
}
func deleteProduct(w http.ResponseWriter,r *http.Request)  {
	var p Product
	json.NewDecoder(r.Body).Decode(&p)
	p.deleteProduct()
}
func getInvoice(w http.ResponseWriter,r *http.Request)  {
	var u User
	json.NewDecoder(r.Body).Decode(&u)
	var ret []Invoice
	for i,j := range invoices{
		if j.User == u {
			ret = append(ret,invoices[i])
		}
	}
	json.NewEncoder(w).Encode(ret)
}
func searchProduct(w http.ResponseWriter,r *http.Request)  {
	var p Product
	json.NewDecoder(r.Body).Decode(&p)
	json.NewEncoder(w).Encode(p.search())
}
func buyProduct(w http.ResponseWriter,r *http.Request)  {
	var I Invoice
	json.NewDecoder(r.Body).Decode(&I)
	I.PurchaseDate = time.Now()
	invoices = append(invoices,I)
	I.Product.changeCount(I.Quantity)
}