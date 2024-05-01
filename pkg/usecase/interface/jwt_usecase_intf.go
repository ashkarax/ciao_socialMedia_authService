package interfaceUseCase_authSvc

type IJwtUseCase interface {
	VerifyAccessToken(token *string) (*string, error)
	AccessRegenerator(accessToken *string, refreshToken *string) (*string, error)
}
