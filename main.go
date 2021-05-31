package main

import (
	"context"
	"database/sql"
	// "encoding/json"
	"fmt"
	"html/template"
	"log"
	"net/http"

	// "strconv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

const (
	CONN_HOST        = "localhost"
	CONN_PORT        = "8080"
	DRIVER_NAME      = "mysql"
	DATA_SOURCE_NAME = "root:@/mydb"
)

var db *sql.DB
var ctx context.Context
var err error

var connectionError error

func init() {
	db, connectionError = sql.Open(DRIVER_NAME, DATA_SOURCE_NAME)
	if connectionError != nil {
		log.Fatal("error connecting to database : ", connectionError)
	}
}

func getCurrentDb(w http.ResponseWriter, r *http.Request) {
	rows, err := db.Query("SELECT DATABASE() as db")
	if err != nil {
		log.Print("error executing query :: ", err)
		return
	}
	var db string
	for rows.Next() {
		rows.Scan(&db)
	}
	fmt.Fprintf(w, "Current Database is :: %s", db)
}

func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/home.html")
		t.Execute(w, nil)
	}
}

// Signup Page
func Signup(w http.ResponseWriter, r *http.Request) {

	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/signup.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		name := r.FormValue("name")
		phone := r.FormValue("phone")
		email := r.FormValue("email")
		batch := r.FormValue("batch")
		password := r.FormValue("password")

		_, err = db.Exec("INSERT INTO alumni (Name, Phone, Email, Batch, Password, Is_verified) VALUES (?, ?, ?, ?, ?, ?)", name, phone, email, batch, password, "0")
		if err != nil {
			fmt.Println(err)
			http.Error(w, "Server error, unable to create your account.", 500)
			return
		}
		fmt.Println("User created successfully")
	}
}

// Login Page
func Login(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/login.html")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		// logic part of log in
		fmt.Println("Email:", r.Form["email"])
		fmt.Println("Password:", r.Form["password"])
		// getCurrentDb()
	}
}

type Alumni struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
	Batch string `json:"batch"`
}

// Alums Page
func Alums(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method:", r.Method) //get request method
	if r.Method == "GET" {
		t, _ := template.ParseFiles("templates/alums.html")
		t.Execute(w, nil)

		log.Print("reading records from database")
		rows, err := db.Query("SELECT Id, Name, Email, Phone, Batch fROM alumni;")
		// fmt.Println(rows)
		if err != nil {
			log.Print("error executing select query :: ", err)
			return
		}

		alumnees := []Alumni{}
		for rows.Next() {
			var Id int
			var Name string
			var Email string
			var Phone string
			var Batch string
			err = rows.Scan(&Id, &Name, &Email, &Phone, &Batch)
			alumni := Alumni{Id: Id, Name: Name, Email: Email, Phone: Phone, Batch: Batch}
			// fmt.Println(alumni)
			alumnees = append(alumnees, alumni)
		}

		// json.NewEncoder(w).Encode(alumnees)

		for i := range alumnees {
			emp := alumnees[i]
			myvar := map[string]interface{}{"Name": emp.Name, "Batch" : emp.Batch, "Phone" : emp.Phone, "Email" : emp.Email}
			t, _ := template.ParseFiles("templates/card.html")
			fmt.Println(emp)
			// t.Execute(w, myvar)
			if err := t.Execute(w, myvar); err != nil {
        			http.Error(w, err.Error(), 500)
        			return
    			}
			// fmt.Fprintf(w, "Id :%d \n Name : %s\n Email: %s \n Phone: %s \n Batch : %s\n", emp.Id, emp.Name, emp.Email, emp.Phone, emp.Batch)
			// fmt.Fprintln(w)
		}
		
		

		// fmt.Println(alumnees)
	}
	t, _ := template.ParseFiles("templates/procedure-add.html")
		t.Execute(w, nil)
}

func main() {
	r := mux.NewRouter()
	r.HandleFunc("/", Home)
	r.HandleFunc("/signup", Signup)
	r.HandleFunc("/login", Login)
	r.HandleFunc("/alums", Alums)
	http.Handle("/", r)
	r.PathPrefix("/").Handler(http.StripPrefix("/",
		http.FileServer(http.Dir("templates/"))))

	defer db.Close()
	err := http.ListenAndServe(CONN_HOST+":"+CONN_PORT, nil)
	if err != nil {
		log.Fatal("error starting http server : ", r)
		return
	}
}
