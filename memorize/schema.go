package memorize

import (
	"bytes"
	"encoding/json"
	"fmt"
	"os"
)

type WordList struct {
	Name  string `json:"name"`
	Words []Word `json:"words"`
}

type Word struct {
	Word    string   `json:"word"`
	Prompts []Prompt `json:"prompts"`
}

type Prompt struct {
	Content string  `json:"content"`
	Hint    *string `json:"hint"`
}

func ReadWordList(path string) WordList {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading file: %v\n", err)
		os.Exit(1)
	}

	var wordList WordList
	fmt.Fprintf(os.Stdout, "%v", path)
	if err := json.Unmarshal(bytes.TrimPrefix(data, []byte("\xef\xbb\xbf")), &wordList); err != nil {
		fmt.Fprintf(os.Stderr, "Error parsing JSON: %v\n", err)
		os.Exit(1)
	}

	return wordList
}
