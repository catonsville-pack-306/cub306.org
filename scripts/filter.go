package main


import (
    "fmt"
    "bufio"
    "flag"
    "io"
    "os"
    "strings"
    )

func filter (input io.Reader, output *io.PipeWriter, match string) {
	scanner := bufio.NewScanner(input)
	defer output.Close()
	for scanner.Scan() {
		line := scanner.Text()
		if len(match)>0 && strings.Contains(line, match) {
			continue
		}
		output.Write([]byte(line + "\n"))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}

func main() {
    match := flag.String("match", "", "text to match in each line")
    flag.Parse()

	piper, pipew := io.Pipe()
	go filter(os.Stdin, pipew, *match)
	io.Copy(os.Stdout, piper)
}