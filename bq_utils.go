package app

import (
	"net/http"
	"time"

	"golang.org/x/net/context"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"google.golang.org/api/bigquery/v2"
	"google.golang.org/api/googleapi"

	"appengine"
	"appengine/urlfetch"
)

func createTable(c appengine.Context, ctx context.Context, dataSetID, tableID string, schema *bigquery.TableSchema) error {
	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(ctx, bigquery.BigqueryScope),
			Base:   &urlfetch.Transport{Context: c},
		},
	}

	bqs, err := bigquery.New(client)
	if err != nil {
		return err
	}

	_, err = bqs.Tables.Get(appengine.AppID(c), dataSetID, tableID).Do()
	if err == nil {
		// 作成済み
		c.Debugf("%s is already exists", tableID)
		return nil
	} else if gerr, ok := err.(*googleapi.Error); ok && gerr.Code != http.StatusNotFound {
		// 404エラー以外の場合は想定外の何かなので処理やめる
		c.Warningf("%#v", err)
		return err
	}

	c.Debugf("insert %s", tableID)
	table := &bigquery.Table{
		TableReference: &bigquery.TableReference{
			ProjectId: appengine.AppID(c),
			DatasetId: dataSetID,
			TableId:   tableID,
		},
		Description:  "",
		FriendlyName: "",
		Schema:       schema,
	}
	_, err = bqs.Tables.Insert(appengine.AppID(c), dataSetID, table).Do()
	if err != nil {
		return err
	}

	return nil
}

func insertToBq(c appengine.Context, ctx context.Context, dataSetID, tableID string, jsonRow map[string]bigquery.JsonValue) error {
	client := &http.Client{
		Transport: &oauth2.Transport{
			Source: google.AppEngineTokenSource(ctx, bigquery.BigqueryInsertdataScope),
			Base:   &urlfetch.Transport{Context: c},
		},
		Timeout: 1500 * time.Millisecond,
	}
	bqs, err := bigquery.New(client)
	if err != nil {
		return err
	}
	row := &bigquery.TableDataInsertAllRequestRows{
		Json: jsonRow,
	}
	data := &bigquery.TableDataInsertAllRequest{Rows: []*bigquery.TableDataInsertAllRequestRows{row}}

	var resp *bigquery.TableDataInsertAllResponse

	// Insertに3回失敗した場合はエラーとする。2回までリトライを行う。
	var i int
	for i = 0; i < 3; i++ {
		resp, err = bqs.Tabledata.InsertAll(appengine.AppID(c), dataSetID, tableID, data).Do()
		if err != nil {
			c.Debugf("retry insert. err: %#v", err)
		} else {
			break
		}
	}
	if err != nil {
		return err
	}
	if i != 0 {
		c.Debugf("recover insert error")
	}

	for _, err1 := range resp.InsertErrors {
		for _, err2 := range err1.Errors {
			c.Warningf("%#v", err2)
		}
	}
	c.Debugf("bq insertion done.")

	return nil
}
