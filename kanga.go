package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/alexstory/kanga/cmd"
	"github.com/alexstory/kanga/data"
)

const (
	kanga int = iota
	egg
	misty
)

func main() {
	kangaFlag := flag.Bool("kanga", false, "Operate on kanga table")
	eggFlag := flag.Bool("egg", false, "Operate on exeggutor table")
	mistyFlag := flag.Bool("misty", false, "Operate on misty table")
	flag.Parse()

	tables := map[data.TableType]bool{
		data.Kanga: *kangaFlag,
		data.Egg:   *eggFlag,
		data.Misty: *mistyFlag,
	}

	db, err := data.Init()
	if err != nil {
		log.Fatalf("Failed to initialize database: %v", err)
	}
	defer db.Close()

	if flag.NArg() == 0 {
		cmd.PrintHelp("")
		return
	}

	command := flag.Arg(0)
	folder := "."
	if flag.NArg() >= 2 {
		folder = flag.Arg(1)
	}

	switch command {
	case "heads":
		cmd.Heads(db)
	case "tails":
		cmd.Tails(db)
	case "stats":
		cmd.Stats(db)
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
		cmd.Egg(db)
	case "misty":
		cmd.Misty(db)
	case "reset":
		data.Reset(db)
		fmt.Printf("Data reset\n")
	case "undo":
		data.Undo(db)
		fmt.Printf("Last flip undone\n")
	case "dump-csv":
		err := data.DumpCsv(db, folder, tables)
		if err != nil {
			fmt.Printf("Failed to dump CSV: %v\n", err)
		} else {
			fmt.Printf("Data dumped to %s\n", folder)
		}
	case "read-csv":
		var table string
		if *kangaFlag {
			table = "kanga"
		} else if *eggFlag {
			table = "egg"
		}
		cmd.ReadCsv(db, folder, table)
	case "help":
		if flag.NArg() < 2 {
			cmd.PrintHelp("")
		} else {
			cmd.PrintHelp(flag.Arg(1))
		}
	default:
		cmd.PrintHelp("")
	}
}
