package token

import (
	"testing"
	"time"

	"github.com/dgrijalva/jwt-go"
)

const privateKey = `-----BEGIN PRIVATE KEY-----
MIIEvQIBADANBgkqhkiG9w0BAQEFAASCBKcwggSjAgEAAoIBAQCdZ62D4gBWy/M6
Tk51+Xd98qXLNGMd4Uk5DWOao8aNT8KX2L+tqQPA1Irx1ckg0s3IfxQf0BhvHRuY
VtTHR5aY+B5pmhip/uejEsqz/9TltOBCNEbGYg96hnoJi+J4S87967kj2jp1P1Yk
dnEoalhedrRF2jOy1kuIJQkTiAOa3xNG/jDjJ1+zTx9ZJV3APtIvIDotpw+sg9tT
+8Ziv5qLQ/i10g3pSt8bsCfJTIaQ7J0gYh51G0xKZSyjmKmirWBNXm0WiZ1UXit+
MZC4iKsLUo3rAVBhYU2aS1LSZoDDARKcvX7YBOskibd8f92g/TevjoAMo+Mh7Imm
CaweJPqpAgMBAAECggEBAIw0tZIr1TF7KYReC/V56L3/TT7bww3yhk6TZo1wJIPq
7+Jh5xrA2d8Bc2JGk4jxPOvChiJwMdOHkfT4Iz/+vF41ZKGb6SxDKgFP087RqsmR
e9B80C4VWsRA1KN8PpX4sL/tIFSXJksZx5ljBxiA4YYDJkCyRCqgR1dV5efH164y
cCXACRxf1UyNpJt9zZoOPg19IQTUGwteKSop9oU/WxVhjddz3AlBb+okDCfsrjn+
kzxu4fa2nEF1tpN44X8goHeMkXryuIdsehBkmMl6jaH1t64mwp/P+cGT+10AmwL9
k/EjggfbnuUgEZiLri/6MQ+vPasH/SaSEo+LajQz3AECgYEAzz+zEObWat05Q5NA
p8r1XxKVjUVfcLRWQsoDBhvy9JlCM1ha/WsqmrCizs1GKtcS9GHP65zpGR/ZfCP2
vYJy4qSdIvk//TFF6L5WUaN6rhpJzoCBZfYtCFGhtQ57bYvaQazPdq53oESj0PZN
0KTmeLbvSNOONSlJ7fIW8UueB9kCgYEAwm5vuK6jdsyZn/Ohha9HlvwlbsRVF+cN
Z6rCr3CynYEx7PgjYFqiGp5lacWSjscFNY3U5fQMpewN6sQd4jKITnjPSHe2gfTV
Yt0xDdWTM6atxUIFVSWxi3pZkgsVhntrnlF2dJa9gHpfZfFQyfXXqY/34cBE+HVk
t405utvwF1ECgYBRxJ8gvwK//PJ37+Qlj5UJ4qowp7tFG1GhXlSdF2/fA4yz91tG
+v4/NAu4LhNOGbc3xlOjcTAioodLTGEwWgR72VjKEK8ndUZQ0q/529cuU97k45yq
HtubmaGEbudRzEjbepQMDj/ScuJzMop3FGh+HicAg79qyBSMFeTpZN0/2QKBgHDc
Ximb5fMdzMcWStoo5qtz7d6gRKy9SAC3FI92IZhf2DUvzIkv0w0UiNWfA/Ww/Qsb
K0vYIEdoAKQX9yjIIGs8oUX1h5FkJ0FeGA1pviqrRA9OxX2phafq+3dUy8fmeI/L
xbDjl1iusBWiwDybYfZhRYhbbS20JySM68fVx0YhAoGAAT79QqgEKJRz7N9YPqig
HewOEOB7BXzHjF+6nkNQFgAe/tA1dYYtMqGaGDnpOjnJYnW4NWtwsjIcSXgOw37W
Y0KOpJ/Uk9xaf4e3mEJuRc2S5tkeQ+sdywpRqxC87xL6VzaJmtZJ0UiOutBum6XX
kjoFFa/UJpPAVwOWISPAOdI=
-----END PRIVATE KEY-----`

func TestGenerateToken(t *testing.T) {
	key, err := jwt.ParseRSAPrivateKeyFromPEM([]byte(privateKey))
	if err != nil {
		t.Fatalf("cannot parse private key: %v", err)
	}
	g := NewJWTTokenGen("coolcar/auth", key)
	g.nowFunc = func() time.Time {
		return time.Unix(1516239022, 0)
	}
	tkn, err := g.GenerateToken("62f70e5bad3ee27004f68358", 2*time.Hour)
	if err != nil {
		t.Errorf("cannot generate token: %v", err)
	}
	want := "eyJhbGciOiJSUzUxMiIsInR5cCI6IkpXVCJ9.eyJleHAiOjE1MTYyNDYyMjIsImlhdCI6MTUxNjIzOTAyMiwiaXNzIjoiY29vbGNhci9hdXRoIiwic3ViIjoiNjJmNzBlNWJhZDNlZTI3MDA0ZjY4MzU4In0.P7BjvXCyTtacFSaw79ViUmzXMpBMkNu_EhDPL8gevNw5xlSOv0O3rQWJCbGxfQonxB5cX87rQKdL0aaMf0mBCwSz2w8eeI_atbJCnFZb7Wp50TGDFPBIiMdo-IS-2TPa-uVPqYf6oCLl-Mr_4qDIm_-3atzyrmZ4WNcPPB34y8W8vmSZ2LxUP0SWpVvsPWQdaOoZnBLDlZvfD0L8IF3NyfBE2efdlXLBQATkoKrdyPN-a_7-_oh06SziqccIFgVydOC0-zEpCXIFpfDuBE8WTaRS2r4VdySQN0EOyWdPFsOJzBTCVDlN3jTZiTfXDl8lBF-VVy__BElZM16m-RlNTA"
	if tkn != want {
		t.Errorf("wrong token generated. want: %q; got: %q", want, tkn)
	}
}
