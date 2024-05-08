package usecase_authSvc

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	config_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/config"
	requestmodels_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/models/requestmodels"
	responsemodels_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/models/responsemodels"
	"github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/pb"
	interfaceRepository_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/repository/interface"
	interfaceUseCase_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/usecase/interface"
	interface_smtp_authSvc "github.com/ashkarax/ciao_socialMedia_authService/utils/go_smtp/interface"
	interface_hash_authSvc "github.com/ashkarax/ciao_socialMedia_authService/utils/hash_password/interface"
	interface_jwt_authSvc "github.com/ashkarax/ciao_socialMedia_authService/utils/jwt.go/interface"
	interface_randnumgene_authSvc "github.com/ashkarax/ciao_socialMedia_authService/utils/random_number/interface"
	interface_regex_authSvc "github.com/ashkarax/ciao_socialMedia_authService/utils/regex/interface"
)

type UserUseCase struct {
	UserRepo         interfaceRepository_authSvc.IUserRepo
	SmtpUtil         interface_smtp_authSvc.ISmtp
	JwtUtil          interface_jwt_authSvc.IJwt
	RandNumUtil      interface_randnumgene_authSvc.IRandGene
	RegexUtli        interface_regex_authSvc.IRegexUtil
	tokenSecurityKey *config_authSvc.Token
	HashUtil         interface_hash_authSvc.IhashPassword
	PostNrelClient   pb.PostNrelServiceClient
}

func NewUserUseCase(userRepo interfaceRepository_authSvc.IUserRepo,
	smtpUtil interface_smtp_authSvc.ISmtp,
	jwtUtil interface_jwt_authSvc.IJwt,
	randNumUtil interface_randnumgene_authSvc.IRandGene,
	regexUtli interface_regex_authSvc.IRegexUtil,
	config *config_authSvc.Token,
	hashUtil interface_hash_authSvc.IhashPassword,
	postNrelClient *pb.PostNrelServiceClient) interfaceUseCase_authSvc.IUserUseCase {
	return &UserUseCase{UserRepo: userRepo,
		SmtpUtil:         smtpUtil,
		JwtUtil:          jwtUtil,
		RandNumUtil:      randNumUtil,
		RegexUtli:        regexUtli,
		tokenSecurityKey: config,
		HashUtil:         hashUtil,
		PostNrelClient:   *postNrelClient,
	}
}

func (r *UserUseCase) UserSignUp(userData *requestmodels_authSvc.UserSignUpReq) (responsemodels_authSvc.SignupData, error) {

	var resSignUp responsemodels_authSvc.SignupData

	if isUserExist := r.UserRepo.IsUserExist(userData.Email); isUserExist {
		return resSignUp, errors.New("user exists, try again with another email id")
	}

	stat, message := r.RegexUtli.IsValidPassword(userData.Password)
	if !stat {
		return resSignUp, errors.New(message)
	}

	stat, message = r.RegexUtli.IsValidUsername(userData.UserName)
	if !stat {
		return resSignUp, errors.New(message)
	}
	if isUserExistUserName := r.UserRepo.IsUserExistWithSameUserName(userData.UserName); isUserExistUserName {
		return resSignUp, errors.New("user exists, try again with another username")
	}

	errRemv := r.UserRepo.DeleteRecentOtpRequestsBefore5min()
	if errRemv != nil {
		return resSignUp, errRemv
	}

	otp := r.RandNumUtil.RandomNumber()
	errOtp := r.SmtpUtil.SendVerificationEmailWithOtp(otp, userData.Email, userData.Name)
	if errOtp != nil {
		return resSignUp, errOtp
	}

	expiration := time.Now().Add(5 * time.Minute)

	errTempSave := r.UserRepo.TemporarySavingUserOtp(otp, userData.Email, expiration)
	if errTempSave != nil {
		fmt.Println("Cant save temporary data for otp verification in db")
		return resSignUp, errors.New("OTP verification down,please try after some time")
	}

	hashedPassword := r.HashUtil.HashPassword(userData.ConfirmPassword)
	userData.Password = hashedPassword

	errCreateUsr := r.UserRepo.CreateUser(userData)
	if errCreateUsr != nil {
		return resSignUp, errCreateUsr
	}

	tempToken, err := r.JwtUtil.TempTokenForOtpVerification(r.tokenSecurityKey.TempVerificationKey, userData.Email)
	if err != nil {
		fmt.Println("error creating temp token for otp verification")
		return resSignUp, errors.New("error creating temp token for otp verification")
	}

	resSignUp.Token = tempToken

	return resSignUp, nil

}

