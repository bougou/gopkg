package sysfsutil

import (
	"fmt"
	"testing"
	"time"

	"github.com/google/uuid"
)

func Test_UUID(t *testing.T) {

	s := "31303232-534d-ec3c-ef2e-958900000000"
	u, err := uuid.Parse(s)
	if err != nil {
		t.Errorf("%s", err)
	}
	fmt.Println(u.String())

	/* GUID time-stamp is a 60-bit value representing the
	 * count of 100ns intervals since 00:00:00.00, 15 Oct 1582 */
	gregorianTimestamp := int64(u.Time()) * 100 / (10 * 9)
	fmt.Printf("%d\n", gregorianTimestamp)

	/* Seconds from 15 Oct 1582 to 1 Jan 1970 00:00:00 */
	var epochSinceGregorian int64 = 12219292800

	unixTimestamp := gregorianTimestamp - epochSinceGregorian

	a := time.Unix(unixTimestamp, 0)
	fmt.Println(a)

}
