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

func handleConfigHome(cs *ConfigServer) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		b := new(bytes.Buffer)
		err := cs.homeT.Execute(b, cs.config)
		Check(err)
		_, err = rw.Write(b.Bytes())
		Check(err)
	}
}

func handleGetRoute(cs *ConfigServer) func(rw http.ResponseWriter, r *http.Request) {
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
func handleUpdateRoute(cs *ConfigServer) func(rw http.ResponseWriter, r *http.Request) {
	return func(rw http.ResponseWriter, r *http.Request) {
		path := r.PostFormValue("path")
		statusVal := r.PostFormValue("status")
		method := r.PostFormValue("method")
		body := r.PostFormValue("body")
		action := r.PostFormValue("action")

		status, err := strconv.Atoi(statusVal)
		if err != nil {
			return
		}

		b := new(bytes.Buffer)
		if body == "" || path == "" || method == "" && action != "delete" {
			cs.homeT.Execute(b, cs.config)
			rw.Write(b.Bytes())
			return
		}

		key := path + ":" + method
		switch action {
		case "submit":
			cs.config.Routes[key] = Item{
				Method: method,
				Path:   path,
				Response: Response{
					Status: status,
					Body:   body,
				},
			}
		case "delete":
			delete(cs.config.Routes, key)
		}

		WriteConfig(cs.config)
		cs.homeT.Execute(b, cs.config)
		rw.Write(b.Bytes())
	}
}

func (m ConfigServer) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	key := r.URL.String() + ":" + r.Method
	val, ok := m.config.Routes[key]
	if ok == false {
		rw.WriteHeader(http.StatusNotFound)
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
	mux.HandleFunc("/config", handleConfigHome(cs))
	mux.HandleFunc("/getFromPath", handleGetRoute(cs))
	mux.HandleFunc("/updateItem", handleUpdateRoute(cs))

	configServerUrl := "http://localhost:8081/config"
	mockServerUrl := "http://localhost:8080"

	fmt.Println("Config server :", configServerUrl)
	fmt.Println("Mock server URL:", mockServerUrl)

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
