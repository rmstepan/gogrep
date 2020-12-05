package main

import "fmt"
import "os"
import "flag"
import "strings"
import "path/filepath"

type Color string

const (
	ColorBlack  Color = "\u001b[30m"
	ColorRed          = "\u001b[31m"
	ColorGreen        = "\u001b[32m"
	ColorYellow       = "\u001b[33m"
	ColorReset        = "\u001b[0m"
)

func colorizePattern(color Color, src string, pattern string){
	tmp := strings.Split(src, pattern)
	format := strings.Join(tmp, string(ColorRed) + pattern + string(ColorReset))
	fmt.Println(format)
}

func colorize(color Color, message string) {
	fmt.Println(string(color) + message + string(ColorReset))
}

func main() {
	args := os.Args[1:]
	rflag := flag.Bool("R", false, "Recursive search")
	//fflag := flag.Bool("F", false, "Deep search (read file content)")

	flag.Parse()

	args = ClearArgs(args)

	if len(args) == 0 {
		PrintUsage()
		return 
	}

	files := []string{}
	rootDir := "."

	if *rflag {
		err := filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			files = append(files, path)
			return nil
		})

		if err != nil {
			panic(err)
		}

		for _, file := range files {
			for _, pattern := range args {
				if strings.Contains(file, pattern) {
					colorizePattern(ColorRed, file, pattern)
				}
			}
		}
	}
}

func PrintUsage() {
	fmt.Println("Usage: gogrep <flag> <search_string1> <search_string2> ... <search_stringN>")
	fmt.Println("Example: gogrep -R #define import #include *.py")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("\t-R\t:search recursively for the given pattern")
	fmt.Println("\t-F\t:Deep search(read files content)")

}

func ClearArgs(args []string) []string {
	for i := 0; i < len(args); i++ {
		// remove flags from args list
		if strings.HasPrefix(args[i], "-") {
			args = append(args[:i], args[i+1:]...)
			i--
		}
	}
	return args
}

