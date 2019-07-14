package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/kataras/iris"
)

func main() {
	db, err := sql.Open("mysql", "root:root@tcp(localhost:3306)/yami?charset=utf8")
	if err != nil {
		panic(err.Error())
	}
	defer db.Close()

	app := iris.Default()
	// >>> api >>>
	app.Post("/table2json", func(ctx iris.Context) {
		sqlString := ctx.FormValue("s")
		fmt.Println("sqlString[" + sqlString + "]")
		result, err := getJSON(db, sqlString)
		if err != nil {
			ctx.JSON(err.Error())
			panic(err.Error())
		} else {
			ctx.Write([]byte(result))
		}
	})
	// <<< api <<<
	fmt.Print(db)
	fmt.Println("main finish.")
	app.Run(iris.Addr(":7080"))
	os.Exit(0)
}

func getJSON(db *sql.DB, sqlString string) (string, error) {
	fmt.Print("db=")
	fmt.Println(db)
	rows, err := db.Query(sqlString)
	if err != nil {
		return "", err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i]
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i]
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return "", err
	}
	fmt.Println(string(jsonData))
	return string(jsonData), nil
}