func (r *UserUseCase) VerifyOtp(otp string, TempVerificationToken *string) (responsemodels_authSvc.OtpVerifResult, error) {
	var otpveriRes responsemodels_authSvc.OtpVerifResult

	email, unbindErr := r.JwtUtil.UnbindEmailFromClaim(*TempVerificationToken, r.tokenSecurityKey.TempVerificationKey)
	if unbindErr != nil {
		return otpveriRes, unbindErr
	}

	userOTP, expiration, errGetInfo := r.UserRepo.GetOtpInfo(email)
	if errGetInfo != nil {
		return otpveriRes, errGetInfo
	}

	if otp != userOTP {
		return otpveriRes, errors.New("invalid OTP")
	}
	if time.Now().After(expiration) {
		return otpveriRes, errors.New("OTP expired")
	}

	changeStatErr := r.UserRepo.ChangeUserStatusActive(email)
	if changeStatErr != nil {
		return otpveriRes, changeStatErr
	}

	userId, fetchErr := r.UserRepo.GetUserId(email)
	if fetchErr != nil {
		return otpveriRes, fetchErr
	}

	accessToken, aTokenErr := r.JwtUtil.GenerateAcessToken(r.tokenSecurityKey.UserSecurityKey, userId)
	if aTokenErr != nil {
		otpveriRes.AccessToken = aTokenErr.Error()
		return otpveriRes, aTokenErr
	}
	refreshToken, rTokenErr := r.JwtUtil.GenerateRefreshToken(r.tokenSecurityKey.UserSecurityKey)
	if rTokenErr != nil {
		otpveriRes.RefreshToken = rTokenErr.Error()
		return otpveriRes, rTokenErr
	}

	otpveriRes.Otp = "verified"
	otpveriRes.AccessToken = accessToken
	otpveriRes.RefreshToken = refreshToken

	return otpveriRes, nil
}

func (r *UserUseCase) UserLogin(loginData *requestmodels_authSvc.UserLoginReq) (responsemodels_authSvc.UserLoginRes, error) {
	var resLogin responsemodels_authSvc.UserLoginRes

	stat, message := r.RegexUtli.IsValidPassword(loginData.Password)
	if !stat {
		return resLogin, errors.New(message)
	}

	hashedPassword, userId, status, errr := r.UserRepo.GetHashPassAndStatus(loginData.Email)
	if errr != nil {
		return resLogin, errr
	}

	passwordErr := r.HashUtil.CompairPassword(hashedPassword, loginData.Password)
	if passwordErr != nil {
		return resLogin, passwordErr
	}

	if status == "blocked" {
		return resLogin, errors.New("user is blocked by the admin")
	}

	if status == "pending" {
		return resLogin, errors.New("user is on status pending,OTP not verified")
	}

	accessToken, err := r.JwtUtil.GenerateAcessToken(r.tokenSecurityKey.UserSecurityKey, userId)
	if err != nil {
		return resLogin, err
	}

	refreshToken, err := r.JwtUtil.GenerateRefreshToken(r.tokenSecurityKey.UserSecurityKey)
	if err != nil {
		return resLogin, err
	}

	resLogin.AccessToken = accessToken
	resLogin.RefreshToken = refreshToken
	return resLogin, nil

}

func (r *UserUseCase) ForgotPasswordRequest(email *string) (*string, error) {

	_, _, status, err := r.UserRepo.GetHashPassAndStatus(*email)
	if err != nil {
		return nil, err
	}

	if status == "blocked" {
		return nil, errors.New("user is blocked by the admin")
	}

	if status == "pending" {
		return nil, errors.New("user is on status pending,OTP not verified")
	}

	err = r.UserRepo.DeleteRecentOtpRequestsBefore5min()
	if err != nil {
		return nil, err
	}

	otp := r.RandNumUtil.RandomNumber()
	err = r.SmtpUtil.SendRestPasswordEmailOtp(otp, *email)
	if err != nil {
		return nil, err
	}

	expiration := time.Now().Add(5 * time.Minute)

	errTempSave := r.UserRepo.TemporarySavingUserOtp(otp, *email, expiration)
	if errTempSave != nil {
		fmt.Println("Cant save temporary data for otp verification in db")
		return nil, errors.New("OTP verification down,please try after some time")
	}

	tempToken, err := r.JwtUtil.TempTokenForOtpVerification(r.tokenSecurityKey.TempVerificationKey, *email)
	if err != nil {
		fmt.Println("----------", err)
		return nil, errors.New("error creating temp token for otp verification")
	}

	return &tempToken, nil
}

