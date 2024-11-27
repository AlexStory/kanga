package data

import (
	"database/sql"
	"fmt"
)

type MistyStats struct {
	TotalEntries int
	TotalHeads   int
}

func InsertMisty(db *sql.DB, heads int) {
	stmt := `
	INSERT INTO misty (heads, created_at)
	VALUES (?, CURRENT_TIMESTAMP)`

	_, err := db.Exec(stmt, heads)
	if err != nil {
		fmt.Printf("Failed to insert misty entry: %v\n", err)
	}
}

func GetMistyStats(db *sql.DB) (MistyStats, error) {
	var stats MistyStats
	err := db.QueryRow("SELECT COUNT(*), IFNULL(SUM(heads), 0) FROM misty").Scan(&stats.TotalEntries, &stats.TotalHeads)
	if err != nil {
		return MistyStats{}, fmt.Errorf("failed to get misty stats: %v", err)
	}
	return stats, nil
}

func UndoMisty(db *sql.DB) {
	stmt := `
	DELETE FROM misty
	WHERE id = (SELECT MAX(id) FROM misty)
	`
	_, err := db.Exec(stmt)
	if err != nil {
		fmt.Printf("Failed to undo misty entry: %v\n", err)
	}
}
