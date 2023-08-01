package util

import (
	"strings"

	"github.com/shopspring/decimal"
)

func DecimalJoin(a []decimal.Decimal, sep string) string {
	switch len(a) {
	case 0:
		return ""
	case 1:
		return a[0].String()
	}
	n := len(sep) * (len(a) - 1)
	for i := 0; i < len(a); i++ {
		n += len(a[i].String())
	}

	var b strings.Builder
	b.Grow(n)
	b.WriteString(a[0].String())
	for _, s := range a[1:] {
		b.WriteString(sep)
		b.WriteString(s.String())
	}
	return b.String()
}
