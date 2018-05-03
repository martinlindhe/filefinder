package finder

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrettyDataSize(t *testing.T) {
	data := map[int64]string{
		0:             "0 bytes",
		100:           "100 bytes",
		1024:          "1.0 KiB",
		16384:         "16.0 KiB",
		1048576:       "1.0 MiB",
		4194304:       "4.0 MiB",
		1073741824:    "1.0 GiB",
		4294967296:    "4.0 GiB",
		1099511627776: "1.0 TiB",
		4398046511104: "4.0 TiB",
	}
	for n, expect := range data {
		assert.Equal(t, expect, prettyDataSize(n))
	}
}

func TestParseDataSize(t *testing.T) {
	data := map[string]int64{
		"":        0,
		"100":     100,
		"1K":      1024,
		"1 k":     1024,
		"1 MB":    1048576,
		"2.8 MiB": 2936012,
		"4GiB":    4294967296,
		"4TB":     4398046511104,
	}
	for s, expect := range data {
		assert.Equal(t, expect, parseDataSize(s))
	}
}
