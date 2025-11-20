package bot_test

import (
	"bot-gateway/api/bot"
	"fmt"
	"log"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCutSlice(t *testing.T) {
	slice := []string{}

	for i := range 6 {
		slice = append(slice, fmt.Sprintf("%d", i+1))

		have := bot.CutSlice(slice, 3)
		log.Println(have)
		rowSlice := []string{}
		for j := range i + 1 {
			rowSlice = append(rowSlice, fmt.Sprintf("%d", j+1))
		}

		var want [][]string
		if len(rowSlice) < 4 {
			want = [][]string{rowSlice}
		} else {
			want = [][]string{rowSlice[:3], rowSlice[3:]}
		}

		assert.Equal(t, want, have)
	}
}
