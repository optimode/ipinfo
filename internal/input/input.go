package input

import (
	"bufio"
	"io"
	"strings"
)

// FromReader reads IPs from a reader, one per line.
// Empty lines and lines starting with # are skipped.
func FromReader(r io.Reader) []string {
	var ips []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		line = strings.TrimRight(line, "\r")
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		ips = append(ips, line)
	}
	return ips
}
