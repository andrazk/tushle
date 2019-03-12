package nkbm

import (
	"strconv"
	"strings"

	"tushle"
	"tushle/lib/browser"
)

const (
	acronym = "nkbm"
)

type localCred struct {
	Pass  string
	Store string
}

// nkbm struct
type nkbm struct {
	browser browser.Browser

	cli   tushle.Cli
	creds tushle.Credentials
}

// NewNkbm function
func NewNkbm(cli tushle.Cli, creds tushle.Credentials, b browser.Browser) tushle.Banker {
	return &nkbm{
		browser: b,
		cli:     cli,
		creds:   creds,
	}

	// if err := b.Run(login(pass, localStorage)); err != nil {
	// 	return nil, err
	// }
}

// Acronym retunrs bank acronym.
func (n *nkbm) Acronym() string {
	return acronym
}

func toFloat(s string) (float64, error) {
	cut := strings.Index(s, " ")
	if cut > -1 {
		s = s[0:cut]
	}
	s = strings.Replace(s, ".", "", -1)
	s = strings.Replace(s, ",", ".", -1)
	return strconv.ParseFloat(s, 64)
}
