package app

import (
	"net/http"
	"strconv"

	"github.com/dustin/gojson"
	"github.com/laco0416/aespy"
	ae "google.golang.org/appengine"
)

func aespySetup() {
	http.HandleFunc("/aespy-sample", handleAeSpySample)
}

func handleAeSpySample(w http.ResponseWriter, r *http.Request) {
	c := aespy.NewContext(r)
	ctx := ae.NewContext(r)

	personHandler := &PersonDatastoreHandler{ctx}
	c.AddDatastoreHandler(personHandler)

	p := &Person{
		Name: r.URL.Query().Get("name"),
	}
	age, err := strconv.ParseInt(r.URL.Query().Get("age"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p.Age = age
	p.Save(c)
	resp, err := json.Marshal(p)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Write(resp)
}
