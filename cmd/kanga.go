package cmd

import (
	"database/sql"
	"fmt"

	"github.com/alexstory/kanga/data"
)

func Heads(db *sql.DB) {
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
}

func Tails(db *sql.DB) {
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
}

func Stats(db *sql.DB) {
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
}
