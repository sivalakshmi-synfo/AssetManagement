package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"text/template"

	_ "github.com/go-sql-driver/mysql"
)

type Employee struct {
	IdEmployees int
	Name        string
	DOB         string
	Email       string
	Mobile      string
	Address     string
}

func dbConn() (db *sql.DB) {
	dbDriver := "mysql"
	dbUser := "root"
	dbPass := "admin@123"
	dbName := "AMS"
	db, err := sql.Open(dbDriver, dbUser+":"+dbPass+"@/"+dbName)
	if err != nil {
		panic(err.Error())
	}
	return db
}

var tmpl = template.Must(template.ParseGlob("view/*"))

func employees(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")

	db := dbConn()

	selDB, err := db.Query("SELECT * FROM Employees")
	if err != nil {
		panic(err.Error())
	}
	emp := Employee{}
	res := []Employee{}
	for selDB.Next() {
		err = selDB.Scan(&emp.IdEmployees, &emp.Name, &emp.DOB, &emp.Email, &emp.Mobile, &emp.Address)
		if err != nil {
			panic(err.Error())
		}
		res = append(res, emp)
	}
	tmpl.ExecuteTemplate(w, "list_mployees", res)

	//tmpl.ExecuteTemplate(w, "Index", res) // defined name not file name
	defer db.Close()
}
func show(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	id := r.URL.Query().Get("id")
	db := dbConn()

	selDB := db.QueryRow("SELECT * FROM Employees WHERE IdEmployees=?", id)

	emp := Employee{}
	err := selDB.Scan(&emp.IdEmployees, &emp.Name, &emp.DOB, &emp.Email, &emp.Mobile, &emp.Address)
	if err != nil {
		panic(err.Error())
	}

	if mode := r.URL.Query().Get("mode"); mode == "view" {
		// t, _ := template.ParseFiles("view/show.html")
		// t.Execute(w, emp)
		tmpl.ExecuteTemplate(w, "showView", emp)

	} else if mode == "edit" {
		// t, _ := template.ParseFiles("view/edit.html")
		// t.Execute(w, emp)
		tmpl.ExecuteTemplate(w, "edit", emp)
	}

	//tmpl.ExecuteTemplate(w, "Index", res) // defined name not file name
	defer db.Close()
}
func update(w http.ResponseWriter, r *http.Request) {

	db := dbConn()
	if err := r.ParseForm(); err != nil {
		fmt.Fprintf(w, "ParseForm() err: %v", err)
		return
	}

	if r.Method == "POST" {
		IDEmployees := r.FormValue("IdEmployees")
		Name := r.FormValue("Name")
		DOB := r.FormValue("DOB")

		Email := r.FormValue("Email")
		Mobile := r.FormValue("Mobile")
		Address := r.FormValue("Address")
		update, err := db.Prepare("UPDATE Employees SET Name=?, DOB=?,Email=?, Mobile=?,Address=? WHERE IdEmployees=?")
		if err != nil {
			panic(err.Error())
		}
		update.Exec(Name, DOB, Email, Mobile, Address, IDEmployees)

	}
	http.Redirect(w, r, "/employees", http.StatusSeeOther)
	defer db.Close()
}

func adding(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	db := dbConn()

	if r.Method == "GET" {
		//t, _ := template.ParseFiles("view/new.html")

		//	t.Execute(w, nil)
		tmpl.ExecuteTemplate(w, "addemployee", nil)

	}
	if r.Method == "POST" {
		Name := r.FormValue("Name")
		DOB := r.FormValue("DOB")

		Email := r.FormValue("Email")
		Mobile := r.FormValue("Mobile")
		Address := r.FormValue("Address")
		ins, err := db.Prepare("insert into Employees ( Name, DOB,Email, Mobile,Address) values (?,?,?,?,?) ")
		if err != nil {
			panic(err.Error())
		}
		ins.Exec(Name, DOB, Email, Mobile, Address)
		http.Redirect(w, r, "/employees", 301)
	}

	defer db.Close()

}
func delete(w http.ResponseWriter, r *http.Request) {
	db := dbConn()
	id := r.URL.Query().Get("id")
	delte, err := db.Prepare("delete from Employees  WHERE IdEmployees=?")
	if err != nil {
		panic(err.Error())
	}
	delte.Exec(id)
	defer db.Close()
	http.Redirect(w, r, "/employees", http.StatusSeeOther)
}
func dashboard(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	if r.Method == "GET" {
		db := dbConn()
		selDB := db.QueryRow("SELECT count(*) as count FROM Employees")
		var count int
		err := selDB.Scan(&count)
		if err != nil {
			panic(err.Error())
		}
		tmpl.ExecuteTemplate(w, "dashboard", count)
		// t.Execute(w, count)
		defer db.Close()
	} else if r.Method == "POST" {
		http.Redirect(w, r, "/employees", http.StatusSeeOther)
	}

}
func main() {
	log.Println("Server started on: http://localhost:9000")
	// mux:=http.NewServeMux();
	// mux.HandleFunc("/add",adding)
	http.HandleFunc("/add", adding)
	http.HandleFunc("/employees", employees)
	http.HandleFunc("/", dashboard)
	http.HandleFunc("/show", show)
	http.HandleFunc("/edit", show)
	http.HandleFunc("/delete", delete)
	http.HandleFunc("/update", update)

	http.ListenAndServe(":9000", nil)
}
