package jwt

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"time"
)

const TokenSignedSecurity = "cde842f51c676f1d3d9bcb3893dea36a0bc7c452366036bf57966cd9b0132e83"

type UserJwt struct {
	Token string
	Exp time.Time
}

func (u *UserJwt) JwtSignedString() (tokenString string, ok bool) {
	JWTToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"token": u.Token,
		//"exp": u.Exp,
	})
	fmt.Println("jwtsigned: ", u.Token, u.Exp)
	tokenString, err := JWTToken.SignedString([]byte(TokenSignedSecurity))
	if err != nil {
		fmt.Println("JwtSignedString: ", err)
		return "", false
	}
	return tokenString, true
}

func (u *UserJwt) JwtParse(tokenString string) {
	// Parse takes the token string and a function for looking up the key. The latter is especially
	// useful if you use multiple keys for your application.  The standard is to use 'kid' in the
	// head of the token to identify which key to use, but the parsed token (head and claims) is provided
	// to the callback, providing flexibility.
	token, _ := jwt.Parse(tokenString, func(token *jwt.Token) (i interface{}, e error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(TokenSignedSecurity), nil
	})

	if claims, ok := token.Claims.(jwt.MapClaims); ok {
		u.Token = claims["token"].(string)
		//layout := "2006-01-02T15:04:05.000Z"
		//str := claims["exp"].(string)
		//t, err := time.Parse(layout, str)
		//
		//if err != nil {
		//	fmt.Println(err)
		//}
		//
		//u.Exp = t
	}

}