package main

import (
	"bufio"
	"fmt"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
)

type grid struct {
	grid [][]string
	size int
}

const black = "__black__"

func (g *grid) String() string {
	out := "\n\n"

	for i := range g.grid {
		for j := range g.grid[i] {
			out += g.grid[i][j]
			out += "\t"
		}
		out += "\n"
	}

	return out
}

func (g *grid) FillAt(row, col int, word string, isHorizontally bool) error {

	if isHorizontally {
		// TODO: Iterate based on character + diacritic
		for i := 0; i < len(word); i++ {
			if i >= g.size {
				return fmt.Errorf("Overflow: %d, %d, %d, %s, %s", i, row, col, word, g)
			}
			log.Println(row, col, i, word)
			g.grid[row][i+col] = word[i : i+1]
			log.Println(row, i)
		}
	}

	return nil
}

func newGrid(gridSize int) *grid {
	g := make([][]string, 0)
	for i := 0; i < gridSize; i++ {
		g = append(g, make([]string, 0))
		for j := 0; j < gridSize; j++ {
			g[i] = append(g[i], "?")
		}
	}

	return &grid{grid: g, size: gridSize}
}

type words struct {
	words     []string
	usedWords []string

	m   map[int][]string
	rnd *rand.Rand
}

func newWords(path string) (*words, error) {

	m := make(map[int][]string)

	file, err := os.Open(path)
	if err != nil {
		//handle error
		return nil, err
	}
	defer file.Close()

	lines := make([]string, 0)

	s := bufio.NewScanner(file)
	for s.Scan() {
		line := strings.TrimSpace(s.Text())
		if line == "" {
			continue
		}
		lines = append(lines, line)

		t := len(line)

		// if t == 4 {
		// 	log.Println(t, line, m[t])
		// }

		if m[t] == nil {
			m[t] = make([]string, 0)
		}

		m[t] = append(m[t], line)

		if false {
			if t == 4 {
				log.Println("-------")
				for _, i := range m[t] {
					log.Println(i)
				}
				log.Println("-------\n\n")
			}
		}
	}

	s1 := rand.NewSource(time.Now().UnixNano())
	r1 := rand.New(s1)

	return &words{
		words:     lines,
		usedWords: []string{},
		m:         m,
		rnd:       r1}, nil
}

// thread unsafe
func (w *words) randWord(prefix string, minLen, maxLen int) (string, error) {
	strLen := w.rnd.Intn((maxLen-minLen)+1) + minLen
	words := w.m[strLen]
	// log.Println(strLen, words[1])
	// log.Println(words)

	for i := 0; i < len(words); i++ {
		alreadyUsed := false
		for j := 0; j < len(w.usedWords); j++ {
			if w.usedWords[j] == words[i] {
				alreadyUsed = true
				break
			}
		}
		// log.Println(alreadyUsed, i, words[i])

		if !alreadyUsed {
			w.usedWords = append(w.usedWords, words[i])
			return words[i], nil
		}
	}

	return "", fmt.Errorf("Insufficient words")
}

func main() {

	log.SetFlags(log.LstdFlags | log.Lshortfile)

	size := 8
	if size < 8 {
		log.Fatal("We need at least a 8x8 grid to generate crossword")
		return
	}

	words, err := newWords("./words_alpha.txt")
	if err != nil {
		log.Fatal(err)
		return
	}
	// log.Println(words.m[4])

	grid := newGrid(size)
	log.Println(grid)

	row := 0
	col := 0
	// Fill first cell
	randWord, err := words.randWord("", size/2, size/2)
	if err != nil {
		log.Fatal(err)
		return
	}

	err = grid.FillAt(row, col, randWord, true)
	if err != nil {
		log.Fatal(err)
		return
	}

	col += (size / 2)

	log.Println(grid)
}
