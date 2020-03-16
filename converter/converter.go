package converter

import (
	"bitbucket.org/olik636/qsls/parser"
	"github.com/pkg/errors"
	"strconv"
	"time"
)

type QSLCard struct {
	Callsign string
	QSLVia   string
	RefDXCC  uint
	QSOs     []QSOInfo
}

type QSO struct {
	QSOInfo
	Callsign string
	QSLVia   string
	RefDXCC  uint
}

type QSOInfo struct {
	Band    string
	Freq    string
	Mode    string
	RstSent string
	QSLRcvd bool
	Time    time.Time
}

func ConvertMapToQSO(oQso parser.QSO) (*QSO, error) {
	var qso = &QSO{}
	var year, month, day, hour, minute int
	for i, v := range oQso {
		switch i {
		case "CALL":
			qso.Callsign = v
		case "BAND":
			qso.Band = v
		case "FREQ":
			qso.Freq = v[:5]
		case "MODE":
			qso.Mode = v
		case "QSL_VIA":
			qso.QSLVia = v
		case "QSL_RCVD":
			qso.QSLRcvd = v == "Y"
		case "RST_SENT":
			qso.RstSent = v
		case "QSO_DATE":
			i, err := strconv.Atoi(v)
			if err != nil {
				return nil, errors.Wrapf(err, "invalid date '%s'", v)
			}
			year, month, day = i/10000, (i/100)%100, i%100
		case "TIME_ON":
			i, err := strconv.Atoi(v)
			if err != nil {
				return nil, errors.Wrapf(err, "invalid time '%s'", v)
			}
			hour, minute = i/10000, (i/100)%100
		}
	}
	qso.Time = time.Date(year, time.Month(month), day, hour, minute, 0, 0, time.UTC)
	return qso, nil
}
