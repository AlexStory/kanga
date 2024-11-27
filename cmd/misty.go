package cmd

import (
	"database/sql"
	"flag"
	"fmt"
	"strconv"

	"github.com/alexstory/kanga/data"
)

func InsertMisty(db *sql.DB, heads int) {
	data.InsertMisty(db, heads)
	fmt.Printf("Entry logged...\n")
}

func MistyStats(db *sql.DB) {
	var funStat LabelValuePair

	stats, err := data.GetMistyStats(db)
	if err != nil {
		fmt.Printf("Failed to get misty stats: %v\n", err)
		return
	}

	if stats.TotalEntries == 0 {
		funStat = LabelValuePair{"Soul", "clean"}
	} else {
		funStat = LabelValuePair{"Sins", "infinite"}
	}

	dataPairs := []LabelValuePair{
		{"Total attempts", fmt.Sprintf("%d", stats.TotalEntries)},
		{"Total heads", fmt.Sprintf("%d", stats.TotalHeads)},
		funStat,
	}
	printTable("MISTY STATS", dataPairs)
}

func Misty(db *sql.DB) {
	if flag.NArg() < 2 {
		fmt.Println("Usage: kanga misty <command>")
		fmt.Println("See `kanga help misty` for more info")
		return
	}

	arg := flag.Arg(1)

	heads, err := strconv.Atoi(arg)
	if err == nil {
		InsertMisty(db, heads)
		return
	}

	switch arg {
	case "stats":
		MistyStats(db)
	case "undo":
		data.UndoMisty(db)
		fmt.Println("Last misty flip undone...")
	default:
		fmt.Println("Invalid argument for misty command.")
		fmt.Println("See `kanga help misty` for more info")
	}

}
