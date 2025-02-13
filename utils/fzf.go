package utils

import (
	"log"

	"github.com/koki-develop/go-fzf"
)

func FuzzySearch(prompt string, items []string) []int {
	f, err := fzf.New(fzf.WithPrompt(prompt), fzf.WithStyles(fzf.WithStylePrompt(fzf.Style{Bold: true})))
	if err != nil {
		log.Fatal(err)
	}

	idxs, err := f.Find(items, func(i int) string { return items[i] })
	if err != nil {
		log.Fatal(err)
	}

	return idxs
}
