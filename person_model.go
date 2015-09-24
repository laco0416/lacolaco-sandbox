package app

import (
	"github.com/mjibson/goon"

	"appengine"
)

type Person struct {
	ID   int64  `goon:"id" datastore:"-" json:"-"`
	Age  int64  `json:"age"`
	Name string `datastore:",noindex" json:"name"`
	Sex  string `json:"sex"`
}

func (p *Person) Save(c appengine.Context) error {
	g := goon.FromContext(c)
	if _, err := g.Put(p); err != nil {
		return err
	}
	return nil
}
