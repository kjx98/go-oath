package main

import (
	"encoding/binary"
	"encoding/json"
	"github.com/kjx98/go-oath"
	"github.com/kjx98/go-oath/prompt"
	"hash/crc32"
	"io"
	"os"
)

type account struct {
	Name   string `json:"name"`
	Issuer string `json:"issuer"`
	Secret string `json:"secret"`
	otp    *oath.Oath
}

func parseAcct(tData []byte) []account {
	var accts []account
	if err := json.Unmarshal(tData, &accts); err != nil {
		logg.Error("Parse acct file", err)
		os.Exit(3)
	}
	for idx := range accts {
		accts[idx].otp = oath.New(accts[idx].Secret)
	}
	return accts
}

func readAcct(fin string, bForce bool) []account {
	var accts []account
	var acctDB string
	if homePtr := os.Getenv("HOME"); homePtr == "" {
		acctDB = ".authK/accts.db"
	} else {
		if bForce {
			os.Mkdir(homePtr+"/.authK", 0777)
		}
		acctDB = homePtr + "/.authK/accts.db"
	}
	if fin == "" {
		if fp, err := os.Open(acctDB); err != nil {
			logg.Error("account DB", err)
			os.Exit(3)
		} else if ss, err := io.ReadAll(fp); err != nil {
			logg.Error("Read account DB", err)
			os.Exit(3)
		} else if key, err := prompt.Stdin.PromptPassword("DB Passwd:"); err != nil {
			logg.Error("input password")
			os.Exit(3)
		} else if tData, err := New([]byte(key)).Decrypt(ss[4:]); err != nil {
			logg.Error("Decrypt", err)
			os.Exit(3)
		} else {
			crcSum := crc32.ChecksumIEEE(tData)
			if crcSum != binary.BigEndian.Uint32(ss) {
				logg.Error("Password mismatch")
				os.Exit(3)
			}
			return parseAcct(tData)
		}
	} else if fp, err := os.Open(fin); err != nil {
		logg.Error("Open acct file", err)
		os.Exit(3)
	} else if ss, err := io.ReadAll(fp); err != nil {
		logg.Error("Read acct file", err)
		os.Exit(3)
	} else {
		fp.Close()
		accts = parseAcct(ss)
		if bForce {
			if fp, err := os.Create(acctDB); err != nil {
				return accts
			} else {
				defer fp.Close()
				if key, err := prompt.Stdin.PromptPassword("DB Passwd:"); err != nil {
					logg.Error("input password")
					return accts
				} else if tData, err := New([]byte(key)).Encrypt(ss); err == nil {
					tCrc := make([]byte, 4)
					binary.BigEndian.PutUint32(tCrc, crc32.ChecksumIEEE(ss))
					fp.Write(tCrc)
					fp.Write(tData)
				}
			}
		}
	}
	return accts
}
