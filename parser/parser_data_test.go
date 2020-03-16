package parser

import (
	"io"
	"strings"
)

func getTestADIF() io.Reader {
	return strings.NewReader(`
ADIF Export from RumLogNG by DL2RUM
For further info visit: http://www.dl2rum.de

<ADIF_VER:5>3.0.9
<CREATED_TIMESTAMP:15>20191228 082943
<PROGRAMID:8>RUMlogNG
<PROGRAMVERSION:5>4.4.1
<EOH>


<call:5>M6MAU  <qso_date:8>20191227 <time_on:6>071524 <band:3>40m  <freq:8>7.149000  <mode:3>SSB  <rst_sent:2>47  <qsl_rcvd:1>N  <eor>
<call:7>OR75USL  <qso_date:8>20191227 <time_on:6>103452 <band:3>40m  <freq:8>7.149000  <mode:3>SSB  <rst_sent:2>59  <qsl_rcvd:1>N<qsl_via:6>ON4GDA  <eor>
<call:5>DH5LA  <qso_date:8>20191228 <time_on:6>075841 <band:3>40m  <freq:8>7.197500  <mode:3>SSB  <rst_sent:2>59  <qsl_rcvd:1>N  <eor>
`)
}

func getTestADIFData() (Header, []QSO) {
	return Header{
			"ADIF_VER":          "3.0.9",
			"CREATED_TIMESTAMP": "20191228 082943",
			"PROGRAMID":         "RUMlogNG",
			"PROGRAMVERSION":    "4.4.1",
		},
		[]QSO{
			{
				"CALL":     "M6MAU",
				"FREQ":     "7.149000",
				"BAND":     "40m",
				"MODE":     "SSB",
				"QSO_DATE": "20191227",
				"QSL_RCVD": "N",
				"RST_SENT": "47",
				"TIME_ON":  "071524",
			},
			{
				"CALL":     "OR75USL",
				"FREQ":     "7.149000",
				"BAND":     "40m",
				"MODE":     "SSB",
				"QSO_DATE": "20191227",
				"QSL_RCVD": "N",
				"QSL_VIA":  "ON4GDA",
				"RST_SENT": "59",
				"TIME_ON":  "103452",
			},
			{
				"CALL":     "DH5LA",
				"FREQ":     "7.197500",
				"BAND":     "40m",
				"MODE":     "SSB",
				"QSO_DATE": "20191228",
				"QSL_RCVD": "N",
				"RST_SENT": "59",
				"TIME_ON":  "075841",
			},
		}
}
