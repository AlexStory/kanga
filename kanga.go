package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/alexstory/kanga/data"
)

type LabelValuePair struct {
	Label string
	Value string
}

func main() {
	kangaFlag := flag.Bool("kanga", false, "Operate on kanga table")
	eggFlag := flag.Bool("egg", false, "Operate on exeggutor table")
	flag.Parse()

	db, err := data.Init()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	if flag.NArg() == 0 {
		printHelp("")
		return
	}

	command := flag.Arg(0)
	folder := "."
	if flag.NArg() >= 2 {
		folder = flag.Arg(1)
	}

	switch command {
	case "heads":
		totalFlips, headsCount, err := data.HeadsInfo(db)
		if err != nil {
			fmt.Printf("Failed to get heads info: %v\n", err)
			return
		}
		dataPairs := []LabelValuePair{
			{"Total flips", fmt.Sprintf("%d", totalFlips)},
			{"Heads count", fmt.Sprintf("%d", headsCount)},
			{"Heads percentage", fmt.Sprintf("%.2f%%", percentage(headsCount, totalFlips))},
		}
		printTable("HEADS INFO", dataPairs)
	case "tails":
		totalFlips, tailsCount, err := data.TailsInfo(db)
		if err != nil {
			fmt.Printf("Failed to get tails info: %v\n", err)
			return
		}
		dataPairs := []LabelValuePair{
			{"Total flips", fmt.Sprintf("%d", totalFlips)},
			{"Tails count", fmt.Sprintf("%d", tailsCount)},
			{"Tails percentage", fmt.Sprintf("%.2f%%", percentage(tailsCount, totalFlips))},
		}
		printTable("TAILS INFO", dataPairs)
	case "stats":
		stats, err := data.Flips(db)
		if err != nil {
			fmt.Printf("Failed to get stats: %v\n", err)
			return
		}
		dataPairs := []LabelValuePair{
			{"Total flips", fmt.Sprintf("%d", stats.TotalFlips)},
			{"Double heads", fmt.Sprintf("%d", stats.DoubleHeads)},
			{"Double tails", fmt.Sprintf("%d", stats.DoubleTails)},
			{"Total heads", fmt.Sprintf("%d", stats.TotalHeads)},
			{"Total tails", fmt.Sprintf("%d", stats.TotalTails)},
			{"Heads percentage", fmt.Sprintf("%.2f%%", percentage(stats.TotalHeads, stats.TotalFlips))},
			{"Tails percentage", fmt.Sprintf("%.2f%%", percentage(stats.TotalTails, stats.TotalFlips))},
			{"Double heads percentage", fmt.Sprintf("%.2f%%", percentage(stats.DoubleHeads, stats.TotalFlips))},
			{"Double tails percentage", fmt.Sprintf("%.2f%%", percentage(stats.DoubleTails, stats.TotalFlips))},
		}
		printTable("STATISTICS", dataPairs)
	case "TT", "tt":
		data.InsertFlip(db, data.TT)
		fmt.Printf("flip logged... RIP\n")
	case "HH", "hh":
		data.InsertFlip(db, data.HH)
		fmt.Printf("flip logged...\n")
	case "HT", "ht":
		data.InsertFlip(db, data.HT)
		fmt.Printf("flip logged...\n")
	case "TH", "th":
		data.InsertFlip(db, data.TH)
		fmt.Printf("flip logged...\n")
	case "egg":
		handleEggCommand(db)
	case "reset":
		data.Reset(db)
		fmt.Printf("Data reset\n")
	case "undo":
		data.Undo(db)
		fmt.Printf("Last flip undone\n")
	case "dump-csv":
		table := ""
		if *kangaFlag {
			table = "kanga"
		} else if *eggFlag {
			table = "egg"
		}
		err := data.DumpCsv(db, folder, table)
		if err != nil {
			fmt.Printf("Failed to dump CSV: %v\n", err)
		} else {
			fmt.Printf("Data dumped to %s\n", folder)
		}
	case "read-csv":
		table := ""
		if *kangaFlag {
			table = "kanga"
		} else if *eggFlag {
			table = "egg"
		}
		err := data.ReadCsv(db, folder, table)
		if err != nil {
			fmt.Printf("Failed to read CSV: %v\n", err)
		} else {
			fmt.Printf("Data read from %s\n", folder)
		}
	case "help":
		if flag.NArg() < 2 {
			printHelp("")
		} else {
			printHelp(flag.Arg(1))
		}
	default:
		printHelp("")
	}
}

