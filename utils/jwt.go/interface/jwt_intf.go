package interface_jwt_authSvc

type IJwt interface {
	TempTokenForOtpVerification(securityKey string, email string) (string, error)
	GenerateRefreshToken(secretKey string) (string, error)
	GenerateAcessToken(securityKey string, id string) (string, error)
	UnbindEmailFromClaim(tokenString string, tempVerificationKey string) (string, error)
	VerifyRefreshToken(accesToken string, securityKey string) error
	VerifyAccessToken(token string, secretkey string) (string, error)
}
