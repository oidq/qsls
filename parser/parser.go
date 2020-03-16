package parser

import (
	"bufio"
	"fmt"
	"io"
	"regexp"
	"strconv"
	"strings"
	"unicode/utf8"
)

var allowedSpace = regexp.MustCompile(`\s`)
var allowedAttribute = regexp.MustCompile(`^[a-zA-Z0-9_]+$`)

type cursor struct {
	lineNumber uint
	charOnLine uint
}

type attribute struct {
	name string
	size uint
}

type Header map[string]string

type QSO map[string]string

func (h Header) IsHeader() bool {
	return true
}

func (q QSO) IsHeader() bool {
	return false
}

type ADIFParser struct {
	r        *bufio.Reader
	position cursor
}

type Record interface {
	IsHeader() bool
}

func NewParser(r io.Reader) *ADIFParser {
	return &ADIFParser{r: bufio.NewReader(r)}
}

func (p *ADIFParser) Parse(c chan<- Record) error {
	defer close(c)
	err := p.skipComment()
	if err != nil {
		return p.formatError(err)
	}
	header, err := p.parseHeader()
	if err != nil {
		return p.formatError(err)
	}
	c <- header
	for {
		qso, err := p.parseQSO()
		if err == io.EOF {
			break
		} else if err != nil {
			return p.formatError(err)
		}
		c <- qso
	}
	return nil
}

func (p *ADIFParser) parseQSO() (QSO, error) {
	var qso = QSO{}
	for {
		attr, err := p.scanAttribute()
		if err != nil {
			return nil, err
		}
		if attr.name == "EOR" {
			return qso, nil
		}
		val, err := p.scanValue(attr.size)
		if err != nil {
			return nil, err
		}
		qso[attr.name] = val
	}
}

func (p *ADIFParser) parseHeader() (Header, error) {
	var header = Header{}
	for {
		attr, err := p.scanAttribute()
		if err != nil {
			return nil, err
		}
		if attr.name == "EOH" {
			return header, nil
		}
		val, err := p.scanValue(attr.size)
		if err != nil {
			return nil, err
		}
		header[attr.name] = val
	}
}

func (p *ADIFParser) readRune() (rune, error) {
	r, _, err := p.r.ReadRune()
	if err != nil {
		return ' ', err
	}
	switch r {
	case '\n':
		p.position.lineNumber++
		p.position.charOnLine = 0
	case '\r':
	default:
		p.position.charOnLine++
	}
	return r, nil
}

func (p *ADIFParser) unreadRune() error {
	err := p.r.UnreadRune()
	if err != nil {
		return err
	}
	p.position.charOnLine--
	if p.position.charOnLine < 0 {
		p.position.charOnLine = 0
		p.position.lineNumber--
	}
	return nil
}

func (p *ADIFParser) skipComment() error {
	for {
		r, err := p.readRune()
		if err != nil {
			return err
		}
		if r == '<' {
			return p.unreadRune()
		}
	}
}

func (p *ADIFParser) scanAttribute() (*attribute, error) {
	var attr = &attribute{}
	var hasSize bool
	err := p.findAttribute()
	if err != nil {
		return nil, err
	}
	attr.name, hasSize, err = p.scanAttributeName()
	if err != nil {
		return nil, err
	}
	if !hasSize {
		attr.size = 0
		return attr, nil
	}
	attr.size, err = p.scanAttributeSize()
	if err != nil {
		return nil, err
	}
	return attr, nil
}

func (p *ADIFParser) scanValue(size uint) (string, error) {
	var value string
	for i := 0; i < int(size); i++ {
		r, err := p.readRune()
		if err != nil {
			return "", err
		}
		value += string(r)
	}
	if !utf8.ValidString(value) {
		return "", fmt.Errorf("invalid value - '%s'", value)
	}
	return value, nil
}

func (p *ADIFParser) scanAttributeName() (string, bool, error) {
	var name string
	for {
		r, err := p.readRune()
		if err != nil {
			return "", false, err
		}
		if r == ':' || r == '>' {
			if allowedAttribute.MatchString(name) {
				return strings.ToUpper(name), r == ':', nil
			} else {
				return "", r == ':', fmt.Errorf("invalid attribute format '%s'", name)
			}
		}
		name += string(r)
	}
}

func (p *ADIFParser) scanAttributeSize() (uint, error) {
	var size string
	for {
		r, err := p.readRune()
		if err != nil {
			return 0, err
		}
		if r == '>' {
			n, err := strconv.Atoi(size)
			if err != nil {
				return 0, fmt.Errorf("invalid size character")
			}
			return uint(n), nil
		}
		size += string(r)
	}
}

func (p *ADIFParser) findAttribute() error {
	for {
		r, err := p.readRune()
		if err != nil {
			return err
		}
		if !allowedSpace.MatchString(string(r)) {
			if r == '<' {
				return nil
			} else {
				return fmt.Errorf("unknown character - '%s'", string(r))
			}
		}
	}
}

func (p *ADIFParser) formatError(err error) error {
	return fmt.Errorf("%s - %d:%d", err.Error(), p.position.lineNumber, p.position.charOnLine)
}
