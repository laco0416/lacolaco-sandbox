package app

import (
	"github.com/laco0416/aespy"
	"golang.org/x/net/context"

	"appengine"
	"appengine/datastore"
)

type PersonDatastoreHandler struct {
	*BbqHandler
}

func NewPersonBbqHandler(ctx context.Context, bqDataSetID, bqTableID string) *PersonDatastoreHandler {
	bh := &BbqHandler{ctx, bqDataSetID, bqTableID}
	return &PersonDatastoreHandler{bh}
}

func (h *PersonDatastoreHandler) PostPut(c appengine.Context, key *datastore.Key, entity *aespy.Entity) error {
	if entity.Kind == "Person" {
		if err := h.InsertToBq(c, key, entity); err != nil {
			return err
		}
	}
	return nil
}
