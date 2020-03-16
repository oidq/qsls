package parser

import (
	requirement "github.com/stretchr/testify/require"
	"strings"
	"testing"
)

type attrTest struct {
	attr  attribute
	isErr bool
}

var attrData = map[string]attrTest{
	"   <hugo:3>":         {attribute{"HUGO", 3}, false},
	" \n <HUGO_BOSS:314>": {attribute{"HUGO_BOSS", 314}, false},
	"<hugo:3a>":           {attribute{"", 0}, true},
}

func TestScanAttribute(t *testing.T) {
	require := requirement.New(t)
	for i, a := range attrData {
		rd := strings.NewReader(i)
		p := NewParser(rd)
		attr, err := p.scanAttribute()
		if a.isErr {
			require.Errorf(err, "scanAttributeName - %s", i)
		} else {
			require.NoErrorf(err, "scanAttributeName - %s", i)
			require.Equalf(a.attr.name, attr.name, "attrName - %s", i)
			require.Equalf(a.attr.size, attr.size, "attrSize - %s", i)
		}
	}
}

func TestParser(t *testing.T) {
	require := requirement.New(t)
	var header Header
	var qsos []QSO
	p := NewParser(getTestADIF())
	var out = make(chan Record, 10)
	err := p.Parse(out)
	for val := range out {
		if val.IsHeader() {
			header = val.(Header)
		} else {
			qsos = append(qsos, val.(QSO))
		}
	}
	require.NoError(err)
	expHeader, expQSOs := getTestADIFData()
	require.Equal(expHeader, header, "header")
	require.Equal(expQSOs, qsos, "qsos")
}
