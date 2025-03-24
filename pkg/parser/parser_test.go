package parser

import (
	"strings"
	"testing"
	"time"
)

func TestParser(t *testing.T) {
	str := "fdfd ffgfg erfgg cggf"
	for field := range strings.FieldsSeq(str) {
		t.Log(field)
	}
	d := strings.Fields(str)
	t.Log(d[0])
	t.Log(time.Now().Year()-2000)
}
