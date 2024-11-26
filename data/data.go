package data

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type Stats struct {
	TotalFlips  int
	DoubleHeads int
	DoubleTails int
	TotalHeads  int
	TotalTails  int
}

func Init() (*sql.DB, error) {
	exePath, err := os.Executable()
	if err != nil {
		return nil, err
	}

	dbPath := filepath.Join(filepath.Dir(exePath), "kanga.db")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, err
	}

	// Enable WAL mode
	_, err = db.Exec("PRAGMA journal_mode=WAL;")
	if err != nil {
		return nil, err
	}

	createTableSQL := `CREATE TABLE IF NOT EXISTS flips (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		heads1 INTEGER NOT NULL,
		heads2 INTEGER NOT NULL
	);`
	_, err = db.Exec(createTableSQL)
	if err != nil {
		return nil, err
	}

	return db, nil
}

func HeadsInfo(db *sql.DB) (totalFlips, headsCount int, err error) {
	var rowCount int
	err = db.QueryRow("SELECT COUNT(*) FROM flips").Scan(&rowCount)
	if err != nil {
		return
	}
	totalFlips = rowCount * 2

	err = db.QueryRow("SELECT IFNULL(SUM(heads1 + heads2), 0) FROM flips WHERE heads1 = 1 OR heads2 = 1").Scan(&headsCount)
	return
}

func TailsInfo(db *sql.DB) (totalFlips, tailsCount int, err error) {
	var rowCount int
	err = db.QueryRow("SELECT COUNT(*) FROM flips").Scan(&rowCount)
	if err != nil {
		return
	}
	totalFlips = rowCount * 2

	err = db.QueryRow("SELECT IFNULL(SUM((1 - heads1) + (1 - heads2)), 0) FROM flips WHERE heads1 = 0 OR heads2 = 0").Scan(&tailsCount)
	return
}

func Flips(db *sql.DB) (stats Stats, err error) {
	var rowCount int
	err = db.QueryRow("SELECT COUNT(*) FROM flips").Scan(&rowCount)
	if err != nil {
		return
	}
	stats.TotalFlips = rowCount * 2

	err = db.QueryRow("SELECT IFNULL(COUNT(*), 0) FROM flips WHERE heads1 = 1 AND heads2 = 1").Scan(&stats.DoubleHeads)
	if err != nil {
		return
	}

	err = db.QueryRow("SELECT IFNULL(COUNT(*), 0) FROM flips WHERE heads1 = 0 AND heads2 = 0").Scan(&stats.DoubleTails)
	if err != nil {
		return
	}

	err = db.QueryRow("SELECT IFNULL(SUM(heads1 + heads2), 0) FROM flips").Scan(&stats.TotalHeads)
	if err != nil {
		return
	}

	err = db.QueryRow("SELECT IFNULL(SUM((1 - heads1) + (1 - heads2)), 0) FROM flips").Scan(&stats.TotalTails)
	return
}

func TT(db *sql.DB) {
	stmt := `
	INSERT INTO flips (heads1, heads2)
	VALUES (0, 0)	
`
	_, err := db.Exec(stmt)
	if err != nil {
		fmt.Printf("Failed to insert flip: %v\n", err)
	}
}

func HH(db *sql.DB) {
	stmt := `
	INSERT INTO flips (heads1, heads2)
	VALUES (1, 1)	
`
	_, err := db.Exec(stmt)
	if err != nil {
		fmt.Printf("Failed to insert flip: %v\n", err)
	}
}

func HT(db *sql.DB) {
	stmt := `
	INSERT INTO flips (heads1, heads2)
	VALUES (1, 0)	
`
	_, err := db.Exec(stmt)
	if err != nil {
		fmt.Printf("Failed to insert flip: %v\n", err)
	}
}

func TH(db *sql.DB) {
	stmt := `
	INSERT INTO flips (heads1, heads2)
	VALUES (0, 1)	
`
	_, err := db.Exec(stmt)
	if err != nil {
		fmt.Printf("Failed to insert flip: %v\n", err)
	}
}

func Reset(db *sql.DB) {
	_, err := db.Exec("DELETE FROM flips")
	if err != nil {
		fmt.Printf("Failed to reset database: %v\n", err)
	}
}

func DumpCsv(db *sql.DB, path string) error {
	rows, err := db.Query("SELECT heads1, heads2 FROM flips")
	if err != nil {
		return err
	}
	defer rows.Close()

	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for rows.Next() {
		var heads1, heads2 int
		err := rows.Scan(&heads1, &heads2)
		if err != nil {
			return err
		}
		record := []string{fmt.Sprintf("%d", heads1), fmt.Sprintf("%d", heads2)}
		err = writer.Write(record)
		if err != nil {
			return err
		}
	}

	return rows.Err()
}

func ReadCsv(db *sql.DB, path string) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	if err != nil {
		return err
	}

	tx, err := db.Begin()
	if err != nil {
		return err
	}

	stmt, err := tx.Prepare("INSERT INTO flips (heads1, heads2) VALUES (?, ?)")
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, record := range records {
		if len(record) != 2 {
			tx.Rollback()
			return fmt.Errorf("invalid record: %v", record)
		}
		heads1 := record[0]
		heads2 := record[1]
		_, err := stmt.Exec(heads1, heads2)
		if err != nil {
			tx.Rollback()
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}
	return nil
}
