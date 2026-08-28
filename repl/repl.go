package repl

import "strings"

func Split(input string) []string {
    return strings.Fields(strings.ToLower(input))
}
