package main

import (
	"encoding/json"
	"github.com/kjx98/go-oath"
	"io"
	"os"
)

type account struct {
	Name   string `json:"name"`
	Issuer string `json:"issuer"`
	Secret string `json:"secret"`
	otp    *oath.Oath
}

func readAcct(fin string) []account {
	var accts []account
	var ff *os.File
	if fin == "" {
		ff = os.Stdin
	} else if fp, err := os.Open(fin); err == nil {
		ff = fp
	} else {
		logg.Error("Open acct file", err)
		panic(err)
	}
	if ss, err := io.ReadAll(ff); err != nil {
		logg.Error("Read acct file", err)
		panic(err)
	} else if err = json.Unmarshal(ss, &accts); err != nil {
		logg.Error("Parse acct file", err)
		panic(err)
	}
	for idx := range accts {
		accts[idx].otp = oath.New(accts[idx].Secret)
	}
	return accts
}
