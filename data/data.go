package data

import (
	"database/sql"
	"encoding/csv"
	"fmt"
	"os"
	"path/filepath"

	_ "modernc.org/sqlite"
)

type FlipType int

const (
	TT FlipType = iota
	HH
	TH
	HT
)

type EggType int

const (
	H EggType = iota
	HX
	T
	TX
)

type Stats struct {
	TotalFlips  int
	DoubleHeads int
	DoubleTails int
	TotalHeads  int
	TotalTails  int
}

type EggStats struct {
	TotalEntries     int
	TotalHeads       int
	TotalTails       int
	TotalNotMattered int
	HeadsMattered    int
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

	createFlipsTableSQL := `CREATE TABLE IF NOT EXISTS flips (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		heads1 INTEGER NOT NULL,
		heads2 INTEGER NOT NULL,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createFlipsTableSQL)
	if err != nil {
		return nil, err
	}

	createExeggutorTableSQL := `CREATE TABLE IF NOT EXISTS exeggutor (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		heads INTEGER NOT NULL,
		mattered BOOLEAN DEFAULT TRUE,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP
	);`
	_, err = db.Exec(createExeggutorTableSQL)
	if err != nil {
		return nil, err
	}

	// Create indexes
	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_heads1 ON flips (heads1);")
	if err != nil {
		return nil, err
	}

	_, err = db.Exec("CREATE INDEX IF NOT EXISTS idx_heads2 ON flips (heads2);")
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

func InsertFlip(db *sql.DB, flipType FlipType) {
	var heads1, heads2 int
	switch flipType {
	case TT:
		heads1, heads2 = 0, 0
	case HH:
		heads1, heads2 = 1, 1
	case TH:
		heads1, heads2 = 0, 1
	case HT:
		heads1, heads2 = 1, 0
	}

	stmt := `
	INSERT INTO flips (heads1, heads2, created_at)
	VALUES (?, ?, CURRENT_TIMESTAMP)
`
	_, err := db.Exec(stmt, heads1, heads2)
	if err != nil {
		fmt.Printf("Failed to insert flip: %v\n", err)
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

func Reset(db *sql.DB) {
	_, err := db.Exec("DELETE FROM flips")
	if err != nil {
		fmt.Printf("Failed to reset flips table: %v\n", err)
	}
	_, err = db.Exec("DELETE FROM exeggutor")
	if err != nil {
		fmt.Printf("Failed to reset exeggutor table: %v\n", err)
	}
}

func Undo(db *sql.DB) {
	stmt := `
	DELETE FROM flips
	WHERE id = (SELECT MAX(id) FROM flips)
`
	_, err := db.Exec(stmt)
	if err != nil {
		fmt.Printf("Failed to undo flip: %v\n", err)
	}
}

func DumpCsv(db *sql.DB, folder string, table string) error {
	if table == "" || table == "kanga" {
		err := dumpTable(db, folder, "kanga")
		if err != nil {
			return err
		}
	}
	if table == "" || table == "egg" {
		err := dumpTable(db, folder, "egg")
		if err != nil {
			return err
		}
	}
	return nil
}

func dumpTable(db *sql.DB, folder string, table string) error {
	var query, filename string
	switch table {
	case "kanga":
		query = "SELECT heads1, heads2, created_at FROM flips"
		filename = "kanga.csv"
	case "egg":
		query = "SELECT heads, mattered, created_at FROM exeggutor"
		filename = "exeggutor.csv"
	default:
		return fmt.Errorf("unknown table: %s", table)
	}

	// Create the folder if it doesn't exist
	if err := os.MkdirAll(folder, os.ModePerm); err != nil {
		return err
	}

	rows, err := db.Query(query)
	if err != nil {
		return err
	}
	defer rows.Close()

	filePath := filepath.Join(folder, filename)
	file, err := os.Create(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	for rows.Next() {
		var record []string
		if table == "kanga" {
			var heads1, heads2 int
			var createdAt string
			err := rows.Scan(&heads1, &heads2, &createdAt)
			if err != nil {
				return err
			}
			record = []string{fmt.Sprintf("%d", heads1), fmt.Sprintf("%d", heads2), createdAt}
		} else {
			var heads int
			var mattered bool
			var createdAt string
			err := rows.Scan(&heads, &mattered, &createdAt)
			if err != nil {
				return err
			}
			record = []string{fmt.Sprintf("%d", heads), fmt.Sprintf("%t", mattered), createdAt}
		}
		err = writer.Write(record)
		if err != nil {
			return err
		}
	}

	return rows.Err()
}

func ReadCsv(db *sql.DB, folder string, table string) error {
	if table == "" || table == "kanga" {
		err := readTable(db, folder, "kanga")
		if err != nil {
			return err
		}
	}
	if table == "" || table == "egg" {
		err := readTable(db, folder, "egg")
		if err != nil {
			return err
		}
	}
	return nil
}

func readTable(db *sql.DB, folder string, table string) error {
	var query, filename string
	switch table {
	case "kanga":
		query = "INSERT INTO flips (heads1, heads2, created_at) VALUES (?, ?, ?)"
		filename = "kanga.csv"
	case "egg":
		query = "INSERT INTO exeggutor (heads, mattered, created_at) VALUES (?, ?, ?)"
		filename = "exeggutor.csv"
	default:
		return fmt.Errorf("unknown table: %s", table)
	}

	filePath := filepath.Join(folder, filename)
	file, err := os.Open(filePath)
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

	stmt, err := tx.Prepare(query)
	if err != nil {
		return err
	}
	defer stmt.Close()

	for _, record := range records {
		if table == "kanga" {
			if len(record) != 3 {
				tx.Rollback()
				return fmt.Errorf("invalid record: %v", record)
			}
			heads1 := record[0]
			heads2 := record[1]
			createdAt := record[2]
			_, err := stmt.Exec(heads1, heads2, createdAt)
			if err != nil {
				tx.Rollback()
				return err
			}
		} else {
			if len(record) != 3 {
				tx.Rollback()
				return fmt.Errorf("invalid record: %v", record)
			}
			heads := record[0]
			mattered := record[1]
			createdAt := record[2]
			_, err := stmt.Exec(heads, mattered, createdAt)
			if err != nil {
				tx.Rollback()
				return err
			}
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
