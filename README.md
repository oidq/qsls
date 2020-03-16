# QSLs

`qsls` is simple cli program for converting adif files, sorting QSOs and
exporting records for printing. Converting to PDF is done by LaTeX.

### Usage

Write a LaTeX template and following program will generate sorted QSLs from ADIF

`qsls export -a YourAdif.adif -o PDFOutput.pdf`

| Basic Parameters  |                                 |
| ----------------- | ------------------------------- |
| `-a`              | input adif file                 |
| `-o`              | output PDF file                 |
| `-s`              | advanced sorter (DXCC sorting)  |
| `--template`      | LaTeX template file             | 


### Templating
`qsls` uses Go build-in `text/template` package. Simple struct is carried to the template:
```json
{
  "QSLs": [
    {
      "Callsign": "OM3XXX",
      "QSLVia": "OK7XXX",
      "QSOs": {
        "Band": "40m",
        "Freq": "7.140",
        "Mode": "SSB",
        "RstSent": "59",
        "Time": "2020-03-19T07:22Z"
      }
    } 
  ],
  "User": {
    "var": "Variables specified in config"
  }
}
```
Simple example of LaTeX template can be found
[here](example/template.tex) and running 
`qsls export -a example/example.adif -o example/qsls-output.pdf --template example/template.tex`
will create [this PDF](example/qsls-output.pdf)

### Sorting
`qsls` can sort either by alphabet or by DXCC packets. To see how the final QSOs
will be sorted use `qsls show`

## Config
Default config file is located in `~/.config/qsls/cfg.json`
```json
{
  "data-dir": "/home/ham/.config/qsls/data",
  "template": "/home/ham/.config/qsls/qsls-template.tex",
  "user-var": {
    "qth": "Your QTH"
  },
  "prior-prefixes": ["^OM", "^OK"]
}
```
All files in config must be *absolute*. Running command with `--save-config` saves
supported configuration into config file.

##### template
`template` is LaTeX that will be templated with golang `text/template` package. 
[More about templating](README.md#Templating)

##### user-var
`user-var` is used in templating. Values can be accessed by `{{ index .User "qth" }}`

##### prior-prefixes
`prior-prefixes` is ordered list of regular expressions.
If callsign (or prefix in _advanced mode_) matches the expression it will be sorted
before other callsigns (or prefixes). 

`OK7XXX OM2OML OK8XXX OH2STA KW1TLK 9A3ESR IT2LPO ZA2LPA`

with `"prior-prefixes": ["^OK", "^9A"]` will be sorted as

`OK2OML OK8XXX 9A3ESR IT2LPO KW1TLK OH2STA OM2OML ZA2LPA`

## TODO
* validate QSLVia with QSL Manager Database
* csv output
