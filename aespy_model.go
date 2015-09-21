package app

import (
	"github.com/laco0416/aespy"
	"github.com/mjibson/goon"
	"golang.org/x/net/context"
	"google.golang.org/api/bigquery/v2"

	"appengine"
	"appengine/datastore"
)

type Person struct {
	ID   int64  `goon:"id" datastore:"-" json:"-"`
	Age  int64  `json:"id"`
	Name string `datastore:",noindex" json:"name"`
}

func (l *Person) Save(c appengine.Context) error {
	g := goon.FromContext(c)
	if _, err := g.Put(l); err != nil {
		return err
	}
	return nil
}

func (p *Person) InsertToBq(c appengine.Context, ctx context.Context) error {
	err := createTable(c, ctx, "aespy", "Person", &bigquery.TableSchema{
		Fields: []*bigquery.TableFieldSchema{
			&bigquery.TableFieldSchema{
				Name: "id",
				Type: "INTEGER",
			},
			&bigquery.TableFieldSchema{
				Name: "age",
				Type: "INTEGER",
			},
			&bigquery.TableFieldSchema{
				Name: "name",
				Type: "STRING",
			},
		},
	})
	if err != nil {
		return err
	}
	jsonRow := make(map[string]bigquery.JsonValue)
	jsonRow["id"] = p.ID
	jsonRow["age"] = p.Age
	jsonRow["name"] = p.Name
	err = insertToBq(c, ctx, "aespy", "Person", jsonRow)
	if err != nil {
		return err
	}
	return nil
}

// PersonDatastoreHandler はPersonについてDatastoreHandlerを実装する
type PersonDatastoreHandler struct {
	ctx context.Context
}

func (h *PersonDatastoreHandler) PrePut(c appengine.Context, entity *aespy.Entity) error {
	return nil
}

func (h *PersonDatastoreHandler) PostPut(c appengine.Context, key *datastore.Key) error {
	c.Debugf("PostPut")
	g := goon.FromContext(c)
	p := &Person{ID: key.IntID()}
	if err := g.Get(p); err != nil {
		return err
	}
	if err := p.InsertToBq(c, h.ctx); err != nil {
		return err
	}
	return nil
}

func (h *PersonDatastoreHandler) PreGet(c appengine.Context, key *datastore.Key) error {
	return nil
}

func (h *PersonDatastoreHandler) PostGet(c appengine.Context, entity *aespy.Entity) error {
	return nil
}

func (h *PersonDatastoreHandler) PreDelete(c appengine.Context, key *datastore.Key) error {
	return nil
}

func (h *PersonDatastoreHandler) PostDelete(c appengine.Context) error {
	return nil
}
