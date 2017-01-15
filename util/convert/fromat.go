package convert

import (
	"fmt"
	"time"
)

const (
	timeFormat = "2006-01-02 15:04:05.999999"
)

//2006-01-02 15:04:05.999999
func Str2time(str string) (t time.Time, err error) {
	base := "0000-00-00 00:00:00.0000000"
	switch len(str) {
	case 10, 19, 21, 22, 23, 24, 25, 26: // up to "YYYY-MM-DD HH:MM:SS.MMMMMM"
		if str == base[:len(str)] {
			return
		}
		t, err = time.Parse(timeFormat[:len(str)], str)
	default:
		err = fmt.Errorf("invalid time string: %s", str)
		return
	}
	return
}
