package utils

import (
	"database/sql"
	"encoding/json"
	"reflect"
)

// queries every table to json
func QueryToJson(db *sql.DB, query string, args ...interface{}) ([]byte, error) {
	var objects []map[string]interface{}

	rows, err := db.Query(query, args...)
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		// get column types
		columns, err := rows.ColumnTypes()
		if err != nil {
			return nil, err
		}

		// Create row statics
		values := make([]interface{}, len(columns))
		object := map[string]interface{}{}

		// fetch row into object
		for i, column := range columns {
			object[column.Name()] = reflect.New(column.ScanType()).Interface()
			values[i] = object[column.Name()]
		}

		err = rows.Scan(values...)
		if err != nil {
			return nil, err
		}

		objects = append(objects, object)
	}

	defer rows.Close()
	return json.MarshalIndent(objects, "", "\t")
}
