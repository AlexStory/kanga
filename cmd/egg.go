package cmd

import (
	"database/sql"
	"flag"
	"fmt"

	"github.com/alexstory/kanga/data"
)

func Egg(db *sql.DB) {
	if flag.NArg() < 2 {
		fmt.Println("Usage: kanga egg <H|HX|T|TX|stats>")
		fmt.Println("Log an exeggutor entry or show stats")
		fmt.Println("  H   - Log a heads")
		fmt.Println("  T   - Log a tails")
		fmt.Println("  HX  - Log a heads, but... the result didn't really matter")
		fmt.Println("  TX  - Log a tails, but... the result didn't really matter")
		fmt.Println("  stats - Show exeggutor statistics")
		return
	}
	arg := flag.Arg(1)
	if arg == "stats" {
		stats, err := data.GetEggStats(db)
		if err != nil {
			fmt.Printf("Failed to get egg stats: %v\n", err)
			return
		}
		dataPairs := []LabelValuePair{
			{"Total flips", fmt.Sprintf("%d", stats.TotalEntries)},
			{"Total heads", fmt.Sprintf("%d", stats.TotalHeads)},
			{"Total tails", fmt.Sprintf("%d", stats.TotalTails)},
			{"Heads percentage", fmt.Sprintf("%.2f%%", percentage(stats.TotalHeads, stats.TotalEntries))},
			{"Tails percentage", fmt.Sprintf("%.2f%%", percentage(stats.TotalTails, stats.TotalEntries))},
			{"Heads that mattered", fmt.Sprintf("%d", stats.HeadsMattered)},
			{"Percent when it mattered", fmt.Sprintf("%.2f%%", percentage(stats.HeadsMattered, stats.TotalEntries-stats.TotalNotMattered))},
			{"Average damage", fmt.Sprintf("%d", ((stats.TotalHeads*80)+(stats.TotalTails*40))/stats.TotalEntries)},
		}
		printTable("EXEGGUTOR STATS", dataPairs)
		return
	}
	var eggType data.EggType
	switch arg {
	case "H", "h":
		eggType = data.H
	case "HX", "hx":
		eggType = data.HX
	case "T", "t":
		eggType = data.T
	case "TX", "tx":
		eggType = data.TX
	case "undo":
		data.UndoEgg(db)
		fmt.Println("Last exeggutor flip undone...")
		return
	default:
		fmt.Println("Invalid argument for egg command.")
		fmt.Println("use `kanga help egg` for more info")
		return
	}
	data.InsertExeggutor(db, eggType)
	fmt.Println("Exeggutor entry logged...")
}
