package common

import (
	"bufio"
	"bytes"
	"fmt"
	"regexp"
	"strings"
	"testing"
)

func Test_Dirquota(t *testing.T) {
	output := `
# file: ewn7jf9wc3rf
abc ceph.quota.max_bytes="52428800000" efg ceph.quota.max_bytes="1234" efg ceph.quota.max_bytes="5678"

`

	scannner := bufio.NewScanner(bytes.NewBuffer([]byte(output)))
	scannner.Split(scanLines)

	var d int
	for scannner.Scan() {
		d++
		lineText := strings.TrimSpace(scannner.Text())
		if lineText == "" {
			continue
		}

		matcher := regexp.MustCompile(`ceph.quota.([a-z_]+)="(\d+)"`)

		if matcher.Match([]byte(lineText)) {
			fmt.Printf("line #%d, %s\n", d, lineText)

			var s string

			bl := matcher.Find([]byte(lineText))
			fmt.Println("1 bl.", bl)
			fi := matcher.FindIndex([]byte(lineText))
			fmt.Println("1 bl.", fi)
			s = matcher.FindString(lineText)
			fmt.Println("1 s. FindString", s)

			dd := matcher.FindSubmatch([]byte(lineText))
			fmt.Println("dd", dd, len(dd))
			fmt.Println("dd", string(dd[0]),
				string(dd[1]), string(dd[2]))

			di := matcher.FindSubmatchIndex([]byte(lineText))
			fmt.Println("di", di, len(di))
			fmt.Println("di", (di[0]),
				(di[1]), (di[2]))

			s2 := matcher.FindString(lineText)
			fmt.Println("s2:", s2)

			s3 := matcher.FindStringIndex(lineText)
			fmt.Println("s3:", s3)

			s4 := matcher.FindStringSubmatch(lineText)
			fmt.Println("s4:", s4)

			s5 := matcher.FindStringSubmatchIndex(lineText)
			fmt.Println("s5:", s5)

			s13 := matcher.FindAllString(lineText, -1)
			fmt.Println("s13:", strings.Join(s13, "||"))

			s14 := matcher.FindAllStringIndex(lineText, -1)
			fmt.Println("s14:", s14)

			s15 := matcher.FindAllStringSubmatch(lineText, -1)
			fmt.Println("s15:", s15)

			s16 := matcher.FindAllStringSubmatchIndex(lineText, -1)
			fmt.Println("s16:", s16)
		}
	}
}
