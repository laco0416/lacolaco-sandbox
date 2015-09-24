package app

import (
	"net/http"
	"strconv"

	"encoding/json"
	"time"

	"github.com/laco0416/bbq"
)

func aespySetup() {
	http.HandleFunc("/aespy-sample", handleAeSpySample)
}

func handleAeSpySample(w http.ResponseWriter, r *http.Request) {
	b := bbq.NewBBQ(&bbq.Option{Log: true})
	b.AddKind("Person", "aespy", "person")
	b.AddKind("Child", "aespy", "child")
	c, ch := b.Hook(r)
	p := &Person{}
	age, err := strconv.ParseInt(r.URL.Query().Get("age"), 10, 64)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	p.Age = age
	p.Name = r.URL.Query().Get("name")
	p.Sex = r.URL.Query().Get("sex")
	if err := p.Save(c); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	child := &Child{ID: time.Now().Format("20060102_030405"), Text: "Child"}
	if err := child.Save(c, p); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	resp, err := json.Marshal(child)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if err := <-ch; err != nil {
		c.Debugf(err.Error())
	}
	if err := <-ch; err != nil {
		c.Debugf(err.Error())
	}
	w.Write(resp)
}
