package main

import "fmt"
import "os"
import "flag"
import "strings"
import "path/filepath"
import "io/ioutil"
import "bufio"

type Color string

const (
	ColorBlack  Color = "\u001b[30m"
	ColorRed          = "\u001b[31m"
	ColorGreen        = "\u001b[32m"
	ColorYellow       = "\u001b[33m"
	ColorReset        = "\u001b[0m"
)

func colorizePattern(color Color, src string, pattern string) {
	tmp := strings.Split(src, pattern)
	format := strings.Join(tmp, string(ColorRed)+pattern+string(ColorReset))
	fmt.Println(format)
}

func colorize(color Color, message string) {
	fmt.Println(string(color) + message + string(ColorReset))
}

func main() {
	args := os.Args[1:]
	rflag := flag.Bool("R", false, "Recursive search")
	fflag := flag.Bool("F", false, "Deep search (read file content)")

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
			if *fflag {
				DeepSearch(file, args)
			} else {
				for _, pattern := range args {
					if strings.Contains(file, pattern) {
						colorizePattern(ColorRed, file, pattern)
					}
				}
			}

		}
	} else if *rflag == false {
		files, err := ioutil.ReadDir(rootDir)
		if err != nil {
			panic(err)
		}

		for _, file := range files {
			if *fflag {
				DeepSearch(file.Name(), args)
			} else {
				for _, pattern := range args {
					if strings.Contains(file.Name(), pattern) {
						colorizePattern(ColorRed, file.Name(), pattern)
					}
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

func DeepSearch(fn string, args []string) {
	f, _ := os.Open(fn)
	scanner := bufio.NewScanner(f)
	lineCounter := 0
	fnPrinted := false
	for scanner.Scan() {
		lineCounter += 1
		tmp := strings.Trim(strings.Trim(scanner.Text(), " "), "\t")
		for _, pattern := range args {
			if strings.Contains(tmp, pattern) {
				// check and print filename only once
				if fnPrinted == false {
					fmt.Println()
					dir, _ := os.Getwd()
					fmt.Println(dir + "/" + fn)
					fnPrinted = true
				}

				// Print line count and the matching text-pattern
				fmt.Print(lineCounter, ":  ")
				colorizePattern(ColorRed, tmp, pattern)
			}
		}

	}
}
