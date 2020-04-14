#!/bin/sh
#
# Create a directory and generate template files for a given day. Ex:
#
#   ./wake-up.sh 01

set -e

ROOT=$(dirname "$0")

if [ $# -ne 1 ]; then
    echo "usage $(basename "$0") DAY" > /dev/stderr
    exit 1
fi

DAY=$1

DAYDIRNAME="day${DAY}"
DAYDIRPATH="${ROOT}/${DAYDIRNAME}"
MAINFILENAME="main.go"
TESTFILENAME="main_test.go"
MAINFILEPATH="${DAYDIRPATH}/${MAINFILENAME}"
TESTFILEPATH="${DAYDIRPATH}/${TESTFILENAME}"

if [ -d "${DAYDIRPATH}" ]; then
    echo "${DAYDIRPATH} already exists" > /dev/stderr
    exit 1
fi
mkdir -p "$DAYDIRPATH"

# input.txt
: > "${DAYDIRPATH}/input.txt"

# README.md
: > "${DAYDIRPATH}/README.md"

# README.part2.md
: > "${DAYDIRPATH}/README.part2.md"

# answer.md
cat <<EOF > "${DAYDIRPATH}/answer.md"
Your puzzle answer was \`?\`.
EOF

# answer.part2.md
cat <<EOF > "${DAYDIRPATH}/answer.part2.md"
Your puzzle answer was \`?\`.
EOF

# Main file
cat <<EOF > "$MAINFILEPATH"
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
)

// TODO
func main() {
	lines, err := Parse(os.Stdin)
	if err != nil {
		log.Fatalf("input error: %s\n", err)
	}
	fmt.Print(lines)
}

// TODO
func Parse(r io.Reader) ([]string, error) {
	var lines []string
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return lines, nil
}
EOF

# Test file
cat <<EOF > "$TESTFILEPATH"
package main

import "testing"

func TestMain(t *testing.T) {
	// TODO
}
EOF
