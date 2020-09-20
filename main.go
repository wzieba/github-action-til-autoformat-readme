package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"regexp"
	"sort"
	"strconv"
	"strings"
	"text/template"
	"time"
)

var repoPath = os.Getenv("REPO_PATH")
var repoName = os.Getenv("GITHUB_REPOSITORY") // Default environment on Github
var templatePath = os.Getenv("TEMPLATE_PATH")
var inputDescription = os.Getenv("INPUT_DESCRIPTION")
var inputFooter = os.Getenv("INPUT_FOOTER")
var inputListMostRecent = os.Getenv("INPUT_LIST_MOST_RECENT")
var inputDateFormat = os.Getenv("INPUT_DATE_FORMAT")
var inputTilsCounterFormat = os.Getenv("INPUT_TILS_COUNTER_FORMAT")

var re = regexp.MustCompile(`^Date:\s*`)
var re2 = regexp.MustCompile(`^#\s*`)

type Til struct {
	Title        string
	AbsoluteLink string
	Category     string
	DateAdded    time.Time
}

// sort TILs by DateAdded (DESC) and return n most recent
func cmdTrimMostRecentTils(tils *[]Til, n int) {
	if n <= 0 {
		n = 0
	}
	if n > len(*tils) {
		n = len(*tils)
	}
	sort.Slice(*tils, func(i, j int) bool {
		first := (*tils)[i].DateAdded
		second := (*tils)[j].DateAdded
		return first.After(second)
	})
	*tils = (*tils)[0:n]
}

// run a git cli command, the capture and parse the output to extract the date
// a file was added to the repository
func cmdGetDate(file string) time.Time {
	c1 := exec.Command("git", "log", "--diff-filter=A", "--date=rfc", "--", file)
	c1.Dir = repoPath
	var commandOutput bytes.Buffer
	var commandErrorOutput bytes.Buffer
	c1.Stdout = &commandOutput
	c1.Stderr = &commandErrorOutput

	err := c1.Start()
	if err != nil {
		fmt.Println("start error")
		fmt.Println(commandErrorOutput.String())
		fmt.Println(file)
		fmt.Println(err)
		return time.Time{}
	}

	err = c1.Wait()
	if err != nil {
		fmt.Println("finish error")
		fmt.Println(commandErrorOutput.String())
		fmt.Println(file)
		fmt.Println(err)
		return time.Time{}
	}

	date := time.Time{}
	for _, outputLine := range strings.Split(commandOutput.String(), "\n") {
		if re.MatchString(outputLine) {
			// strip "Date: " substring from matching line
			var strippedDate = re.ReplaceAllString(outputLine, "")
			date, _ = time.Parse(time.RFC1123Z, strippedDate)
			break
		}
	}
	return date
}

func main() {
	var absoluteLinkFormat = "https://github.com/" + repoName + "/%s"
	// map of all categories and respective TILs
	tilsMap := make(map[string][]Til)
	// list of all (non-grouped by category) TILs for use with `list_most_recent` feature
	var tilsSlice []Til
	// tils = TIL markdown files
	tils, _ := filepath.Glob(repoPath + "/**/*.md")

	for _, til := range tils {
		// grab the "category" and the "fileName"
		// ex: html/div-tags.md -- category "html" fileName "div-tags.md"
		splitResult := strings.Split(til, "/")
		length := len(splitResult)
		category := splitResult[length-2]
		fileName := splitResult[length-1]

		if strings.ToLower(fileName) == "readme.md" {
			continue
		}

		// Read the first line of each fileName, use the string as a title
		f, err := os.Open(til)
		if err != nil {
			log.Panic(err)
		}
		reader := bufio.NewReader(f)
		linkTitle, err := reader.ReadString('\n')
		if err != nil {
			log.Println(fmt.Sprintf("ERROR: file \"%s\" does not have > 1 line of text (no title)", fileName))
			log.Panic(err)
		}

		// strip "# " from beginning of line
		linkTitle = re2.ReplaceAllString(linkTitle, "")
		linkTitle = strings.TrimSpace(linkTitle)

		var absoluteLink = fmt.Sprintf(absoluteLinkFormat, category+"/"+fileName)

		// if category first encountered in loop so far, append new map key, otherwise
		// add to existing
		tilStruct := Til{
			Title:        linkTitle,
			AbsoluteLink: absoluteLink,
			Category:     category,
			DateAdded:    cmdGetDate(category + "/" + fileName),
		}

		if _, exists := tilsMap[category]; exists {
			tilsMap[category] = append(tilsMap[category], tilStruct)
		} else {
			tilsMap[category] = []Til{tilStruct}
		}

		tilsSlice = append(tilsSlice, tilStruct)
	}

	n, err := strconv.Atoi(inputListMostRecent)
	if err != nil {
		n = 0
	}

	cmdTrimMostRecentTils(&tilsSlice, n)

	// load and execute template, write results to README.md
	t, err := template.New(path.Base(templatePath)).ParseFiles(templatePath)
	if err != nil {
		log.Panic(err)
	}

	var output bytes.Buffer
	err = t.Execute(&output, struct {
		Tils                   map[string][]Til
		AllTils                []string
		InputDescription       string
		InputFooter            string
		MostRecentTils         []Til
		InputDateFormat        string
		InputTilsCounterFormat string
	}{
		Tils:                   tilsMap,
		AllTils:                tils,
		InputDescription:       inputDescription,
		InputFooter:            inputFooter,
		MostRecentTils:         tilsSlice,
		InputDateFormat:        inputDateFormat,
		InputTilsCounterFormat: inputTilsCounterFormat,
	})

	if err != nil {
		log.Panic(err)
	}

	fmt.Println("------------------------------------------------------------")
	fmt.Println("------------------------------------------------------------")
	fmt.Print(output.String())
	fmt.Println("------------------------------------------------------------")
	fmt.Println("------------------------------------------------------------")

	// truncates before writing
	ioutil.WriteFile(repoPath+"/README.md", []byte(output.String()), 0644)
}
