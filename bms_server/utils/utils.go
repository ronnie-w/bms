package utils

import (
	"fmt"
	"net/http"

	jsoniter "github.com/json-iterator/go"
	"github.com/ronnie-w/bms-server/database"
	"github.com/ronnie-w/bms-server/schema"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	db   = database.Conn()
)

func GetColumns(rw http.ResponseWriter, r *http.Request) {
	json.NewEncoder(rw).Encode(schema.Pragma(schema.GeneralDecoder(r).Args1.Arg1Val))
}

func DefaultInsert(r *http.Request) (string, schema.InsertRes) {
	var (
		values []interface{}
		args   = schema.GeneralDecoder(r)
		table  = args.Args1.Arg1Val
	)

	for _, a := range args.Args2 {
		values = append(values, a.Arg2Val1)
	}

	return table, schema.InsertQuery(table, values)
}

func InsertInto(rw http.ResponseWriter, r *http.Request) {
	_, insertRes := DefaultInsert(r)
	if insertRes.Err != nil {
		ErrResponse(rw, nil, insertRes.Err)
	} else {
		SuccessResponse(rw, "Insert operation complete without errors")
	}
}

func DefaultSelect(rw http.ResponseWriter, r *http.Request) [][]schema.QueryData {
	table := schema.GeneralDecoder(r).Args1.Arg1Val
	rows, err := db.Query("SELECT * FROM " + table + "")
	if err != nil {
		fmt.Println(err)
		ErrResponse(rw, err, "Could not fetch data from table "+table+"")
	}

	return schema.SelectRowsScan(r, rows)
}

func SelectQuery(rw http.ResponseWriter, r *http.Request) {
	selectData := DefaultSelect(rw, r)
	json.NewEncoder(rw).Encode(selectData)
}

func UniqueSelectQuery(rw http.ResponseWriter, r *http.Request) {
	uniqueData := schema.GeneralDecoder(r)
	table := uniqueData.Args1.Arg1Val
	for _, uniqueVal := range uniqueData.Args2 {
		rows, err := db.Query("SELECT * FROM "+table+" WHERE "+uniqueVal.Arg2Val1+"=?", uniqueVal.Arg2Val2)
		if err != nil {
			fmt.Println(err)
			ErrResponse(rw, err, "Could not fetch data from table "+table+"")
		}

		json.NewEncoder(rw).Encode(schema.SelectRowsScan(r, rows))
	}
}

func LikeSelectQuery(rw http.ResponseWriter, r *http.Request) {
	uniqueData := schema.GeneralDecoder(r)
	table := uniqueData.Args1.Arg1Val
	for _, uniqueVal := range uniqueData.Args2 {
		rows, err := db.Query("SELECT * FROM " + table + " WHERE " + uniqueVal.Arg2Val1 + " LIKE '%" + uniqueVal.Arg2Val2 + "%'")
		if err != nil {
			fmt.Println(err)
			ErrResponse(rw, err, "Could not fetch data from table "+table+"")
		}

		json.NewEncoder(rw).Encode(schema.SelectRowsScan(r, rows))
	}
}

func DropTable(rw http.ResponseWriter, r *http.Request) {
	tableName := schema.GeneralDecoder(r).Args1.Arg1Val
	_, err := db.Exec("DROP TABLE " + tableName + "")
	if err != nil {
		ErrResponse(rw, err, "Could not remove table "+tableName+" from database")
	} else {
		SuccessResponse(rw, "Removed table "+tableName+" from database")
	}
}

func Exists(rw http.ResponseWriter, r *http.Request) {
	data := schema.GeneralDecoder(r)
	json.NewEncoder(rw).Encode(database.Exists(data.Args1.Arg1Val, data.Args2[0].Arg2Val1, data.Args2[0].Arg2Val2))
}

func ErrOrWarn(rw http.ResponseWriter, returnVal interface{}, successMsg, errMsg string) {
	if returnVal == nil {
		SuccessResponse(rw, successMsg)
	} else {
		errOrWarning := returnVal.(interface{})
		switch e := errOrWarning.(type) {
		case error:
			ErrResponse(rw, e, errMsg)
		default:
			ErrResponse(rw, nil, fmt.Sprintf("%v", e))
		}
	}
}

func ErrResponse(rw http.ResponseWriter, err error, errMsg interface{}) {
	if err != nil {
		json.NewEncoder(rw).Encode(schema.Error{Err: fmt.Sprintf("%v", errMsg) + ": " + err.Error()})
	} else {
		json.NewEncoder(rw).Encode(schema.Error{Err: fmt.Sprintf("%v", errMsg)})
	}
}

func SuccessResponse(rw http.ResponseWriter, successMsg string) {
	json.NewEncoder(rw).Encode(schema.Success{OK: successMsg})
}
