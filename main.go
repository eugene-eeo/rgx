package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"regexp"
	"strconv"
)

const BUFLEN = 1024 * 256

var matchGroup = regexp.MustCompile(`\$(\$|\d+)`)

func formatMatch(buf []byte, indices []int, format []byte, dst io.Writer) {
	length := len(indices) / 2
	dst.Write(
		matchGroup.ReplaceAllFunc(format, func(b []byte) []byte {
			s := string(b)
			if s == "$$" {
				return []byte("$")
			}
			idx, _ := strconv.Atoi(s[1:]) // ignore leading $
			if idx >= length {
				return []byte{}
			}
			i := idx * 2
			return buf[indices[i]:indices[i+1]]
		}))
	dst.Write([]byte("\n"))
}

func max(a, b int) int {
	if a < b {
		return b
	}
	return a
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func main() {
	if len(os.Args) < 3 {
		help()
	}
	flags := "ms"
	if len(os.Args) > 3 {
		flags = os.Args[3]
	}
	// buffer size
	buflen := BUFLEN
	if len(os.Args) > 4 {
		bl, err := strconv.Atoi(os.Args[4])
		if err != nil {
			die("Cannot parse buffer length")
		}
		buflen = bl
		if buflen <= 1 {
			buflen = 1024
		}
	}
	pattern := "(?" + flags + ")" + os.Args[1]
	format := []byte(os.Args[2])

	regex, err := regexp.Compile(pattern)
	if err != nil {
		die("Cannot compile regex:", err)
	}

	r := bufio.NewReader(os.Stdin)
	w := bufio.NewWriter(os.Stdout)
	b := make([]byte, buflen)
	s := 0
	defer w.Flush()

	for {
		n, err := r.Read(b[s:])
		if err != nil && err != io.EOF {
			break
		}
		indices := regex.FindAllSubmatchIndex(b[:s+n], -1)
		for _, idxs := range indices {
			formatMatch(b, idxs, format, w)
		}
		if err == io.EOF {
			break
		}
		// if we don't have any matches then
		s := max(min((s+n)/8, 256), 1)
		if len(indices) > 0 {
			s = s + n - indices[len(indices)-1][1]
		}
		// we may need to copy some unmatched characters
		// over to the new buffer
		if s > 0 {
			copy(b, b[s:])
		}
	}
}

func help() {
	os.Stdout.WriteString("Usage: rgx <pattern> <format> [<flags>]\n")
	os.Exit(1)
}

func die(a ...interface{}) {
	for _, x := range a {
		fmt.Fprintln(os.Stderr, x)
	}
	os.Exit(1)
}
