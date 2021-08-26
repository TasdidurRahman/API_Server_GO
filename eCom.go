package main

import (
	"github.com/go-chi/chi/v5"
	"log"
	"net/http"
)

type User struct{
	Name string
	Pass string
	Type string
}

type Product struct{
	Category string
	Brand string
	Model string
	Price int
	Count int
}

type Invoice struct{
	User User
	Product Product
	Quantity int
	PurchaseDate string
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

func (u *User)searchInvoice() (invList []Invoice){
	for _,j := range invoices{
		if *u == j.User{
			invList = append(invList,j)
		}
	}
	return
}

func match(p Product,q Product) bool{
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
		if match(*p,j) {
			ps = append(ps,j)
		}
	}
	return
}

func (p *Product)increaseCount() {
	for i,j := range products{
		if match(*p,j){
			products[i].Count++
			return
		}
	}
}

func (p *Product)deleteProduct() {
	for i,j := range products{
		if match(*p,j){
			products = append(products[:i],products[i+1:]...)
			return
		}
	}
}

func addProduct(p Product){
	ret := p.search()
	if len(ret) != 0 {
		p.increaseCount()
	}else{
		p.Count = 1
		products = append(products,p)
	}
}

func main() {
	r := chi.NewRouter()

	

	log.Fatal(http.ListenAndServe(":8081",r))
}