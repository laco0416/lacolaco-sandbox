package app

import (
	"reflect"
	"time"

	"github.com/laco0416/aespy"
	"golang.org/x/net/context"
	"google.golang.org/api/bigquery/v2"

	"appengine"
	"appengine/datastore"
)

type BbqHandler struct {
	Context     context.Context
	BqDataSetID string
	BqTableID   string
}

func (h BbqHandler) InsertToBq(c appengine.Context, key *datastore.Key, entity *aespy.Entity) error {
	fields := make([]*bigquery.TableFieldSchema, 0, len(entity.Properties))
	jsonRow := make(map[string]bigquery.JsonValue)
	fields = append(fields, &bigquery.TableFieldSchema{
		Name: "__ID__",
		Type: "INTEGER",
	})
	fields = append(fields, &bigquery.TableFieldSchema{
		Name: "__Name__",
		Type: "STRING",
	})
	jsonRow["__ID__"] = key.IntID()
	jsonRow["__Name__"] = key.StringID()
	for key, value := range entity.Properties {
		var typeStr string

		switch reflect.ValueOf(value).Kind() {
		case reflect.String:
			typeStr = "STRING"
		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			typeStr = "INTEGER"
		case reflect.Float32, reflect.Float64:
			typeStr = "FLOAT"
		case reflect.Bool:
			typeStr = "BOOLEAN"
		case reflect.TypeOf(time.Time{}).Kind():
			typeStr = "TIMESTAMP"
		default:
			typeStr = "RECORD"
		}
		fields = append(fields, &bigquery.TableFieldSchema{
			Name: key,
			Type: typeStr,
		})
		jsonRow[key] = value
	}
	err := createTable(c, h.Context, h.BqDataSetID, h.BqTableID, &bigquery.TableSchema{
		Fields: fields,
	})
	if err != nil {
		return err
	}
	err = insertToBq(c, h.Context, h.BqDataSetID, h.BqTableID, jsonRow)
	if err != nil {
		return err
	}
	return nil
}
