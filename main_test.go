package main

import (
	"io/ioutil"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setup() {
	inputDescription = "# TIL\n> Today I Learned\n\n\nA collection of concrete writeups of small things I learn daily.\n\n\n"
	inputFooter = ""
	inputTilsCounterFormat = "_%d TILs and counting..._"
	inputRepoName = "user/repository"
	inputPresentationType = "list"
	listPath = "./README.md.tmpl"
	tablePath = "./README_TABLE.md.tmpl"
}

func TestOneTil(t *testing.T) {
	setup()
	repoPath = "./test_data/1_til"
	main()
	expected, err := ioutil.ReadFile(repoPath + "/README.md.expected")
	if err != nil {
		t.Error(err)
	}
	actual, err := ioutil.ReadFile(repoPath + "/README.md")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, string(expected), string(actual))
}

func TestZeroTil(t *testing.T) {
	setup()
	repoPath = "./test_data/zero_til"
	main()
	expected, err := ioutil.ReadFile(repoPath + "/README.md.expected")
	if err != nil {
		t.Error(err)
	}
	actual, err := ioutil.ReadFile(repoPath + "/README.md")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, string(expected), string(actual))
}

func TestManyTil(t *testing.T) {
	setup()
	repoPath = "./test_data/many_til"
	main()
	expected, err := ioutil.ReadFile(repoPath + "/README.md.expected")
	if err != nil {
		t.Error(err)
	}
	actual, err := ioutil.ReadFile(repoPath + "/README.md")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, string(expected), string(actual))
}

func TestOneTilInputs(t *testing.T) {
	setup()
	repoPath = "./test_data/many_with_inputs"
	inputDescription = "This is a placeholder description used for testing.\n\n"
	inputFooter = "here is where the markdown footer links would go"
	main()
	expected, err := ioutil.ReadFile(repoPath + "/README.md.expected")
	if err != nil {
		t.Error(err)
	}
	actual, err := ioutil.ReadFile(repoPath + "/README.md")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, string(expected), string(actual))
}

func TestOneTilInputsAndMostRecent(t *testing.T) {
	setup()
	repoPath = "./test_data/many_with_inputs_and_most_recent"
	inputDescription = "# TIL\n> Today I Learned\n\nThis is a placeholder description used for testing.\n\n"
	inputFooter = "here is where the markdown footer links would go"
	inputListMostRecent = "3"
	inputDateFormat = time.RFC822
	main()
	expected, err := ioutil.ReadFile(repoPath + "/README.md.expected")
	if err != nil {
		t.Error(err)
	}
	actual, err := ioutil.ReadFile(repoPath + "/README.md")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, string(expected), string(actual))
}

func TestNoCounterLabel(t *testing.T) {
	setup()
	repoPath = "./test_data/no_counter_label"
	inputListMostRecent = "0"
	inputTilsCounterFormat = ""
	main()
	expected, err := ioutil.ReadFile(repoPath + "/README.md.expected")
	if err != nil {
		t.Error(err)
	}
	actual, err := ioutil.ReadFile(repoPath + "/README.md")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, string(expected), string(actual))

}
func TestTablePresentationManyTil(t *testing.T) {
	setup()
	repoPath = "./test_data/table_presentation_many_til"
	inputTilsCounterFormat = ""
	inputDescription = "Header\n"
	inputPresentationType = "table"
	inputListMostRecent = "5"
	inputDateFormat = "02/01/2006"
	main()
	expected, err := ioutil.ReadFile(repoPath + "/README.md.expected")
	if err != nil {
		t.Error(err)
	}
	actual, err := ioutil.ReadFile(repoPath + "/README.md")
	if err != nil {
		t.Error(err)
	}
	assert.Equal(t, string(expected), string(actual))
}
