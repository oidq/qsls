package sorter

import (
	"bitbucket.org/olik636/qsls/converter"
	"bitbucket.org/olik636/qsls/dxcc"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

var okCallRegexp = regexp.MustCompile(`^(OL|OK)`)
var suffixRegexp = regexp.MustCompile(`^[0-9A-z][A-z]*$`)

type Packet struct {
	Prefix string
	DXCC   uint
	QSLs   []*converter.QSLCard
}

type Packets []*Packet

type Sorter struct {
	priorPrefixRegExp []*regexp.Regexp
	priorPrefixSize   int
}

func NewSorter(regexps []*regexp.Regexp) *Sorter {
	c, i := 1, len(regexps)+1
	for i != 0 {
		i /= 10
		c++
	}
	return &Sorter{
		priorPrefixSize:   c,
		priorPrefixRegExp: regexps,
	}
}

func (s *Sorter) SortQSLsByDXCC(e *dxcc.Entities, qsos []*converter.QSLCard) Packets {
	var packets = map[uint]*Packet{}
	for _, v := range qsos {
		_, ok := packets[v.RefDXCC]
		if !ok {
			ent, _ := e.LookupEntityCode(int64(v.RefDXCC))
			packets[v.RefDXCC] = &Packet{Prefix: ent.Prefix, DXCC: v.RefDXCC, QSLs: []*converter.QSLCard{v}}
		} else {
			packets[v.RefDXCC].QSLs = append(packets[v.RefDXCC].QSLs, v)
		}
	}
	var qs []*Packet
	for _, v := range packets {
		sort.Slice(v.QSLs, func(i, j int) bool {
			return getRefCall(v.QSLs[i]) > getRefCall(v.QSLs[j])
		})
		qs = append(qs, v)
	}
	sort.Slice(qs, func(i, j int) bool {
		p1, p2 := qs[i].Prefix, qs[j].Prefix
		if okCallRegexp.MatchString(p1) {
			p1 = "0" + p1
		}
		if okCallRegexp.MatchString(p2) {
			p2 = "0" + p2
		}
		return p1 > p2
	})
	return qs
}

func (s *Sorter) SortQSLsByAlphabet(qsls []*converter.QSLCard) []*converter.QSLCard {
	sort.Slice(qsls, func(i, j int) bool {
		p1, p2 := removePrefix(getRefCall(qsls[i])), removePrefix(getRefCall(qsls[j]))
		if okCallRegexp.MatchString(p1) {
			p1 = "0" + p1
		}
		if okCallRegexp.MatchString(p2) {
			p2 = "0" + p2
		}
		return p1 > p2
	})
	return qsls
}

func (s *Sorter) getSortingRefCall(call string) string {
	for i, v := range s.priorPrefixRegExp {
		if v.MatchString(call) {
			p := strconv.Itoa(i)
			for len(p) < s.priorPrefixSize {
				p = "0" + p
			}
			return p + call
		}
	}
	return "Z" + call
}

func getRefCall(q *converter.QSLCard) string {
	if q.QSLVia != "" {
		return q.QSLVia
	}
	return q.Callsign
}

func (p Packets) Convert() []*converter.QSLCard {
	var qsls []*converter.QSLCard
	for _, v := range p {
		for _, q := range v.QSLs {
			qsls = append(qsls, q)
		}
	}
	return qsls
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
