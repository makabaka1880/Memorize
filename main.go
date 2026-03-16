package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"memorize/memorize"
	"os"
	"strconv"
	"strings"
)

func maskWord(w string) string {
	vowels := "aeiouAEIOU"
	var out []string

	for _, r := range w {
		if strings.ContainsRune(vowels, r) {
			out = append(out, "_")
		} else {
			out = append(out, string(r))
		}
	}

	return strings.Join(out, "")
}

func choices(deck *memorize.Deck, correct string) []string {
	var opts []string
	opts = append(opts, correct)

	for len(opts) < 3 {
		i := rand.Intn(len(deck.Queue()))
		w := deck.Queue()[i].Word.Word

		if w != correct {
			duplicate := false
			for _, o := range opts {
				if o == w {
					duplicate = true
				}
			}

			if !duplicate {
				opts = append(opts, w)
			}
		}
	}

	rand.Shuffle(len(opts), func(i, j int) {
		opts[i], opts[j] = opts[j], opts[i]
	})

	return opts
}

func main() {
	// rand.Seed(time.Now().UnixNano())

	args := os.Args
	if len(args) < 2 {
		fmt.Println("usage: memorize <wordlist.json>")
		os.Exit(1)
	}
	const cachePath = "deck.cache.json"

	var deck *memorize.Deck

	// Try loading saved deck
	deck, err := memorize.LoadDeckCache(cachePath)
	if err != nil {
		// fallback: read wordlist and create new deck
		list := memorize.ReadWordList(args[1])
		deck = memorize.NewDeck(list)
		fmt.Println("Starting new deck.")
	} else {
		fmt.Println("Loaded deck from cache.")
	}

	reader := bufio.NewReader(os.Stdin)

	for {
		word := deck.Next()
		if word == nil {
			fmt.Println("done!")
			return
		}

		opts := choices(deck, word.Word)

		prompt := word.Prompts[rand.Intn(len(word.Prompts))]

		fmt.Println("Prompt:", prompt.Content)
		fmt.Println("Options:")

		for i, o := range opts {
			fmt.Printf("%d) %s\n", i+1, maskWord(o))
		}

		fmt.Print("> ")

		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		var q int

		if n, err := strconv.Atoi(line); err == nil {
			q = n

		} else if line == "" {
			q = 0

		} else if strings.EqualFold(line, word.Word) {
			q = 5

		} else if strings.Contains(strings.ToLower(word.Word), strings.ToLower(line)) {
			q = 4

		} else {
			q = 2
		}

		fmt.Println("Answer:", word.Word)
		fmt.Println("Score:", q)
		fmt.Println("Hint:", *prompt.Hint)
		fmt.Println("----------------")

		deck.Review(q)
		if err := deck.SaveDeckCache(cachePath); err != nil {
			fmt.Println("Warning: failed to save deck cache:", err)
		}
	}
}
