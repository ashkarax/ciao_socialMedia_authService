package jwttoken_authSvc

import (
	"errors"
	"fmt"
	"time"

	interface_jwt_authSvc "github.com/ashkarax/ciao_socialMedia_authService/utils/jwt.go/interface"
	"github.com/golang-jwt/jwt"
)

type JwtUtil struct{}

func NewJwtUtil() interface_jwt_authSvc.IJwt {
	return &JwtUtil{}
}

func (jwtUtil *JwtUtil) TempTokenForOtpVerification(securityKey string, email string) (string, error) {
	claims := jwt.MapClaims{
		"email": email,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(securityKey))
	if err != nil {
		fmt.Println(err, "error at creating jwt token ")
	}
	return tokenString, err
}

func (jwtUtil *JwtUtil) GenerateRefreshToken(secretKey string) (string, error) {

	claims := jwt.MapClaims{
		"exp": time.Now().Unix() + 604800,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	signedToken, err := token.SignedString([]byte(secretKey))
	if err != nil {
		fmt.Println("Error occured while creating token:", err)
		return "", err
	}

	return signedToken, nil

}

func (jwtUtil *JwtUtil) GenerateAcessToken(securityKey string, id string) (string, error) {
	claims := jwt.MapClaims{
		"exp": time.Now().Unix() + 3600,
		"id":  id,
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(securityKey))
	if err != nil {
		fmt.Println(err, "Error creating acesss token ")
		return "", err
	}
	return tokenString, nil
}

func (jwtUtil *JwtUtil) UnbindEmailFromClaim(tokenString string, tempVerificationKey string) (string, error) {

	secret := []byte(tempVerificationKey)
	parsedToken, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil || !parsedToken.Valid {
		fmt.Println(err)
		return "", err
	}

	claims := parsedToken.Claims.(jwt.MapClaims)
	email := claims["email"].(string)

	return email, nil
}

func (jwtUtil *JwtUtil) VerifyRefreshToken(accesToken string, securityKey string) error {
	key := []byte(securityKey)
	_, err := jwt.Parse(accesToken, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	if err != nil {
		NewResp := err.Error() + ":RefreshToken"
		fmt.Println("-----------", NewResp)
		return errors.New(NewResp)
	}
	return nil
}

func (jwtUtil *JwtUtil) VerifyAccessToken(token string, secretkey string) (string, error) {
	key := []byte(secretkey)
	parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		return key, nil
	})
	// If parsing failed, check the specific error and handle accordingly
	if err != nil {
		if ve, ok := err.(*jwt.ValidationError); ok {
			if ve.Errors&jwt.ValidationErrorMalformed != 0 {
				// Token is malformed
				fmt.Println("malformed token")
				return "", errors.New("malformed token")
			} else if ve.Errors&(jwt.ValidationErrorExpired|jwt.ValidationErrorNotValidYet) != 0 {
				claims, ok := parsedToken.Claims.(jwt.MapClaims)
				if !ok {
					fmt.Println("failed to extract token claims")
					return "", errors.New("failed to extract claims")
				}

				id, ok := claims["id"].(string)
				if !ok {
					fmt.Println("id calim not found or not a string")
					return "", errors.New("ID claim not found or not a string")
				}

				// Token is expired or not valid yet
				fmt.Println("token expired")
				return id, errors.New("expired token")
			} else {
				// Other validation errors
				fmt.Println("validation error")
				return "", errors.New("validation error")
			}
		} else {
			// Other parsing errors
			fmt.Println("other error:", err)
			return "", err
		}
	}

	// If the token is valid, extract claims and return the ID
	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		fmt.Println("token valid,but failed to extract claims")
		return "", errors.New("failed to extract claims")
	}

	id, ok := claims["id"].(string)
	if !ok {
		fmt.Println("token valid,id claim not found or not a string")
		return "", errors.New("ID claim not found or not a string")
	}

	return id, nil
}
