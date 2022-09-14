package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"

	_ "github.com/lib/pq"
)

var (
	DatabaseURL = "postgres://postgres:1234@db:5432/restapi_dev?sslmode=disable"
)

func NewDbExplorer(db *sql.DB) (http.Handler, error) {

	srv := http.NewServeMux()
	srv.Handle("/", handlerMain(db))

	return srv, nil
}
func handlerMain(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		table := strings.Split(strings.Split(r.URL.Path, "/")[1], "?")[0]
		result, code, err := readAll(db, table)

		if err != nil {
			Error(w, r, code, err)
			return
		}
		Respond(w, r, code, result)

	}
}

func Error(w http.ResponseWriter, r *http.Request, code int, err error) {
	Respond(w, r, code, map[string]string{"error": err.Error()})

}
func Respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}

func readAll(db *sql.DB, tableName string) ([]map[string]interface{}, int, error) {
	rows, err := db.Query("SELECT column_name, data_type FROM information_schema.columns WHERE table_name = $1",
		tableName)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer rows.Close()

	queueName := []string{}
	queueType := []string{}
	queuePtr := []interface{}{}

	for rows.Next() {
		var colName string
		var colType string
		err = rows.Scan(&colName, &colType)
		if err != nil {
			return nil, http.StatusInternalServerError, err
		}
		queueName = append(queueName, colName)
		if strings.Contains(colType, "int") {
			queuePtr = append(queuePtr, new(int))
			queueType = append(queueType, "int")
		} else {
			queuePtr = append(queuePtr, new(string))
			queueType = append(queueType, "string")
		}
	}

	result := []map[string]interface{}{}
	query := fmt.Sprintf("SELECT * FROM %s\n", tableName)
	rows, err = db.Query(query)
	if err != nil {
		return nil, http.StatusBadRequest, errors.New("table does not exist")
	}
	defer rows.Close()
	for rows.Next() {
		err = rows.Scan(queuePtr...)
		if err != nil && !strings.Contains(err.Error(), "NULL") {
			return nil, http.StatusInternalServerError, err
		}

		record := map[string]interface{}{}
		for i, v := range queueName {
			if queueType[i] == "int" {
				record[v] = *queuePtr[i].(*int)
				queuePtr[i] = new(int)
			} else {
				record[v] = *queuePtr[i].(*string)
				queuePtr[i] = new(string)
			}
		}
		result = append(result, record)
	}
	return result, http.StatusOK, nil
}

func main() {
	db, err := sql.Open("postgres", DatabaseURL)
	if err != nil {
		panic(err)
	}

	if err != nil {
		panic(err)
	}

	handler, err := NewDbExplorer(db)
	if err != nil {
		panic(err)
	}

	fmt.Println("starting server at :8080")
	http.ListenAndServe(":8080", handler)
}
