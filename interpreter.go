package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"unicode/utf8"

	"golang.org/x/term"
)

func interpreter() {
	var bracketPositions []int64
	var tape []int
	var pointer int

	pointer = 0
	tape = append(tape, 0)

	file, err := os.Open(filepath)

	if err != nil {
		panic(err)
	}

	defer file.Close()

	for {
		r, err := readRuneFromFile(file)

		if err != nil {
			break
		}

		switch r {
		case '[':
			if tape[pointer] == 0 {
				seekToLoopEnd(file)
			} else {
				seek, err := file.Seek(0, io.SeekCurrent)
				if err != nil {
					panic(err)
				}
				bracketPositions = append(bracketPositions, seek)
			}
		case ']':
			if len(bracketPositions) == 0 {
				panic("Unmatched ']'")
			}
			if tape[pointer] != 0 {
				seekPos := bracketPositions[len(bracketPositions)-1]
				_, err := file.Seek(seekPos, io.SeekStart)
				if err != nil {
					panic(err)
				}
			} else {
				bracketPositions = bracketPositions[:len(bracketPositions)-1]
			}
		case '+':
			tape[pointer]++
		case '-':
			tape[pointer]--
		case '<':
			if pointer == 0 {
				panic("Can't acces tape below 0")
			}
			pointer--
		case '>':
			pointer++
			if pointer >= len(tape) {
				tape = append(tape, 0)
			}
		case '.':
			fmt.Printf("%c", rune(tape[pointer]))
		case ',':
			tape[pointer] = int(readRune())
		}
	}

	fmt.Println()
}

func seekToLoopEnd(file *os.File) {
	var loops int = 1
	for {
		r, err := readRuneFromFile(file)

		if err != nil {
			break
		}

		switch r {
		case '[':
			loops++
		case ']':
			loops--
		}
		if loops == 0 {
			return
		}
	}
}

func readRuneFromFile(file *os.File) (rune, error) {
	var buf [4]byte
	var runeBytes []byte

	for i := 0; i < 4; i++ {
		_, err := file.Read(buf[i : i+1])
		if err != nil {
			return 0, err
		}
		runeBytes = buf[:i+1]
		r, size := utf8.DecodeRune(runeBytes)
		if r != utf8.RuneError || size > 1 {
			return r, nil
		}
	}
	return utf8.RuneError, nil
}

func readRune() rune {
	if term.IsTerminal(int(os.Stdin.Fd())) {
		oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}
		defer term.Restore(int(os.Stdin.Fd()), oldState)

		var b [4]byte
		n, err := os.Stdin.Read(b[:1])
		if err != nil || n == 0 {
			return 0
		}
		r, size := utf8.DecodeRune(b[:1])
		if size == 1 && r == utf8.RuneError {
			n, err := os.Stdin.Read(b[1:4])
			if err != nil {
				return r
			}
			r, _ = utf8.DecodeRune(b[:1+n])
		}
		return r

	} else {
		reader := bufio.NewReader(os.Stdin)
		r, _, err := reader.ReadRune()
		if err != nil {
			return 0
		}
		return r
	}
}