func (r *UserUseCase) ResetPassword(userData *requestmodels_authSvc.ForgotPasswordData, TempVerificationToken *string) error {

	stat, message := r.RegexUtli.IsValidPassword(userData.Password)
	if !stat {
		return errors.New(message)
	}

	email, err := r.JwtUtil.UnbindEmailFromClaim(*TempVerificationToken, r.tokenSecurityKey.TempVerificationKey)
	if err != nil {
		return err
	}

	userOTP, expiration, err := r.UserRepo.GetOtpInfo(email)
	if err != nil {
		return err
	}

	if userData.Otp != userOTP {
		return errors.New("invalid OTP")
	}
	if time.Now().After(expiration) {
		return errors.New("OTP expired")
	}

	hashedPassword := r.HashUtil.HashPassword(userData.ConfirmPassword)

	err = r.UserRepo.UpdateUserPassword(&email, &hashedPassword)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserUseCase) UserProfile(userId, userBId *string) (*responsemodels_authSvc.UserProfile, error) {
	var actualId *string

	if *userBId == "" {
		actualId = userId
	} else {
		actualId = userBId
	}
	userData, err := r.UserRepo.GetUserDataLite(actualId)
	if err != nil {
		return nil, err
	}

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	respData, err := r.PostNrelClient.GetCountsForUserProfile(context, &pb.RequestUserIdPnR{
		UserId: *actualId,
	})
	if err != nil {
		log.Fatal(err)
	}
	if respData.ErrorMessage != "" {
		return nil, errors.New(respData.ErrorMessage)
	}

	userData.PostsCount = uint(respData.PostCount)
	userData.FollowersCount = uint(respData.FollowerCount)
	userData.FollowingCount = uint(respData.FollowingCount)

	if *userBId != "" {

		respStat, err := r.PostNrelClient.UserAFollowingUserBorNot(context, &pb.RequestFollowUnFollow{
			UserId:  *userId,
			UserBId: *userBId,
		})
		if err != nil {
			log.Fatal(err)
		}
		if respData.ErrorMessage != "" {
			return nil, errors.New(respData.ErrorMessage)
		}
		userData.FollowingStatus = respStat.BoolStat
	}

	return userData, nil
}

func (r *UserUseCase) EditUserDetails(editInput *requestmodels_authSvc.EditUserProfile) error {

	stat, message := r.RegexUtli.IsValidUsername(editInput.UserName)
	if !stat {
		return errors.New(message)
	}

	userData, err := r.UserRepo.GetUserDataLite(&editInput.UserId)
	if err != nil {
		fmt.Println("-------", err)
		return err
	}

	if userData.UserName != editInput.UserName {

		if isUserExistUserName := r.UserRepo.IsUserExistWithSameUserName(editInput.UserName); isUserExistUserName {
			return errors.New("user exists, try again with another username")
		}
	}

	err = r.UserRepo.UpdateUserDetails(editInput)
	if err != nil {
		return err
	}

	return nil
}

func (r *UserUseCase) GetUserDetailsLiteForPostView(userId *string) (*responsemodels_authSvc.UserDataLite, error) {

	respData, err := r.UserRepo.GetUserProfileURLAndUserName(userId)
	if err != nil {
		return nil, err
	}

	return respData, nil
}

func (r *UserUseCase) CheckUserExist(userId *string) (bool, *error) {

	boolStat, err := r.UserRepo.IsUserExistsByID(*userId)
	if err != nil {
		return boolStat, err
	}
	return boolStat, nil
}

func (r *UserUseCase) GetFollowersDetails(userId *string) (*[]responsemodels_authSvc.UserDataForList, *error) {

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	userIdsSlice, err := r.PostNrelClient.GetFollowersIds(context, &pb.RequestUserIdPnR{UserId: *userId})
	if err != nil {
		log.Fatal(err)
	}
	if userIdsSlice.ErrorMessage != "" {
		return nil, &err
	}

	userDetailsSlice, err := r.UserRepo.GetFollowersDetails(&userIdsSlice.UserIds)
	if err != nil {
		return nil, &err
	}

	return userDetailsSlice, nil
}
func (r *UserUseCase) GetFollowingsDetails(userId *string) (*[]responsemodels_authSvc.UserDataForList, *error) {

	context, cancel := context.WithTimeout(context.Background(), time.Second*10)
	defer cancel()

	userIdsSlice, err := r.PostNrelClient.GetFollowingsIds(context, &pb.RequestUserIdPnR{UserId: *userId})
	if err != nil {
		log.Fatal(err)
	}
	if userIdsSlice.ErrorMessage != "" {
		return nil, &err
	}

	userDetailsSlice, err := r.UserRepo.GetFollowingsDetails(&userIdsSlice.UserIds)
	if err != nil {
		return nil, &err
	}

	return userDetailsSlice, nil
}
