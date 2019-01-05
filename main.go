package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"os"
	"regexp"
	"strconv"
)

var matchGroup = regexp.MustCompile(`\$(\$|\d+)`)

func main() {
	if len(os.Args) < 3 {
		help()
	}
	pattern := os.Args[1]
	format := os.Args[2]
	if len(os.Args) > 3 {
		pattern = "(?" + os.Args[3] + ")" + pattern
	}

	regex, err := regexp.Compile(pattern)
	if err != nil {
		die("Cannot compile regex:", err)
	}

	w := bufio.NewWriter(os.Stdout)
	defer w.Flush()

	input, _ := ioutil.ReadAll(os.Stdin)
	for _, match := range regex.FindAllSubmatch(input, -1) {
		w.WriteString(
			matchGroup.ReplaceAllStringFunc(format, func(s string) string {
				if s == "$$" {
					return "$"
				}
				idx, _ := strconv.Atoi(s[1:]) // ignore leading $
				if idx >= len(match) {
					return ""
				}
				return string(match[idx])
			}))
		w.WriteString("\n")
	}
}

func help() {
	os.Stdout.WriteString("Usage: rgx <pattern> <format> [flags]\n")
	os.Exit(1)
}

func die(a ...interface{}) {
	for _, x := range a {
		fmt.Fprintln(os.Stderr, x)
	}
	os.Exit(1)
}
