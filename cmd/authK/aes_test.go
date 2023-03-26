package main

import (
	"bytes"
	"testing"
)

func TestAes(t *testing.T) {
	ciph := New([]byte("123456789012345678"))
	t1 := []byte("test is ok")
	if encrypted, err := ciph.Encrypt(t1); err != nil {
		t.Error("encrypt error", err)
	} else if t2, err := ciph.Decrypt(encrypted); err != nil {
		t.Error("decrypt error", err)
	} else if !bytes.Equal(t1, t2) {
		t.Error("decrypted diff orin")
	}
}
