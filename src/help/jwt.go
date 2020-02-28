package help

import (
	"crypto/rsa"
	"fmt"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
)

//rsa公钥私钥
//http://web.chacuo.net/netrsakeypair
var (
	privatekeyPEM = []byte(`-----BEGIN RSA PRIVATE KEY-----
MIICXAIBAAKBgQCIzQ88thakRQX48bEeqNTjXiIoXN56EWTwJwICcTmWhLHuo2PW
guc8mxlz9PnXSR7ZoD6IC70lQBI6Yyb2HCCBu9XuwPPD+3CCXaAvNh+WaOf4/8GW
FawEo3IZpU+Eg1wTg1Hv2f4vOfqH1PeyN5xG6YzWwCeBbxNK/RsCZdK3gQIDAQAB
AoGABAfeuhA3bJGoEmS0rk2VMDnW0c+isoAOtFDB62aApuG0CG8CAxEKVSvQvSzr
q5847HqfPQzzfsR7hQLk4/2CK6TsPrfBjW9Uc5YiuweU74rhJtH9/GfZorXDmxaf
u9v7UMJohSGYVoHtL3fw1j7sWaUBhl+CA1Qb2HdxBunWMxECQQCMf/x1A2NRkG/I
gOVS8FFbf98rBW04fgGzW2kUrFM1RUd/kQslVc6pnLeu7LhP9hC7rtjO7Sm5A0E7
jVFB/4MNAkEA+UKjakM1/YteeC3ZjleOQfJIWGgOGxYS7SFQrwU7rvoyACMPHTI9
h4ZUCvUTp+9RLVB8R+FGU+9ZQyUtP0C5RQJAdaYeYmVp3zzRPdYhMWgm2DWlTEMJ
CEsLZYLf5P2/11Wh30I3URYfLYwbi5CRbfOgY2iwB+Y0D8aX8yQMrPUmaQJBAKGv
Ly9Ln6byk3njS97APp/aWEE4ZgX94JL+7EZLX7aVxn8+PpySrUTOxo9A/9oMK5z4
O1Wo9CSX+k/Kurnv8v0CQAnVIWp+FmXdl8b2zH5xXzR3tQ1aRnoOzi8la/GD+r+n
FQrMvaIPB8d5KeUjXuAfpuJzY46AbXSE3bMxE70zHrs=
-----END RSA PRIVATE KEY-----`)

	publickeyPEM = []byte(`-----BEGIN PUBLIC KEY-----
MIGfMA0GCSqGSIb3DQEBAQUAA4GNADCBiQKBgQCIzQ88thakRQX48bEeqNTjXiIo
XN56EWTwJwICcTmWhLHuo2PWguc8mxlz9PnXSR7ZoD6IC70lQBI6Yyb2HCCBu9Xu
wPPD+3CCXaAvNh+WaOf4/8GWFawEo3IZpU+Eg1wTg1Hv2f4vOfqH1PeyN5xG6YzW
wCeBbxNK/RsCZdK3gQIDAQAB
-----END PUBLIC KEY-----`)
)

var privateKey *rsa.PrivateKey
var publicKey *rsa.PublicKey

func init() {
	var err error
	privateKey, err = jwt.ParseRSAPrivateKeyFromPEM(privatekeyPEM)
	if err != nil {
		fmt.Println(err)
		panic("privateKey parse error")
	}
	publicKey, err = jwt.ParseRSAPublicKeyFromPEM(publickeyPEM)
	if err != nil {
		panic("publicKey parse error")
	}
}

type JwtPayload struct {
	AppID     string `json:"app_id,omitempty"`
	Email     string `json:"email,omitempty"`
	ExpiresAt int64  `json:"expires_at,omitempty"`
	Type      string `json:"type,omitempty"`
}

type JwtToken struct {
	AccessToken  string
	RefreshToken string
}

var TimeFunc = time.Now

func GenerateJwtToken(payload *JwtPayload, generateRefresh bool) (*JwtToken, error) {
	jwt.TimeFunc = TimeFunc

	createTimeUnix := time.Now().Unix()

	claim := jwt.MapClaims{
		"email": payload.Email,
		"exp":   payload.ExpiresAt,
		"appid": payload.AppID,
		"iat":   createTimeUnix,
	}

	refreshClaim := jwt.MapClaims{
		"email": payload.Email,
		"appid": payload.AppID,
		"iat":   createTimeUnix,
	}

	if payload.Type != "" {
		claim["type"] = payload.Type
		refreshClaim["type"] = payload.Type
	}

	token := jwt.NewWithClaims(jwt.SigningMethodRS256, claim)
	sToken, err := token.SignedString(privateKey)
	if err != nil {
		return nil, err
	}

	if !generateRefresh {
		return &JwtToken{
			AccessToken: sToken,
		}, nil
	}

	token = jwt.NewWithClaims(jwt.SigningMethodRS256, refreshClaim)
	sRefToken, err := token.SignedString(privateKey)
	if err != nil {
		return nil, err
	}

	return &JwtToken{
		AccessToken:  sToken,
		RefreshToken: sRefToken,
	}, nil
}

func ValidateJwtToken(token string) (*JwtPayload, error) {
	jToken, err := jwt.ParseWithClaims(token, jwt.MapClaims{}, func(tok *jwt.Token) (interface{}, error) {
		return publicKey, nil
	})

	if err != nil {
		return nil, err
	}

	claims, ok := jToken.Claims.(jwt.MapClaims)
	if !ok {
		return nil, err
	}

	payload := &JwtPayload{}

	if email, ok := claims["email"].(string); ok {
		payload.Email = email
	}

	if exp, ok := claims["exp"].(float64); ok {
		payload.ExpiresAt = int64(exp)
	}

	if appid, ok := claims["appid"].(string); ok {
		payload.AppID = appid
	}

	if tType, ok := claims["type"].(string); ok {
		payload.Type = tType
	}
	//这里校验超时
	err = claims.Valid()
	if err != nil {
		return payload, err
	}

	return payload, nil
}
