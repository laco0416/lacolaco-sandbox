package app

import (
    "net/http"
    "strconv"

    "encoding/json"
    "github.com/laco0416/bbq"
)

func aespySetup() {
    http.HandleFunc("/aespy-sample", handleAeSpySample)
}

func handleAeSpySample(w http.ResponseWriter, r *http.Request) {
    b := bbq.NewBBQ(&bbq.Option{Log: true})
    b.AddKind("Person", "aespy", "person")
    c, _ := b.Hook(r)
    p := &Person{}
    age, err := strconv.ParseInt(r.URL.Query().Get("age"), 10, 64)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    p.Age = age
    p.Name = r.URL.Query().Get("name")
    p.Sex = r.URL.Query().Get("sex")
    p.Bar = "Foo"
    if err := p.Save(c); err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    resp, err := json.Marshal(p)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
//    <-ch
    w.Write(resp)
}
