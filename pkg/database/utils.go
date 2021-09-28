package database

import (
	"database/sql"
	"encoding/csv"
	"fmt"

	"github.com/spf13/cast"
)

func ReadRowsToCsvFile(rows *sql.Rows, w *csv.Writer, isPrintHead bool) error {
	if rows == nil {
		return fmt.Errorf("rows can't be nil")
	}
	columns, err := rows.Columns()
	if err != nil {
		return err
	}
	length := len(columns)
	if isPrintHead {
		// print csv file headers
		if err = w.Write(columns); err != nil {
			return err
		}
	}
	//unnecessary to put below into rows.Next loop,reduce allocating
	values := make([]interface{}, length)
	for i := 0; i < length; i++ {
		values[i] = new(interface{})
	}

	record := make([]string, length)
	for rows.Next() {
		err = rows.Scan(values...)
		if err != nil {
			return err
		}

		for i := 0; i < len(columns); i++ {
			value := *(values[i].(*interface{}))
			record[i] = cast.ToString(value)
		}

		if err = w.Write(record); err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}
