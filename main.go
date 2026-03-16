package main

import (
	"bufio"
	"fmt"
	"math/rand"
	"memorize/memorize"
	"os"
	"strconv"
	"strings"

	"github.com/fatih/color"
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

	red := color.New(color.FgRed).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()
	bold := color.New(color.Bold).SprintFunc()

	cfg, err := memorize.LoadConfig("config.toml")
	if err != nil {
		fmt.Println("Warning: failed to load config, using defaults:", err)
	} else {
		memorize.SetConfig(cfg)
	}

	args := os.Args
	if len(args) < 2 {
		fmt.Println("usage: memorize <wordlist.json>")
		os.Exit(1)
	}
	const cachePath = "deck.cache.json"

	var deck *memorize.Deck

	// Try loading saved deck
	deck, err = memorize.LoadDeckCache(cachePath)
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

		fmt.Println(bold("Prompt:"), prompt.Content)
		fmt.Println(bold("Options:"))

		for i, o := range opts {
			fmt.Printf("%d) %s\n", i+1, maskWord(o))
		}

		fmt.Print("> ")

		line, _ := reader.ReadString('\n')
		line = strings.TrimSpace(line)

		var q int
		edit, score := memorize.EditDiff(line, word.Word)

		if n, err := strconv.Atoi(line); err == nil {
			q = n
		} else {
			q = memorize.QFromSimilarity(1 - float64(score)/float64(len(line)+len(word.Word)))
		}
		if q == 5 {
			fmt.Println(bold("Answer:"), green(word.Word))
			fmt.Println(bold("Score:"), green(q))
		} else {
			fmt.Println(bold("Answer:"), edit, " -> ", word.Word)
			fmt.Println(bold("Score:"), red(q))
		}
		fmt.Println(bold("Hint:"), *prompt.Hint)
		fmt.Println("----------------")

		deck.Review(q)
		if err := deck.SaveDeckCache(cachePath); err != nil {
			fmt.Println("Warning: failed to save deck cache:", err)
		}
	}
}
