package app

import (
	"github.com/mjibson/goon"

	"appengine"
	"appengine/datastore"
)

type Child struct {
	ID     string         `goon:"id" datastore:"-" json:"-"`
	Parent *datastore.Key `goon:"parent" datastore:"-"`
	Text   string
}

func (src *Child) Save(c appengine.Context, p *Person) error {
	g := goon.FromContext(c)
	src.Parent = g.Key(p)

	if _, err := g.Put(src); err != nil {
		return err
	}
	return nil
}
