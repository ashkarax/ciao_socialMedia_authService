package server_authSvc

import (
	"context"

	requestmodels_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/models/requestmodels"
	"github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/pb"
	interfaceUseCase_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/usecase/interface"
)

type AuthService struct {
	userUseCase interfaceUseCase_authSvc.IUserUseCase
	jwtUseCase  interfaceUseCase_authSvc.IJwtUseCase
	pb.AuthServiceServer
}

func NewAuthServiceServer(userUseCase interfaceUseCase_authSvc.IUserUseCase,
	jwtUseCase interfaceUseCase_authSvc.IJwtUseCase) *AuthService {
	return &AuthService{userUseCase: userUseCase,
		jwtUseCase: jwtUseCase,
	}
}

func (u *AuthService) UserSignUp(ctx context.Context, req *pb.SignUpRequest) (*pb.SignUpResponse, error) {

	var inputData requestmodels_authSvc.UserSignUpReq

	inputData.Name = req.Name
	inputData.UserName = req.UserName
	inputData.Email = req.Email
	inputData.Password = req.Password
	inputData.ConfirmPassword = req.ConfirmPassword

	respData, err := u.userUseCase.UserSignUp(&inputData)
	if err != nil {
		return &pb.SignUpResponse{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.SignUpResponse{
		Token: respData.Token,
	}, nil

}

func (u *AuthService) UserOTPVerication(ctx context.Context, req *pb.RequestOtpVefification) (*pb.ResponseOtpVerification, error) {

	respData, err := u.userUseCase.VerifyOtp(req.Otp, &req.TempToken)
	if err != nil {
		return &pb.ResponseOtpVerification{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseOtpVerification{
		Otp:          respData.Otp,
		AccessToken:  respData.AccessToken,
		RefreshToken: respData.RefreshToken,
	}, nil

}

func (u *AuthService) UserLogin(ctx context.Context, req *pb.RequestUserLogin) (*pb.ResponseUserLogin, error) {

	var loginData requestmodels_authSvc.UserLoginReq

	loginData.Email = req.Email
	loginData.Password = req.Password

	respData, err := u.userUseCase.UserLogin(&loginData)
	if err != nil {
		return &pb.ResponseUserLogin{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseUserLogin{
		AccessToken:  respData.AccessToken,
		RefreshToken: respData.RefreshToken,
	}, nil

}

func (u *AuthService) ForgotPasswordRequest(ctx context.Context, req *pb.RequestForgotPass) (*pb.ResponseForgotPass, error) {

	respData, err := u.userUseCase.ForgotPasswordRequest(&req.Email)
	if err != nil {
		return &pb.ResponseForgotPass{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseForgotPass{
		Token: *respData,
	}, nil

}

func (u *AuthService) ResetPassword(ctx context.Context, req *pb.RequestResetPass) (*pb.ResponseErrorMessage, error) {

	var requestData requestmodels_authSvc.ForgotPasswordData

	requestData.Otp = req.Otp
	requestData.Password = req.Password
	requestData.ConfirmPassword = req.ConfirmPassword

	err := u.userUseCase.ResetPassword(&requestData, &req.TempToken)

	if err != nil {
		return &pb.ResponseErrorMessage{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseErrorMessage{}, nil

}

func (u *AuthService) VerifyAccessToken(ctx context.Context, req *pb.RequestVerifyAccess) (*pb.ResponseVerifyAccess, error) {

	userId, err := u.jwtUseCase.VerifyAccessToken(&req.AccessToken)

	if err != nil {
		return &pb.ResponseVerifyAccess{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseVerifyAccess{
		UserId: *userId,
	}, nil
}

func (u *AuthService) AccessRegenerator(ctx context.Context, req *pb.RequestAccessGenerator) (*pb.ResponseAccessGenerator, error) {

	newAccessToken, err := u.jwtUseCase.AccessRegenerator(&req.AccessToken, &req.RefreshToken)
	if err != nil {
		return &pb.ResponseAccessGenerator{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseAccessGenerator{
		AccesToken: *newAccessToken,
	}, nil

}

func (u *AuthService) GetUserProfile(ctx context.Context, req *pb.RequestGetUserProfile) (*pb.ResponseUserProfile, error) {

	respData, err := u.userUseCase.UserProfile(&req.UserId, &req.UserBId)
	if err != nil {
		return &pb.ResponseUserProfile{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseUserProfile{
		Name:            respData.Name,
		UserName:        respData.UserName,
		Bio:             respData.Bio,
		Links:           respData.Links,
		ProfileImageURL: respData.UserProfileImgURL,
		PostsCount:      uint64(respData.PostsCount),
		FollowerCount:   uint64(respData.FollowersCount),
		FollowingCount:  uint64(respData.FollowingCount),
		FollowingStat:   respData.FollowingStatus,
	}, nil

}

func (u *AuthService) EditUserProfile(ctx context.Context, req *pb.RequestEditUserProfile) (*pb.ResponseErrorMessage, error) {

	var editInput requestmodels_authSvc.EditUserProfile

	editInput.Name = req.Name
	editInput.UserName = req.UserName
	editInput.Bio = req.Bio
	editInput.Links = req.Links
	editInput.UserId = req.UserId

	err := u.userUseCase.EditUserDetails(&editInput)
	if err != nil {
		return &pb.ResponseErrorMessage{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseErrorMessage{}, nil

}

func (u *AuthService) CheckUserExist(ctx context.Context, req *pb.RequestUserId) (*pb.ResponseBool, error) {
	stat, err := u.userUseCase.CheckUserExist(&req.UserId)
	if err != nil {
		return &pb.ResponseBool{
			ErrorMessage: (*err).Error(),
		}, nil
	}
	return &pb.ResponseBool{
		ExistStatus: stat,
	}, nil
}

func (u *AuthService) GetFollowersDetails(ctx context.Context, req *pb.RequestUserId) (*pb.ResponseGetUsersDetails, error) {
	respData, err := u.userUseCase.GetFollowersDetails(&req.UserId)
	if err != nil {
		return &pb.ResponseGetUsersDetails{
			ErrorMessage: (*err).Error(),
		}, nil
	}

	var respLoader []*pb.SingleResponseGetFollowers
	for i := range *respData {
		respLoader = append(respLoader, &pb.SingleResponseGetFollowers{
			UserId:          (*respData)[i].Id,
			Name:            (*respData)[i].Name,
			UserName:        (*respData)[i].UserName,
			ProfileImageURL: (*respData)[i].ProfileImgUrl})

	}
	return &pb.ResponseGetUsersDetails{
		UserData: respLoader,
	}, nil
}
func (u *AuthService) GetFollowingsDetails(ctx context.Context, req *pb.RequestUserId) (*pb.ResponseGetUsersDetails, error) {
	respData, err := u.userUseCase.GetFollowingsDetails(&req.UserId)
	if err != nil {
		return &pb.ResponseGetUsersDetails{
			ErrorMessage: (*err).Error(),
		}, nil
	}

	var respLoader []*pb.SingleResponseGetFollowers
	for i := range *respData {
		respLoader = append(respLoader, &pb.SingleResponseGetFollowers{
			UserId:          (*respData)[i].Id,
			Name:            (*respData)[i].Name,
			UserName:        (*respData)[i].UserName,
			ProfileImageURL: (*respData)[i].ProfileImgUrl})

	}
	return &pb.ResponseGetUsersDetails{
		UserData: respLoader,
	}, nil
}

func (u *AuthService) GetUserDetailsLiteForPostView(ctx context.Context, req *pb.RequestUserId) (*pb.ResponseUserDetailsLite, error) {

	respData, err := u.userUseCase.GetUserDetailsLiteForPostView(&req.UserId)
	if err != nil {
		return &pb.ResponseUserDetailsLite{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseUserDetailsLite{
		UserName:          respData.UserName,
		UserProfileImgURL: respData.UserProfileImgURL,
	}, nil

}

func (u *AuthService) SearchUser(ctx context.Context, req *pb.RequestUserSearch) (*pb.ResponseUserSearch, error) {

	respData, err := u.userUseCase.SearchUser(&req.UserId, &req.SearchText, &req.Limit, &req.Offset)
	if err != nil {
		return &pb.ResponseUserSearch{
			ErrorMessage: err.Error(),
		}, nil
	}

	var respSlice []*pb.SingleResponseGetFollowers
	for i := range *respData {
		respSlice = append(respSlice, &pb.SingleResponseGetFollowers{
			UserId:          (*respData)[i].Id,
			Name:            (*respData)[i].Name,
			UserName:        (*respData)[i].UserName,
			ProfileImageURL: (*respData)[i].ProfileImgUrl})
	}

	return &pb.ResponseUserSearch{
		SearchResult: respSlice,
	}, nil
}

func (u *AuthService) SetUserProfileImage(ctx context.Context, req *pb.RequestSetProfileImg) (*pb.ResponseErrorMessage, error) {

	err := u.userUseCase.SetProfileImage(&req.UserId, &req.ContentType, &req.Img)
	if err != nil {
		return &pb.ResponseErrorMessage{
			ErrorMessage: err.Error(),
		}, nil
	}

	return &pb.ResponseErrorMessage{}, nil
}
