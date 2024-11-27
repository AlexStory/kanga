package cmd

import "fmt"

func PrintHelp(command string) {
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
		fmt.Println("Usage: kanga egg <command>")
		fmt.Println("Log an exeggutor entry or show stats")
		fmt.Println("  H        Log a heads")
		fmt.Println("  T        Log a tails")
		fmt.Println("  HX       Log a heads, but... the result didn't really matter")
		fmt.Println("  TX       Log a tails, but... the result didn't really matter")
		fmt.Println("  stats    Show exeggutor statistics")
		fmt.Println("  undo     Undo the last exeggutor entry")
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
