package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const publicKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAnWetg+IAVsvzOk5Odfl3
ffKlyzRjHeFJOQ1jmqPGjU/Cl9i/rakDwNSK8dXJINLNyH8UH9AYbx0bmFbUx0eW
mPgeaZoYqf7noxLKs//U5bTgQjRGxmIPeoZ6CYvieEvO/eu5I9o6dT9WJHZxKGpY
Xna0RdozstZLiCUJE4gDmt8TRv4w4ydfs08fWSVdwD7SLyA6LacPrIPbU/vGYr+a
i0P4tdIN6UrfG7AnyUyGkOydIGIedRtMSmUso5ipoq1gTV5tFomdVF4rfjGQuIir
C1KN6wFQYWFNmktS0maAwwESnL1+2ATrJIm3fH/doP03r46ADKPjIeyJpgmsHiT6
qQIDAQAB
-----END PUBLIC KEY-----`

func TestVerify(t *testing.T) {
	pubKey, err := jwt.ParseRSAPublicKeyFromPEM([]byte(publicKey))
	if err != nil {
		t.Fatalf("cannot parse public key: %v", err)
	}
	v := &JWTTokenVerifier{
		PublicKey: pubKey,
	}
	cases := []struct {
		name    string
		tkn     string
		now     time.Time
		want    string
		wantErr bool
	}{
		{
			name: "valid_token",
			tkn:  "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjJmNzBlNWJhZDNlZTI3MDA0ZjY4MzU4In0.P7BjvXCyTtacFSaw79ViUmzXMpBMkNu_EhDPL8gevNw5xlSOv0O3rQWJCbGxfQonxB5cX87rQKdL0aaMf0mBCwSz2w8eeI_atbJCnFZb7Wp50TGDFPBIiMdo-IS-2TPa-uVPqYf6oCLl-Mr_4qDIm_-3atzyrmZ4WNcPPB34y8W8vmSZ2LxUP0SWpVvsPWQdaOoZnBLDlZvfD0L8IF3NyfBE2efdlXLBQATkoKrdyPN-a_7-_oh06SziqccIFgVydOC0-zEpCXIFpfDuBE8WTaRS2r4VdySQN0EOyWdPFsOJzBTCVDlN3jTZiTfXDl8lBF-VVy__BElZM16m-RlNTA",
			now:  time.Unix(1516239222, 0),
			want: "62f70e5bad3ee27004f68358",
		},
		{
			name:    "token_expired",
			tkn:     "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjJmNzBlNWJhZDNlZTI3MDA0ZjY4MzU4In0.P7BjvXCyTtacFSaw79ViUmzXMpBMkNu_EhDPL8gevNw5xlSOv0O3rQWJCbGxfQonxB5cX87rQKdL0aaMf0mBCwSz2w8eeI_atbJCnFZb7Wp50TGDFPBIiMdo-IS-2TPa-uVPqYf6oCLl-Mr_4qDIm_-3atzyrmZ4WNcPPB34y8W8vmSZ2LxUP0SWpVvsPWQdaOoZnBLDlZvfD0L8IF3NyfBE2efdlXLBQATkoKrdyPN-a_7-_oh06SziqccIFgVydOC0-zEpCXIFpfDuBE8WTaRS2r4VdySQN0EOyWdPFsOJzBTCVDlN3jTZiTfXDl8lBF-VVy__BElZM16m-RlNTA",
			now:     time.Unix(1516339222, 0),
			wantErr: true,
		},
		{
			name:    "bad_token",
			tkn:     "bad_token",
			now:     time.Unix(1516239222, 0),
			wantErr: true,
		},
		{
			name:    "wrong_signature",
			tkn:     "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NjEwODEwMjUsImlhdCI6MTY2MDQ3NjIyNSwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjJmNzBlNWJhZDNlZTI3MDA0ZjY4MzU3In0.P7BjvXCyTtacFSaw79ViUmzXMpBMkNu_EhDPL8gevNw5xlSOv0O3rQWJCbGxfQonxB5cX87rQKdL0aaMf0mBCwSz2w8eeI_atbJCnFZb7Wp50TGDFPBIiMdo-IS-2TPa-uVPqYf6oCLl-Mr_4qDIm_-3atzyrmZ4WNcPPB34y8W8vmSZ2LxUP0SWpVvsPWQdaOoZnBLDlZvfD0L8IF3NyfBE2efdlXLBQATkoKrdyPN-a_7-_oh06SziqccIFgVydOC0-zEpCXIFpfDuBE8WTaRS2r4VdySQN0EOyWdPFsOJzBTCVDlN3jTZiTfXDl8lBF-VVy__BElZM16m-RlNTA",
			now:     time.Unix(1516239222, 0),
			wantErr: true,
		},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			jwt.TimeFunc = func() time.Time {
				return c.now
			}
			accountID, err := v.Verify(c.tkn)
			if !c.wantErr && err != nil {
				t.Errorf("verification failed: %v", err)
			}
			if c.wantErr && err == nil {
				t.Error("want error got no error")
			}
			if accountID != c.want {
				t.Errorf("wrong accountID. want: %q; got: %q", c.want, accountID)
			}
		})
	}
}
