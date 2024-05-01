package usecase_authSvc

import (
	"errors"

	config_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/config"
	interfaceRepository_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/repository/interface"
	interfaceUseCase_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/usecase/interface"
	interface_jwt_authSvc "github.com/ashkarax/ciao_socialMedia_authService/utils/jwt.go/interface"
)

type JwtUseCase struct {
	JwtKeys  *config_authSvc.Token
	JwtUtil  interface_jwt_authSvc.IJwt
	UserRepo interfaceRepository_authSvc.IUserRepo
}

func NewJwtUseCase(jwtKeys *config_authSvc.Token,
	jwtUtil interface_jwt_authSvc.IJwt,
	userRepo interfaceRepository_authSvc.IUserRepo) interfaceUseCase_authSvc.IJwtUseCase {
	return &JwtUseCase{JwtKeys: jwtKeys,
		JwtUtil:  jwtUtil,
		UserRepo: userRepo,
	}
}

func (r *JwtUseCase) VerifyAccessToken(token *string) (*string, error) {

	userId, err := r.JwtUtil.VerifyAccessToken(*token, r.JwtKeys.UserSecurityKey)
	if err != nil {
		if userId == "" {
			return nil, err
		}
		return nil, err
	}
	return &userId, nil
}

func (r *JwtUseCase) AccessRegenerator(accessToken *string, refreshToken *string) (*string, error) {

	userId, err := r.JwtUtil.VerifyAccessToken(*accessToken, r.JwtKeys.UserSecurityKey)
	if err != nil {
		if userId == "" {
			return nil, err
		}
	}

	err = r.JwtUtil.VerifyRefreshToken(*refreshToken, r.JwtKeys.UserSecurityKey)
	if err != nil {
		return nil, err
	}

	status, err := r.UserRepo.GetUserStatForGeneratingAccessToken(&userId)
	if err != nil || *status == "blocked" {
		if *status == "blocked" {
			return nil, errors.New("user is in blocked status")
		}

		return nil, err
	}

	newAcessToken, err := r.JwtUtil.GenerateAcessToken(r.JwtKeys.UserSecurityKey, userId)
	if err != nil {
		return nil, err
	}

	return &newAcessToken, nil

}
