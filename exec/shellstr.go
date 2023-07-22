package exec

import (
	"fmt"
	"regexp"
	"strings"
)

func ShellStr2List(command string) []string {
	command = strings.Replace(command, "\\\n", "", -1)
	commandParts := strings.Split(command, " ")
	reassembledCommandParts := []string{}

	commandPartBuffer := ""
	commandPartBufferQuoteChar := ""

	// the following reconstructs some command line pieces for quoted args with spaces
	// which is a case our simple command line part split doesn't handle above
	// this is very imperfect, but will do for us here until it doesn't

	partSearch, _ := regexp.Compile(`['"]`)

	for _, commandPart := range commandParts {
		matches := partSearch.FindAllStringIndex(commandPart, -1)
		if len(matches) == 1 {
			// we're either starting or ending a buffered command part
			if commandPartBuffer == "" {
				// starting a buffer
				if strings.Contains(commandPart, `"`) {
					commandPartBufferQuoteChar = `"`
				} else {
					commandPartBufferQuoteChar = `'`
				}
				commandPartBuffer = commandPart
			} else if strings.Contains(commandPart, commandPartBufferQuoteChar) {
				// finishing a buffer
				commandPartBuffer = fmt.Sprintf(`%s %s`, commandPartBuffer, commandPart)
				reassembledCommandParts = append(reassembledCommandParts,
					strings.Replace(commandPartBuffer, commandPartBufferQuoteChar, "", -1))
				commandPartBuffer = ""
				commandPartBufferQuoteChar = ""
			}
		} else if commandPartBuffer != "" {
			// in the middle of a re-assemble
			commandPartBuffer = fmt.Sprintf("%s %s", commandPartBuffer, commandPart)
		} else {
			// just a normal re-append
			reassembledCommandParts = append(reassembledCommandParts, commandPart)
		}
	}

	return reassembledCommandParts
}