func handleEggCommand(db *sql.DB) {
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
		}
		printTable("EGGEGGUTOR STATS", dataPairs)
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
	default:
		fmt.Println("Invalid argument for egg command.")
		fmt.Println("use `kanga help egg` for more info")
		return
	}
	data.InsertExeggutor(db, eggType)
	fmt.Println("Exeggutor entry logged...")
}

func percentage(part, total int) float64 {
	if total == 0 {
		return 0
	}
	return float64(part) / float64(total) * 100
}

func printTable(title string, dataPairs []LabelValuePair) {
	// Pre-format the strings without the right border
	lines := make([]string, len(dataPairs))
	for i, pair := range dataPairs {
		lines[i] = fmt.Sprintf("| %-24s | %6s ", pair.Label, pair.Value)
	}

	// Determine the maximum length
	maxLength := 0
	for _, line := range lines {
		if len(line) > maxLength {
			maxLength = len(line)
		}
	}

	// Print the table with even borders
	border := "+" + strings.Repeat("-", maxLength+1) + "+"
	headerPadding := (maxLength - len(title)) / 2
	fmt.Println(border)
	fmt.Println("|" + strings.Repeat(" ", headerPadding) + title + strings.Repeat(" ", maxLength-len(title)-headerPadding) + " |")
	fmt.Println(border)
	for _, line := range lines {
		fmt.Println(line + strings.Repeat(" ", maxLength-len(line)+2) + "|")
	}
	fmt.Println(border)
}

func printHelp(command string) {
	switch command {
	case "heads":
		fmt.Println("Usage: kanga heads")
		fmt.Println("Show heads info")
	case "tails":
		fmt.Println("Usage: kanga tails")
		fmt.Println("Show tails info")
	case "stats":
		fmt.Println("Usage: kanga stats")
		fmt.Println("Show statistics")
	case "TT", "tt":
		fmt.Println("Usage: kanga TT")
		fmt.Println("Log a double tails flip")
	case "HH", "hh":
		fmt.Println("Usage: kanga HH")
		fmt.Println("Log a double heads flip")
	case "HT", "ht":
		fmt.Println("Usage: kanga HT")
		fmt.Println("Log a heads-tails flip")
	case "TH", "th":
		fmt.Println("Usage: kanga TH")
		fmt.Println("Log a tails-heads flip")
	case "egg":
		fmt.Println("Usage: kanga egg <H|HX|T|TX|stats>")
		fmt.Println("Log an exeggutor entry or show stats")
		fmt.Println("  H   - Log a heads")
		fmt.Println("  T   - Log a tails")
		fmt.Println("  HX  - Log a heads, but... the result didn't really matter")
		fmt.Println("  TX  - Log a tails, but... the result didn't really matter")
		fmt.Println("  stats - Show exeggutor statistics")
	case "reset":
		fmt.Println("Usage: kanga reset")
		fmt.Println("Reset the database")
	case "undo":
		fmt.Println("Usage: kanga undo")
		fmt.Println("Undo the last action")
	case "dump-csv":
		fmt.Println("Usage: kanga dump-csv [folder]")
		fmt.Println("Dump the data to CSV files in the specified folder (default: current directory)")
	case "read-csv":
		fmt.Println("Usage: kanga read-csv [folder]")
		fmt.Println("Read the data from CSV files in the specified folder (default: current directory)")
	default:
		fmt.Println("Usage: kanga [command]")
		fmt.Println("Commands:")
		fmt.Println("  heads       Show heads info")
		fmt.Println("  tails       Show tails info")
		fmt.Println("  stats       Show statistics")
		fmt.Println("  TT, tt      Log a double tails flip")
		fmt.Println("  HH, hh      Log a double heads flip")
		fmt.Println("  HT, ht      Log a heads-tails flip")
		fmt.Println("  TH, th      Log a tails-heads flip")
		fmt.Println("  egg         Log an exeggutor entry (H, HX, T, TX) or show stats")
		fmt.Println("  reset       Reset the database")
		fmt.Println("  undo        Undo the last action")
		fmt.Println("  dump-csv    Dump the data to CSV files")
		fmt.Println("  read-csv    Read the data from CSV files")
		fmt.Println("  help        Show this help message, or help for a specific command")
	}
}
