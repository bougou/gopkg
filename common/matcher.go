package common

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
)

type Matcher struct {
	regExp *regexp.Regexp
}

func NewMatcher(regExp *regexp.Regexp) *Matcher {
	return &Matcher{
		regExp: regExp,
	}
}

func (m Matcher) MatchAndExtract(s string, matchGroupSeq int) (matched bool, extracted string, err error) {
	if m.regExp == nil {
		return false, "", fmt.Errorf("nil regExp passed")
	}

	if m.regExp.MatchString(s) {
		matchedStrs := m.regExp.FindStringSubmatch(s)
		if len(matchedStrs) < matchGroupSeq+1 {
			return true, "", fmt.Errorf("not enough matched results")
		}
		return true, matchedStrs[matchGroupSeq], nil
	}

	return false, "", nil
}

func (m Matcher) MatchAndExtractFromLines(lines string, matchGroupSeq int) (string, error) {
	scannner := bufio.NewScanner(bytes.NewBuffer([]byte(lines)))
	scannner.Split(scanLines)
	for scannner.Scan() {
		line := scannner.Text()
		lineText := strings.TrimSpace(line)
		if lineText == "" {
			continue
		}

		matched, extracted, err := m.MatchAndExtract(lineText, matchGroupSeq)
		if err != nil {
			return "", err
		}
		if matched {
			return extracted, nil
		}
	}

	return "", nil
}
