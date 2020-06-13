package converter

import (
	"github.com/oIdq/qsls/dxcc"
	"github.com/oIdq/qsls/parser"
	"github.com/pkg/errors"
	"io"
	"log"
	"regexp"
	"strings"
)

var suffixRegexp = regexp.MustCompile(`^[0-9A-z][A-z]*$`)

func ParserWrapper(f io.Reader, e *dxcc.Entities) ([]*QSLCard, error) {
	var qs []*QSLCard
	var calls = map[string]int{}
	var c = make(chan parser.Record, 10)
	var done = make(chan bool, 1)
	p := parser.NewParser(f)
	go func() {
		for v := range c {
			if !v.IsHeader() {
				pq := v.(parser.QSO)
				q, err := ConvertMapToQSO(pq)
				if err != nil {
					log.Printf("Cannot convert qso: '%s'", err.Error())
					continue
				}
				i, ok := calls[q.Callsign]
				if !ok {
					ent, found := e.Lookup(removePrefix(getQSORefCall(q)))
					if !found {
						log.Printf("DXCC entity for %v not found", pq)
					}
					qsl := &QSLCard{Callsign: q.Callsign, QSLVia: q.QSLVia, RefDXCC: uint(ent.DXCC), QSOs: []QSOInfo{q.QSOInfo}}
					calls[q.Callsign] = len(qs)
					qs = append(qs, qsl)
				} else {
					qs[i].QSOs = append(qs[i].QSOs, q.QSOInfo)
				}
			}
		}
		close(done)
	}()
	err := p.Parse(c)
	if err != nil {
		return nil, errors.Wrap(err, "parse")
	}
	<-done
	return qs, nil
}

func removePrefix(callsign string) string {
	c := strings.Split(callsign, "/")
	if len(c) == 1 {
		return callsign
	} else if len(c) == 2 {
		if suffixRegexp.MatchString(c[1]) {
			return c[0]
		} else {
			return c[1]
		}
	} else if len(c) == 3 {
		return c[1]
	}
	return callsign
}

func getQSORefCall(q *QSO) string {
	if q.QSLVia != "" {
		return q.QSLVia
	}
	return q.Callsign
}
