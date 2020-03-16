package dxcc

import (
	"encoding/csv"
	"github.com/pkg/errors"
	"io"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

type Entities []Entity

type Entity struct {
	Entity       string
	Prefix       string
	DXCC         int
	Continent    string
	CQZone       int
	ITUZone      int
	Latitude     float64
	Longitude    float64
	Prefixes     []string
	Score        int
	PrefixRegexp *regexp.Regexp
}

func NewEntityDB(f io.Reader) (Entities, error) {
	var ents Entities
	cr := csv.NewReader(f)

	for {
		var ent Entity
		var err error
		record, err := cr.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			return nil, errors.Wrap(err, "read csv file")
		}

		ent.Prefix = record[0]
		ent.Entity = record[1]
		ent.Continent = record[3]
		ent.DXCC, err = strconv.Atoi(record[2])
		if err != nil {
			return nil, errors.Wrapf(err, "invalid dxcc - '%s'", record[2])
		}
		ent.CQZone, err = strconv.Atoi(record[4])
		if err != nil {
			return nil, errors.Wrapf(err, "invalid cqZone - '%s'", record[4])
		}
		ent.ITUZone, err = strconv.Atoi(record[5])
		if err != nil {
			return nil, errors.Wrapf(err, "invalid ituZone - '%s'", record[5])
		}
		ent.Latitude, err = strconv.ParseFloat(record[6], 64)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid latitude - '%s'", record[6])
		}
		ent.Longitude, err = strconv.ParseFloat(record[7], 64)
		if err != nil {
			return nil, errors.Wrapf(err, "invalid longitude - '%s'", record[7])
		}
		ent.Longitude *= -1

		ent.Prefixes = splitPrefixes(record[9])
		ent.PrefixRegexp, err = regexp.Compile(prefixRegexp(record[9]))
		if err != nil {
			return nil, errors.Wrapf(err, "invalid prefixRegexp - '%s'", record[9])
		}

		ents = append(ents, ent)
	}
	sort.Slice(ents, func(i, j int) bool {
		return ents[i].DXCC < ents[j].DXCC
	})
	return ents, nil
}

func splitPrefixes(pfx string) []string {
	var pfxs []string
	sb := strings.Builder{}

	pfx = strings.Replace(pfx, ";", "", -1)
	for _, p := range strings.Split(pfx, " ") {
		pfxs = append(pfxs, p)
		sb.WriteString(p)
	}
	return pfxs
}

func prefixRegexp(pfx string) string {

	initialChars := map[byte]struct{}{}
	pfx = strings.Replace(pfx, ";", "", -1)
	for _, p := range strings.Split(pfx, " ") {
		switch p[0] {
		case '=':
			initialChars[p[1]] = struct{}{}
		default:
			initialChars[p[0]] = struct{}{}
		}
	}
	sorted := []byte{}
	for c := range initialChars {
		sorted = append(sorted, c)
	}
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})
	sb := strings.Builder{}
	sb.WriteString("^[")
	for _, v := range sorted {
		sb.WriteByte(v)
	}
	sb.WriteByte(']')
	return sb.String()
}
