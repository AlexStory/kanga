package data

import (
	"database/sql"
	"fmt"
)

type EggType int

const (
	H EggType = iota
	HX
	T
	TX
)

type EggStats struct {
	TotalEntries     int
	TotalHeads       int
	TotalTails       int
	TotalNotMattered int
	HeadsMattered    int
}

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
func InsertExeggutor(db *sql.DB, eggType EggType) {
	var heads int
	var mattered bool
	switch eggType {
	case H:
		heads = 1
		mattered = true
	case HX:
		heads = 1
		mattered = false
	case T:
		heads = 0
		mattered = true
	case TX:
		heads = 0
		mattered = false
	}

	stmt := `
	INSERT INTO exeggutor (heads, mattered, created_at)
	VALUES (?, ?, CURRENT_TIMESTAMP)
`
	_, err := db.Exec(stmt, heads, mattered)
	if err != nil {
		fmt.Printf("Failed to insert exeggutor entry: %v\n", err)
	}
}

func GetEggStats(db *sql.DB) (stats EggStats, err error) {
	err = db.QueryRow("SELECT COUNT(*), IFNULL(SUM(heads), 0), IFNULL(SUM(CASE WHEN heads = 0 THEN 1 ELSE 0 END), 0) FROM exeggutor").Scan(&stats.TotalEntries, &stats.TotalHeads, &stats.TotalTails)
	if err != nil {
		return
	}

	err = db.QueryRow("SELECT COUNT(*) FROM exeggutor WHERE mattered = 0").Scan(&stats.TotalNotMattered)
	if err != nil {
		return
	}

	err = db.QueryRow("SELECT COUNT(*) FROM exeggutor WHERE heads = 1 AND mattered = 1").Scan(&stats.HeadsMattered)
	return
}
