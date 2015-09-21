package app

import (
	"github.com/mjibson/goon"

	"appengine"
)

// +bbq
type Person struct {
	ID   int64  `goon:"id" datastore:"-" json:"-"`
	Age  int64  `json:"age"`
	Name string `datastore:",noindex" json:"name"`
}

func (l *Person) Save(c appengine.Context) error {
	g := goon.FromContext(c)
	if _, err := g.Put(l); err != nil {
		return err
	}
	return nil
}
