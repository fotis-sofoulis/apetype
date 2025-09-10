package words

import (
	_ "embed"
	"math/rand"
	"strings"
)

//go:embed words.txt
var wordlist string

func GetRandomWords() ([]string, error) {
	words := strings.Split(wordlist, "\n")
	rand.Shuffle(len(words), func(i, j int) {
        words[i], words[j] = words[j], words[i]
    })
	return words[:25], nil
}
