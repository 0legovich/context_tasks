package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"sort"
	"strconv"
	"strings"
)

func main() {
	file, err := os.Open("input.txt")
	defer file.Close()

	if err != nil {
		log.Fatal(err)
	}

	target, sequence := fromFile(file)
	fmt.Println(sequence)

	sort.Ints(sequence)

	var first, s int
	last := len(sequence) - 1

	for first < last {
		if sequence[first] > target {
			printInFile("0")
			return
		}

		if sequence[last] > target {
			last--
			continue
		}

		s = sequence[first] + sequence[last]

		if s == target {
			printInFile("1")
			return
		}

		if s < target {
			first++
		} else {
			last--
		}
	}

	printInFile("0")
}

// reading strings batches considering gluing broken items
func fromFile(file io.Reader) (int, []int) {
	batch := make([]byte, 10000)

	readBytes, err := file.Read(batch); if err != nil {
		log.Fatal(err)
	}

	splitLines := make([]string, 0)
	sequence := make([]int, 0)

	var targetStr, prevElem string
	var target int
	var lastPrevByte, lastByte byte

	for i := 0; readBytes != 0; i++ {
		splitLines = strings.Split(string(batch), "\n")

		if len(splitLines) >= 2 && i == 0 {
			targetStr = splitLines[0]
			target, err = strconv.Atoi(targetStr)
			if err != nil {
				log.Fatal(err)
			}

			if len(splitLines[1]) > 0 {
				lastByte = fillSequence(&sequence, splitLines[1], target, prevElem, lastPrevByte)
				lastPrevByte = lastByte
			}
		} else {
			if len(splitLines[0]) > 0 {
				lastByte = fillSequence(&sequence, splitLines[0], target, prevElem, lastPrevByte)
				lastPrevByte = lastByte
			}
		}

		readBytes, _ = file.Read(batch)
	}

	return target, sequence
}

func fillSequence(sequence *[]int, line string, target int, prevElem string, lastPrevByte byte) byte {
	seqStr := strings.Split(line, " ")
	firstByte := line[0]
	lastByte := line[len(line)-1]

	for idx, elem := range seqStr {
		if firstByte != 32 && idx == 0 && lastPrevByte != 32 {
			elem = prevElem + elem
		}

		if firstByte == 32 && idx == 0 && lastPrevByte != 32 {
			prevElemInt, err := strconv.Atoi(prevElem)
			if err == nil {
				*sequence = append(*sequence, prevElemInt)
			}
		}

		if lastByte != 32 && idx == len(seqStr)-1 {
			prevElem = elem
			continue
		}

		elemInt, err := strconv.Atoi(elem)
		if err != nil {
			continue
		}

		if elemInt <= target {
			*sequence = append(*sequence, elemInt)
		}
	}

	return lastByte
}

func printInFile(line string) {
	file, err := os.OpenFile("output.txt", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
	defer file.Close()

	if err != nil {
		return
	}
	_, err = fmt.Fprintln(file, line)
}
