package words

import (
	"bytes"
	_ "embed"
	"encoding/csv"
	"errors"
	"io"
	"math/rand"
	"strings"
	"time"
)

//go:embed words.txt
var wordlist string

//go:embed quotes.csv
var quoteBytes []byte

type Quote struct {
	Author string
	Text   string
}

type Stats struct {
	RawWPM float64
	NetWPM float64
	Acc    float64
}

func GetRandomWords(num int) ([]string, error) {
	words := strings.Split(wordlist, "\n")
	rand.Shuffle(len(words), func(i, j int) {
		words[i], words[j] = words[j], words[i]
	})
	return words[:num], nil
}

func GetRandomQuote() (Quote, error) {
	r := csv.NewReader(bytes.NewReader(quoteBytes))

	if _, err := r.Read(); err != nil && err != io.EOF {
		return Quote{}, err
	}

	var quotes []Quote
	for {
		record, err := r.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			return Quote{}, err
		}
		if len(record) >= 2 {
			quotes = append(quotes, Quote{
				Author: record[0],
				Text:   record[1],
			})
		}
	}

	if len(quotes) == 0 {
		return Quote{}, errors.New("No Quotes in the file")
	}

	return quotes[rand.Intn(len(quotes))], nil
}

func GetStats(totalTyped int, mistakes int, elapsed time.Duration) (Stats, error) {
	if mistakes < 0 || mistakes > totalTyped {
		return Stats{}, errors.New("More mistakes than total typed")
	}

	if elapsed <= 0 {
		return Stats{}, errors.New("elapsed time must be > 0")
	}

	if totalTyped == 0 {
		return Stats{RawWPM: 0, NetWPM: 0, Acc: 0}, nil
	}

	rawWPM := float64(totalTyped/5) / elapsed.Seconds()
	netWPM := rawWPM - (float64(mistakes) / elapsed.Seconds())
	acc := float64(totalTyped-mistakes) / float64(totalTyped)

	return Stats{RawWPM: rawWPM, NetWPM: netWPM, Acc: acc}, nil
}

