package main

import (
	"os"
	"strings"
	"fmt"
	"sort"
	"io"
)

func dirTree(out io.Writer, path string, force bool) (error) {
	f, err := os.Open(path)
	if err != nil {
		return err
	}

	list, err := f.Readdir(-1)
	f.Close()
	if err != nil {
		return err
	}
	sort.Slice(list, func(i, j int) bool { return list[i].Name() < list[j].Name() })

	fileLevel := strings.Count(path, "/")
	separator := "├───"
	for i := 1; i < fileLevel; i++ {
		separator = "│\t" + separator
	}
	
	for i := 0; i < len(list); i++ {

		if i == len(list)-1 {
			separator = strings.Replace(separator, "├───", "└───", -1)
		}

		if force || list[i].IsDir() {
			if (!list[i].IsDir()) {
				fileSize := list[i].Size()
				if fileSize == 0 {
					fmt.Fprintf(out, "%s%s (empty)\n", separator, list[i].Name())
				} else {
					fmt.Fprintf(out, "%s%s (%db)\n", separator, list[i].Name(), fileSize)
				}
			} else {
				fmt.Fprintf(out, "%s%s\n", separator, list[i].Name())
			}
		}

		if list[i].IsDir() {
			dirTree(out, generatePath(path, list[i].Name()), force)
		}
	}

	return nil
}

func generatePath(parent string, child string) string {
	var paths []string
	paths = append(paths, parent)
	paths = append(paths, child)
	return strings.Join(paths, string(os.PathSeparator))
}

func main() {
	out := os.Stdout
	if !(len(os.Args) == 2 || len(os.Args) == 3) {
		panic("usage go run main.go . [-f]")
	}
	path := os.Args[1]
	printFiles := len(os.Args) == 3 && os.Args[2] == "-f"
	err := dirTree(out, path, printFiles)
	if err != nil {
		panic(err.Error())
	}
}