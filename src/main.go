package main

import (
	"html/template"
	"log"
	"net/http"
)

//deklarasi template
var tpl *template.Template

//deklarasi data variable
type pageData struct {
	Judul string
	Iklan string
	Data1 string
}

func init() {
	tpl = template.Must(template.ParseGlob("./template/*.gohtml"))
}

//directory static
func main() {
	http.Handle("/public/", http.StripPrefix("/public", http.FileServer(http.Dir("./public"))))
	http.HandleFunc("/", idx)
	http.HandleFunc("/about", abt)
	log.Println("Server Jalan di port 3000...!!!")
	http.ListenAndServe(":3000", nil)
}

func idx(w http.ResponseWriter, req *http.Request) {
	pd := pageData{Judul: "coba-coba", Iklan: "halaman ini di persembahkan oleh.......", Data1: "test aja"}
	err := tpl.ExecuteTemplate(w, "index.gohtml", pd)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
}
func abt(w http.ResponseWriter, req *http.Request) {
	pd := pageData{Judul: "halaman about"}
	err := tpl.ExecuteTemplate(w, "about.gohtml", pd)
	if err != nil {
		log.Println(err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
	}
}
