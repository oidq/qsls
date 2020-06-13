package dxcc

import (
	"errors"
	"sort"
	"strconv"
	"strings"
)

/*
	This code was shamelessly copied from github.com/tzneal/ham-go/dxcc
	Modified to load data from csv file not from pre-generated go-file
	Thanks Todd!
 */

func (e Entities) Lookup(callsign string) (Entity, bool) {
	callsign = strings.TrimSpace(strings.ToUpper(callsign))
	matchedEntities := []Entity{}
	for _, ent := range e {
		if ent.PrefixRegexp.MatchString(callsign) {
			if matched, ok := ent.Match(callsign); ok {
				matchedEntities = append(matchedEntities, matched)
			}
		}
	}
	sort.Slice(matchedEntities, func(i, j int) bool {
		return matchedEntities[i].Score > matchedEntities[j].Score
	})
	if len(matchedEntities) > 0 {
		return matchedEntities[0], true
	}
	return Entity{}, false
}

var markers = []byte{'(', '[', '<', '{', '~'}

func (e Entity) Match(callsign string) (Entity, bool) {
	callsign = strings.TrimSpace(strings.ToUpper(callsign))

	for _, pfx := range e.Prefixes {
		if pfx[0] == '=' {
			exactCall := pfx[1:]
			for _, oc := range markers {
				if idx := strings.IndexByte(exactCall, oc); idx != -1 {
					exactCall = exactCall[0:idx]
				}
			}

			if exactCall == callsign {
				ent := e
				applyOverrides(pfx[1:], &ent)
				ent.Score = len(callsign)
				return ent, true
			}
		} else {
			origPfx := pfx
			for _, oc := range markers {
				if idx := strings.IndexByte(pfx, oc); idx != -1 {
					pfx = pfx[0:idx]
				}
			}
			if strings.HasPrefix(callsign, pfx) {
				if len(pfx) > e.Score {
					e.Score = len(pfx)
					applyOverrides(origPfx, &e)
				}
			}
		}
	}
	if e.Score > 0 {
		return e, true
	}
	return e, false
}

func applyOverrides(pfx string, ent *Entity) {
	i := 0
	for i < len(pfx) {
		for _, oc := range markers {
			if pfx[i] == oc {
				ec := byte(')')
				switch oc {
				case '(':
					ec = ')'
				case '[':
					ec = ']'
				case '<':
					ec = '>'
				case '{':
					ec = '}'
				case '~':
					ec = '~'
				}
				i++

				j := i
				for pfx[j] != ec {
					j++
				}

				switch oc {
				case '(':
					value, err := strconv.ParseInt(pfx[i:j], 10, 64)
					if err == nil {
						ent.CQZone = int(value)
					}
				case '[':
					value, err := strconv.ParseFloat(pfx[i:j], 64)
					if err == nil {
						ent.ITUZone = int(value)
					}
				case '<':
				case '{':
				case '~':
				}

			}
		}
		i++
	}
}

func (e Entities) LookupEntity(name string) (Entity, error) {
	for _, v := range e {
		if v.Entity == name {
			return v, nil
		}
	}
	return Entity{}, errors.New("entity not found")
}

func (e Entities) LookupEntityCode(code int64) (Entity, error) {
	for _, v := range e {
		if int64(v.DXCC) == code && !strings.Contains(v.Prefix, "*") {
			return v, nil
		}
	}
	return Entity{}, errors.New("entity not found")
}
