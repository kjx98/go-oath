package oath

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base32"
	"encoding/binary"
	"time"
)

const (
	Interval int64 = 30
)

type Oath struct {
	secret []byte
}

func New(key string) *Oath {
	if secr, err := base32.StdEncoding.DecodeString(key); err == nil {
		secret := make([]byte, len(secr))
		copy(secret, secr)
		return &Oath{secret}
	} else {
		panic("Invalid secret")
	}
}

func calcOTP(msg []byte, key []byte) int {
	hh := hmac.New(sha1.New, key)
	hh.Write(msg)
	h := hh.Sum(nil)
	if len(h) != 20 {
		return 0
	}
	off := uint(h[19] & 0xf)
	u32 := binary.BigEndian.Uint32(h[off:]) & 0x7fffffff
	return int(u32 % 1000000)
}

// return RFC 4226 HOTP
func (totp *Oath) Rfc(numb uint64) int {
	msg := make([]byte, 8)
	binary.BigEndian.PutUint64(msg, numb)
	return calcOTP(msg, totp.secret)
}

// return TOTP
func (totp *Oath) Now() int {
	tt := time.Now().UTC().Unix() / Interval
	return totp.Rfc(uint64(tt))
}

// return TOTP
func TOTP(key []byte) int {
	tt := time.Now().UTC().Unix() / Interval
	msg := make([]byte, 8)
	binary.BigEndian.PutUint64(msg, uint64(tt))
	return calcOTP(msg, key)
}
