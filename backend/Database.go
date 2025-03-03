package backend

import (
	"database/sql"
	"fmt"
	"strconv"

	_ "github.com/mattn/go-sqlite3"
)


func InitDB() {
	databases, _ := sql.Open("sqlite3", "./forum.db")
	statement, _ := databases.Prepare("CREATE TABLE IF NOT EXISTS user (id INTEGER PRIMARY KEY, firstname TEXT, lastname TEXT)")
	statement.Exec()
	statement, _ = databases.Prepare("INSERT INTO user (firstname, lastname) VALUES (?, ?)")
	statement.Exec("theo", "loic")
	rows, _ := databases.Query("SELECT id, firstname, lastname FROM user")
	var id int
	var firstname string
	var lastname string
	for rows.Next(){
		rows.Scan(&id, &firstname, &lastname)
		fmt.Println(strconv.Itoa(id) + ": " + firstname + " " + lastname  )
	}

}
