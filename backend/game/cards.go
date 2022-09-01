package game

import (
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
	go periodicallyLoadCards()
}

var cards []string

func getCardsList() ([]string, error) {
	data, err := os.ReadFile("cards.txt")
	if err != nil {
		return nil, err
	}

	list := strings.Split(string(data), "\n")

	return list[:len(list)-1], nil
}

func periodicallyLoadCards() {
	for {
		l, err := getCardsList()
		if err != nil {
			if cards == nil {
				panic(err)
			} else {
				log.Printf("the list of cards could not be reloaded")
			}
		} else {
			if cards == nil {
				log.Printf("the list of cards has been loaded")
			} else {
				log.Printf("the list of cards has been reloaded")
			}
			cards = l
		}

		time.Sleep(time.Hour)
	}
}

func Shuffle(c []string) []string {
	rand.Shuffle(len(c), func(i, j int) { c[i], c[j] = c[j], c[i] })
	return c
}

func RandomIntSet(n, upper int) []int {
	if n > upper {
		panic("the set size is too big for the upper bound")
	}

	s := make([]int, 0, n)
	m := map[int]bool{}

	for len(s) < n {
		r := rand.Intn(upper)

		if in := m[r]; !in {
			m[r] = true
			s = append(s, r)
		}
	}

	return s
}

func RandomCards(n int) []string {
	list := cards

	if len(list) < n {
		panic("not enough cards")
	}

	cardsI := RandomIntSet(n, len(list))

	c := make([]string, 0, n)
	for _, i := range cardsI {
		c = append(c, list[i])
	}

	return c
}

func partitionCards(cards []string, n int) [][]string {
	l := len(cards)
	partitions := make([][]string, 0, n)

	for i := 0; i < n; i++ {
		p := make([]string, 0, l/(n-i))

		for cap(p) > len(p) {
			p = append(p, cards[l-1])
			l--
		}

		partitions = append(partitions, p)
	}

	return partitions
}
