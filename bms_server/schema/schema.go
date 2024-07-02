package schema

import (
	"crypto/rand"
	"database/sql"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	jsoniter "github.com/json-iterator/go"
	"github.com/ronnie-w/bms-server/database"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	db   = database.Conn()
)

func GeneralDecoder(r *http.Request) Args {
	r.ParseForm()

	var form map[string][]string
	for key := range r.PostForm {
		if err := json.Unmarshal([]byte(key), &form); err != nil {
			fmt.Println(err)
		}
	}

	var (
		args  Args
		args2 []Arg2
		arg1  = Arg1{Arg1Val: form[`Args1[Arg1Val]`][0]}
		arg2  Arg2
	)

	for _, val := range form[`Args2`] {
		strArr := strings.Split(val, " ")
		_, arg2Val1, _ := strings.Cut(strArr[0], "{")
		arg2Val2, _, _ := strings.Cut(strArr[1], "}")

		arg2 = Arg2{Arg2Val1: arg2Val1, Arg2Val2: arg2Val2}

		args2 = append(args2, arg2)
	}

	args = Args{Args1: arg1, Args2: args2}

	return args
}

func QueryParser(i Args) string {
	var queries []string
	if len(i.Args2) < 1 {
		return "null"
	}

	for _, col := range i.Args2 {
		if col.Arg2Val2 == "number" {
			queries = append(queries, col.Arg2Val1+" int not null")
		} else {
			queries = append(queries, col.Arg2Val1+" varchar(255) not null")
		}
	}

	return strings.Join(queries, ",")
}

func CreateTableQuery(table string, columns []string) error {
	createQuery := `CREATE TABLE IF NOT EXISTS ` + table + ` (
				uuid varchar(36) not null,
				id int not null,
				date_created varchar(19) not null,
				date_modified varchar(19) not null,
				` + strings.Replace(strings.Join(columns, ","), ",null", "", -1) + `
	)`

	_, err := db.Exec(createQuery)

	return err
}

type PragmaStruct struct {
	Name  string
	Ctype string
}

func Pragma(table string) []PragmaStruct {
	var (
		cid, name, ctype, notNull, dfltValue, pk string
		PragmaValues                             []PragmaStruct
	)

	tableRows, err := db.Query("PRAGMA table_info(" + table + ")")
	defer tableRows.Close()
	if err != nil {
		fmt.Println(err)
	}

	for tableRows.Next() {
		tableRows.Scan(&cid, &name, &ctype, &notNull, &dfltValue, &pk)
		PragmaValues = append(PragmaValues, PragmaStruct{name, ctype})
	}

	return PragmaValues
}

type InsertRes struct {
	Err  interface{}
	Uuid string
}

func InsertQuery(table string, values []interface{}) InsertRes {
	var (
		columnNames []string

		uuid         = uuid.New().String()
		id, _        = rand.Prime(rand.Reader, 20)
		timestamp    = time.Now().Format(time.DateTime)
		parsedValues = []interface{}{uuid, int(id.Uint64()), timestamp, timestamp}
	)

	for _, col := range Pragma(table) {
		columnNames = append(columnNames, col.Name)
	}

	for i := 0; i < len(values); i++ {
		parsedValues = append(parsedValues, values[i])
	}

	// whitelist:
	// employee_time_table
	// task_assignment
	// sales notifications

	if strings.Contains("salesnotifications", table) || !database.Exists(table, columnNames[4], fmt.Sprintf("%v", parsedValues[4])) {
		db.Exec(`INSERT INTO `+table+`(`+strings.Join(columnNames, ",")+`)
						values(`+strings.Join(strings.Split(strings.Repeat("?", len(columnNames)), ""), ",")+`)`, parsedValues...)
	} else {
		return InsertRes{fmt.Sprintf("%v", parsedValues[4]) + " already exists", ""}
	}

	return InsertRes{nil, uuid}
}

func SelectRowsScan(r *http.Request, rows *sql.Rows) [][]QueryData {
	var (
		selectResults [][]QueryData
		table         = GeneralDecoder(r).Args1.Arg1Val
		cols          = Pragma(table)
		results       = make([]interface{}, len(cols))
	)

	for i := range results {
		results[i] = new(interface{})
	}
	defer rows.Close()

	for rows.Next() {
		var parsedResults []QueryData
		if err := rows.Scan(results[:]...); err != nil {
			fmt.Println(err)
		}

		for i, col := range cols {
			var str string
			val := *results[i].(*interface{})

			if val == nil {
				str = "NULL"
			} else {
				switch v := val.(type) {
				case []byte:
					str = string(v)
				default:
					str = fmt.Sprintf("%v", v)
				}
			}

			parsedResults = append(parsedResults, QueryData{col.Name, str})
		}

		selectResults = append(selectResults, parsedResults)
	}

	return selectResults
}

func UpdateQuery(uuid string, updateData Args) interface{} {
	var (
		warning    string
		table      = updateData.Args1.Arg1Val
		updateCols = []string{"date_modified"}
		updateVals = []interface{}{time.Now().Format(time.DateTime)}
	)

	for _, update := range updateData.Args2 {
		if database.Exists(table, update.Arg2Val1, update.Arg2Val2) {
			warning += `
			` + update.Arg2Val1 + `: ` + update.Arg2Val2 + ` has a duplicate value in table: ` + table + `. 
			Avoid duplicate unique values in records! Ignore this warning if the updated values are not in unique columns.
			`
		}

		updateCols = append(updateCols, update.Arg2Val1)
		updateVals = append(updateVals, update.Arg2Val2)
	}

	for i := 0; i < len(updateCols); i++ {
		updateCols[i] = updateCols[i] + "=?"
	}
	updateVals = append(updateVals, uuid)

	_, err := db.Exec("UPDATE "+table+" SET "+strings.Join(updateCols, ",")+" WHERE uuid=?", updateVals...)
	if len(warning) > 0 {
		return warning
	}

	return err
}

func DeleteQuery(uuid string, deleteData Args) error {
	_, err := db.Exec("DELETE FROM "+deleteData.Args1.Arg1Val+" WHERE uuid=?", uuid)

	return err
}

func DevSelect(cols []string, table, clauseId, clause string, scanVar []interface{}) error {
	selectRow := db.QueryRow("SELECT "+strings.Join(cols, ",")+" FROM "+table+" WHERE "+clauseId+"=?", clause)
	return selectRow.Scan(scanVar...)
}

func DevUpdate(table, clauseId string, cols []string, args []interface{}) error {
	for i := 0; i < len(cols); i++ {
		cols[i] = cols[i] + "=?"
	}

	_, err := db.Exec("UPDATE "+table+" SET "+strings.Join(cols, ",")+" WHERE "+clauseId+"=?", args...)
	return err
}

// server success schema
type Success struct {
	OK string
}

// server error schema
type Error struct {
	Err string
}

// general schema
type Arg2 struct {
	Arg2Val1 string
	Arg2Val2 string
}

type Arg1 struct {
	Arg1Val string
}

type Args struct {
	Args1 Arg1
	Args2 []Arg2
}

// query data schema
type QueryData struct {
	ColName  string
	ColValue string
}
