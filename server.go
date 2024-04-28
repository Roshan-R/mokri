package main

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ConfigServer struct {
	config   *Config
	homeT    *template.Template
	detailsT *template.Template
}

func homePage(cs *ConfigServer) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		b := new(bytes.Buffer)
		err := cs.homeT.Execute(b, cs.config)
		Check(err)
		_, err = rw.Write(b.Bytes())
		Check(err)
	}
}

func getRoute(cs *ConfigServer) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		tn := r.Header.Get("Hx-Trigger-Name")
		b := new(bytes.Buffer)
		if tn == "special@reset" {
			err := cs.detailsT.Execute(b, nil)
			Check(err)
			_, err = rw.Write(b.Bytes())
			Check(err)
			return
		}
		item, ok := cs.config.Routes[tn]
		if !ok {
			_, err := rw.Write([]byte("Route not found"))
			Check(err)
		}
		cs.detailsT.Execute(b, item)
		_, err := rw.Write(b.Bytes())
		Check(err)
		return
	}
}
func updateRoute(cs *ConfigServer) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {

		path := r.PostFormValue("path")
		statusVal := r.PostFormValue("status")
		method := r.PostFormValue("method")
		body := r.PostFormValue("body")

		status, err := strconv.Atoi(statusVal)
		if err != nil {
			return
		}

		action := r.PostFormValue("action")
		b := new(bytes.Buffer)

		if body == "" || path == "" || method == "" && action != "delete" {
			cs.homeT.Execute(b, cs.config)
			rw.Write(b.Bytes())
			return
		}

		key := path + ":" + method
		if action == "submit" {
			cs.config.Routes[key] = Item{
				Method: method,
				Path:   path,
				Response: Response{
					Status: status,
					Body:   body,
				},
			}
		} else if action == "delete" {
			fmt.Println(cs.config)
			delete(cs.config.Routes, key)
			fmt.Println(cs.config)
		}

		WriteConfig(cs.config)
		cs.homeT.Execute(b, cs.config)
		rw.Write(b.Bytes())
	}
}

func (m ConfigServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	key := r.URL.String() + ":" + r.Method
	log.Println(key)
	val, ok := m.config.Routes[key]
	if ok == false {
		rw.Write([]byte("Route not found"))
		return
	}
	rw.WriteHeader(val.Response.Status)
	rw.Write([]byte(val.Response.Body))
	return

}
func main() {

	c, _ := ReadConfig()
	ht := template.Must(template.ParseFiles("templates/page.html"))
	dt := template.Must(template.ParseFiles("templates/details.html"))

	cs := &ConfigServer{
		config:   &c,
		homeT:    ht,
		detailsT: dt,
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/config", homePage(cs))
	mux.HandleFunc("/getFromPath", getRoute(cs))
	mux.HandleFunc("/updateItem", updateRoute(cs))

	s := &http.Server{
		Addr:           ":8080",
		Handler:        cs,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}
	go func() { http.ListenAndServe(":8081", mux) }()
	log.Fatal(s.ListenAndServe())

}
