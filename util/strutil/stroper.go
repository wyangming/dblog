package strutil

import (
	"bytes"
)

//追加一个字符串
func StrAppend(before, end string) string {
	var buf bytes.Buffer
	if len(before) > 0 {
		buf.WriteString(before)
	}
	if len(end) > 0 {
		buf.WriteString(end)
	}
	return buf.String()
}
