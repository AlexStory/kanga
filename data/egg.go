package data

import (
	"database/sql"
	"fmt"
)

func UndoEgg(db *sql.DB) {
	stmt := `
	DELETE FROM exeggutor
	WHERE id = (SELECT MAX(id) FROM exeggutor)
	`
	_, err := db.Exec(stmt)
	if err != nil {
		fmt.Printf("Failed to undo flip: %v\n", err)
	}
}
