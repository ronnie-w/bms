package tabs

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ronnie-w/bms-server/schema"
	"github.com/ronnie-w/bms-server/utils"
)

func CreateNewTab(rw http.ResponseWriter, r *http.Request) {
	tabName := schema.GeneralDecoder(r).Args1.Arg1Val

	if err := schema.CreateTableQuery(tabName, []string{"record_name varchar(255)"}); err != nil {
		utils.ErrResponse(rw, err, "Could not create tab")
	} else {
		if insertErr := schema.InsertQuery("tabs", []interface{}{tabName}); insertErr.Err != nil {
			utils.ErrResponse(rw, nil, insertErr)
		} else {
			utils.SuccessResponse(rw, "Tab: "+tabName+" created successfully")
		}
	}
}

func CreateNewRecord(rw http.ResponseWriter, r *http.Request) {
	tabName := strings.Split(r.URL.Path, "/")[2]
	decodedForm := schema.GeneralDecoder(r)
	recordName := decodedForm.Args1.Arg1Val

	if schema.QueryParser(decodedForm) != "null" {
		if insertErr := schema.InsertQuery(tabName, []interface{}{recordName}); insertErr.Err != nil {
			utils.ErrResponse(rw, nil, insertErr)
		} else if err := schema.CreateTableQuery(recordName, []string{schema.QueryParser(decodedForm)}); err != nil {
			fmt.Println(err)
		} else {
			utils.SuccessResponse(rw, "New record: "+recordName+" created in tab: "+tabName+"")
		}
	} else {
		utils.ErrResponse(rw, nil, "Could not create record with no columns. Try adding new columns")
	}
}

func FetchTabs(rw http.ResponseWriter, r *http.Request) {
	utils.SelectQuery(rw, r)
}

func FetchRecordsInTab(rw http.ResponseWriter, r *http.Request) {
	utils.SelectQuery(rw, r)
}

func FetchRecordData(rw http.ResponseWriter, r *http.Request) {
	utils.SelectQuery(rw, r)
}

func InsertToRecord(rw http.ResponseWriter, r *http.Request) {
	utils.InsertInto(rw, r)
}

func UpdateInRecord(rw http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.Path, "/")[2]
	updateReturn := schema.UpdateQuery(uuid, schema.GeneralDecoder(r))

	utils.ErrOrWarn(rw, updateReturn, "Update operation complete without errors", "Could not update item "+uuid+" in record")
}

func DeleteFromRecord(rw http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.Path, "/")[2]
	err := schema.DeleteQuery(uuid, schema.GeneralDecoder(r))

	if err != nil {
		utils.ErrResponse(rw, err, "Could not delete item "+uuid+" from record")
	} else {
		utils.SuccessResponse(rw, "Delete operation complete without errors")
	}
}
