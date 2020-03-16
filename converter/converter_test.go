package converter

import (
	requirement "github.com/stretchr/testify/require"
	"testing"
	"time"
)

func TestConvertor(t *testing.T) {
	require := requirement.New(t)
	mQsos, qsos := getConvertorTestData()
	for i, v := range mQsos {
		qso, err := ConvertMapToQSO(v)
		require.NoErrorf(err, "ConvertMapToQSO - %d", i)
		require.Equal(qsos[i].Time.Format(time.RFC822), qso.Time.Format(time.RFC822), "time")
		qsos[i].Time = qso.Time //time needs special assertion
		require.Equal(qsos[i], qso, "qso")
	}

}

func getConvertorTestData() ([]map[string]string, []*QSO) {
	return []map[string]string{
			{
				"CALL":     "M6MQZ",
				"FREQ":     "7.149000",
				"BAND":     "40m",
				"MODE":     "SSB",
				"QSO_DATE": "20191227",
				"QSL_RCVD": "N",
				"RST_SENT": "47",
				"TIME_ON":  "071524",
			},
			{
				"CALL":     "OR75USA",
				"FREQ":     "7.149000",
				"BAND":     "40m",
				"MODE":     "SSB",
				"QSO_DATE": "20191227",
				"QSL_RCVD": "N",
				"QSL_VIA":  "ON4GDV",
				"RST_SENT": "59",
				"TIME_ON":  "103452",
			},
			{
				"CALL":     "DH5JN",
				"FREQ":     "7.197500",
				"BAND":     "40m",
				"MODE":     "SSB",
				"QSO_DATE": "20191228",
				"QSL_RCVD": "N",
				"RST_SENT": "59",
				"TIME_ON":  "075841",
			},
		}, []*QSO{
			{
				Callsign: "M6MQZ",
				QSOInfo: QSOInfo{
					Freq:    "7.149",
					Band:    "40m",
					Mode:    "SSB",
					QSLRcvd: false,
					RstSent: "47",
					Time:    time.Date(2019, 12, 27, 7, 15, 0, 0, time.UTC),
				},
			},
			{
				Callsign: "OR75USA",
				QSLVia:   "ON4GDV",
				QSOInfo: QSOInfo{
					Freq:    "7.149",
					Band:    "40m",
					Mode:    "SSB",
					QSLRcvd: false,
					RstSent: "59",
					Time:    time.Date(2019, 12, 27, 10, 34, 0, 0, time.UTC),
				},
			},
			{
				Callsign: "DH5JN",
				QSOInfo: QSOInfo{
					Freq:    "7.197",
					Band:    "40m",
					Mode:    "SSB",
					QSLRcvd: false,
					RstSent: "59",
					Time:    time.Date(2019, 12, 28, 7, 58, 0, 0, time.UTC),
				},
			},
		}
}
