package main

import (
	"flag"
	"fmt"
	"strings"

	"github.com/alexstory/kanga/data"
)

type LabelValuePair struct {
	Label string
	Value string
}

func main() {
	flag.Parse()
	data.Init()

	db, err := data.Init()
	if err != nil {
		fmt.Printf("Failed to initialize database: %v\n", err)
		return
	}
	defer db.Close()

	switch flag.Arg(0) {
	case "heads":
		flips, heads, err := data.HeadsInfo(db)
		if err != nil {
			fmt.Printf("Failed to get heads info: %v\n", err)
		}

		dataPairs := []LabelValuePair{
			{"Total flips", fmt.Sprintf("%d", flips)},
			{"Heads count", fmt.Sprintf("%d", heads)},
			{"Heads percentage", fmt.Sprintf("%.2f%%", percentage(heads, flips))},
		}

		printTable("HEADS INFO", dataPairs)

	case "tails":
		flips, tails, err := data.TailsInfo(db)
		if err != nil {
			fmt.Printf("Failed to get tails info: %v\n", err)
		}

		dataPairs := []LabelValuePair{
			{"Total flips", fmt.Sprintf("%d", flips)},
			{"Tails count", fmt.Sprintf("%d", tails)},
			{"Tails percentage", fmt.Sprintf("%.2f%%", percentage(tails, flips))},
		}

		printTable("TAILS INFO", dataPairs)

	case "stats":
		stats, err := data.Flips(db)
		if err != nil {
			fmt.Printf("Failed to get stats: %v\n", err)
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
		data.TT(db)
		fmt.Printf("flip logged... RIP\n")
	case "HH", "hh":
		data.HH(db)
		fmt.Printf("flip logged...\n")
	case "HT", "ht":
		data.HT(db)
		fmt.Printf("flip logged...\n")
	case "TH", "th":
		data.TH(db)
		fmt.Printf("flip logged...\n")
	case "reset":
		data.Reset(db)
		fmt.Printf("Data reset\n")
	case "undo":
		data.Undo(db)
		fmt.Printf("Last flip undone\n")
	case "dump-csv":
		if flag.NArg() < 2 {
			fmt.Println("Usage: kanga dump-csv <filename>")
			return
		}
		filename := flag.Arg(1)
		err := data.DumpCsv(db, filename)
		if err != nil {
			fmt.Printf("Failed to dump CSV: %v\n", err)
		} else {
			fmt.Printf("data dumped to dumped to %s\n", filename)
		}
	case "read-csv":
		if flag.NArg() < 2 {
			fmt.Println("Usage: kanga read-csv <filename>")
			return
		}

		filename := flag.Arg(1)
		err := data.ReadCsv(db, filename)
		if err != nil {
			fmt.Printf("Failed to read CSV: %v\n", err)
		} else {
			fmt.Printf("data read from %s\n", filename)
		}

	default:
		printHelp()
	}
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

func printHelp() {
	fmt.Println("Usage: kanga [command]")
	fmt.Println("Commands:")
	fmt.Println("  heads       Show heads info")
	fmt.Println("  tails       Show tails info")
	fmt.Println("  stats       Show statistics")
	fmt.Println("  TT, tt      Log a double tails flip")
	fmt.Println("  HH, hh      Log a double heads flip")
	fmt.Println("  HT, ht      Log a heads-tails flip")
	fmt.Println("  TH, th      Log a tails-heads flip")
	fmt.Println("  reset       Reset the data")
	fmt.Println("  undo        Undo the last flip")
	fmt.Println("  dump-csv    Dump the data to a CSV file")
	fmt.Println("  read-csv    Read the data from a CSV file (! this overwrites the current dataset)")
}
