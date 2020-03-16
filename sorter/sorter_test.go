package sorter

import (
	requirement "github.com/stretchr/testify/require"
	"regexp"
	"testing"
)

func TestSortingRefCall(t *testing.T) {
	require := requirement.New(t)
	for _, d := range dataTestSortingRefCall {
		var regexps []*regexp.Regexp
		for _, r := range d.prefixRegexp {
			regexps = append(regexps, regexp.MustCompile(r))
		}
		s := NewSorter(regexps)
		prefixedCall := s.getSortingRefCall(d.call)
		require.Equal(d.prefixedCall, prefixedCall, "Wrong prefix")
	}
}

var dataTestSortingRefCall = []struct {
	call         string
	prefixRegexp []string
	prefixedCall string
}{
	{
		"OM3XXX",
		[]string{"^OM", "^OL"},
		"00OM3XXX",
	},
	{
		"0A1ABC",
		[]string{"^OL", "^.A"},
		"010A1ABC",
	},
	{
		"K1OG",
		[]string{"H", "H", "H", "H", "H", "H", "H", "H", "H", "H", "K"},
		"010K1OG",
	},
	{
		"K1OG",
		[]string{"H", "K", "H", "H", "H", "H", "H", "H", "H", "H"},
		"001K1OG",
	},
	{
		"OM3XXX",
		[]string{"H", "H"},
		"ZOM3XXX",
	},
}
