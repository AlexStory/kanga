package cmd

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/alexstory/kanga/data"
)

type LabelValuePair struct {
	Label string
	Value string
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

func ReadCsv(db *sql.DB, folder string, table string) {
	err := data.ReadCsv(db, folder, table)
	if err != nil {
		fmt.Printf("Failed to read CSV: %v\n", err)
	} else {
		fmt.Printf("Data read from %s\n", folder)
	}
}
