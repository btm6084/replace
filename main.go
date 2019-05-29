package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/AstromechZA/pbars"
	"github.com/btm6084/utilities/fileutil"
)

func main() {
	args := os.Args

	if len(args) < 2 {
		fmt.Println("replace <search> <replace> [directory='./'] [extension filter='']")
		fmt.Println("Missing Search Term")
		os.Exit(1)
	}
	if len(args) < 3 {
		fmt.Println("replace <search> <replace> [directory='./'] [extension filter='']")
		fmt.Println("Missing Replacement Term")
		os.Exit(1)
	}

	search := args[1]
	replace := args[2]
	path := "."
	ext := ""

	if len(args) >= 4 {
		path = strings.TrimRight(args[3], "/")
	}
	if len(args) >= 5 {
		ext = fmt.Sprintf(".%s", strings.TrimLeft(args[4], "."))
	}

	if !fileutil.IsDir(path) {
		log.Printf("%s is not a directory\n", path)
		os.Exit(1)
	}

	fmt.Printf("Search:  %s\n", search)
	fmt.Printf("Replace: %s\n", replace)
	fmt.Printf("Path:    %s\n", path)
	fmt.Printf("Filter:  %s\n", ext)

	files := fileutil.DirToArray(path, false, fileutil.DefaultFileFilter, fileutil.DefaultDirectoryFilter)

	// Filter to files with ext if it's set.
	if ext != "" {
		files = fileutil.FilterExtWhitelist([]string{ext}, files)
	}

	// Setup for Progress Bar
	pb := pbars.NewProgressPrinter("Progress", 50, true)
	pb.ShowRate = false

	numOps := int64(len(files))
	step := int64(1)
	percentDone := int64(0)
	active := 0

	var messages []string
	c := make(chan string)

	for _, file := range files {
		go searchAndReplace(file, search, replace, c)
		active++

		if active >= 5 {
			m := <-c
			percentDone += step
			pb.Update(percentDone, numOps)
			active--

			if m != "" {
				messages = append(messages, m)
			}
		}
	}

	// Wait for the last batch of concurrency to wrap up.
	for i := 0; i < active; i++ {
		m := <-c
		percentDone += step
		pb.Update(percentDone, numOps)

		if m != "" {
			messages = append(messages, m)
		}
	}

	pb.Done()

	if len(messages) > 0 {
		fmt.Println(strings.Join(messages, "\n"))
		fmt.Println()
	}

	fmt.Printf("Searched %d files. Replacements made in %d files.\n", len(files), len(messages))
}

// searchAndReplace will replace any instances of 'search' with 'replace'.
// If replacements have been made, Write the contents back to the file.
func searchAndReplace(file, search, replace string, c chan string) {
	s := regexp.MustCompile(search)

	b, err := ioutil.ReadFile(file)
	if err != nil {
		c <- ""
		return
	}

	contents := s.ReplaceAll(b, []byte(replace))

	if string(contents) != string(b) {
		err := ioutil.WriteFile(file, contents, 0644)
		if err != nil {
			c <- ""
			return
		}

		c <- fmt.Sprintf("\tReplacements in %s", file)
		return
	}

	c <- ""
	return
}
