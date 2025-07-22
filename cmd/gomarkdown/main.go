package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/Powerisinschool/gomarkdown/internal/parser"
	"github.com/Powerisinschool/gomarkdown/internal/renderer"
	"os"
)

func main() {
	input := []byte("")
	if len(os.Args) < 2 {
		//fmt.Fprintln(os.Stderr, "Usage: gomarkdown <string>")
		//os.Exit(1)
		// Create a new scanner to read from standard input.
		scanner := bufio.NewScanner(os.Stdin)

		// Iterate through each line of input.
		for scanner.Scan() {
			line := scanner.Text() // Get the current line as a string.
			//fmt.Printf("Received: %s\n", line)
			input = append(input, line...)
			input = append(input, '\n')
		}

		// Check for any errors that occurred during scanning.
		if err := scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading from stdin: %v\n", err)
		}
	} else {
		input = []byte(os.Args[1])
	}

	p := parser.New(input)
	doc := p.Parse()

	var buf bytes.Buffer
	r := renderer.New(&buf)
	if err := r.Render(doc); err != nil {
		fmt.Fprintf(os.Stderr, "error rendering markdown: %v\n", err)
		os.Exit(1)
	}

	fmt.Print(buf.String())
}
