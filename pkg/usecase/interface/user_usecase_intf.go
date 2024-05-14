package interfaceUseCase_authSvc

import (
	requestmodels_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/models/requestmodels"
	responsemodels_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/models/responsemodels"
)

type IUserUseCase interface {
	UserSignUp(userData *requestmodels_authSvc.UserSignUpReq) (responsemodels_authSvc.SignupData, error)
	VerifyOtp(otp string, TempVerificationToken *string) (responsemodels_authSvc.OtpVerifResult, error)
	UserLogin(loginData *requestmodels_authSvc.UserLoginReq) (responsemodels_authSvc.UserLoginRes, error)
	ForgotPasswordRequest(email *string) (*string, error)
	ResetPassword(userData *requestmodels_authSvc.ForgotPasswordData, TempVerificationToken *string) error
	UserProfile(userId, UserBId *string) (*responsemodels_authSvc.UserProfile, error)
	EditUserDetails(editInput *requestmodels_authSvc.EditUserProfile) error

	GetUserDetailsLiteForPostView(userId *string) (*responsemodels_authSvc.UserDataLite, error)
	CheckUserExist(userId *string) (bool, *error)
	GetFollowersDetails(userId *string) (*[]responsemodels_authSvc.UserDataForList, *error)
	GetFollowingsDetails(userId *string) (*[]responsemodels_authSvc.UserDataForList, *error)

	SearchUser(myId, searchText, limit, offset *string) (*[]responsemodels_authSvc.UserDataForList, error)
	SetProfileImage(userId, contentType *string, Img *[]byte) error
}
