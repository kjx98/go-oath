package oath

import (
	"testing"
)

func TestHOTP(t *testing.T) {
	hotp := New("3ORMAAI2NRT7OGA6")
	expect := 627223
	if otp := hotp.Rfc(123); otp != expect {
		t.Errorf("OATH-HOTP: %06d expect %06d", otp, expect)
	}
	expect = 104563
	if otp := hotp.Rfc(7890123456); otp != expect {
		t.Errorf("OATH-HOTP: %06d expect %06d", otp, expect)
	}
	expect = 989260
	if otp := hotp.Rfc(123456789); otp != expect {
		t.Errorf("OATH-HOTP: %06d expect %06d", otp, expect)
	}
}
