package interfaceRepository_authSvc

import (
	"time"

	requestmodels_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/models/requestmodels"
	responsemodels_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/models/responsemodels"
)

type IUserRepo interface {
	IsUserExist(email string) bool
	IsUserExistWithSameUserName(username string) bool
	DeleteRecentOtpRequestsBefore5min() error
	TemporarySavingUserOtp(otp int, userEmail string, expiration time.Time) error
	CreateUser(userData *requestmodels_authSvc.UserSignUpReq) error

	GetOtpInfo(string) (string, time.Time, error)
	ChangeUserStatusActive(email string) error
	GetUserId(email string) (string, error)

	GetHashPassAndStatus(email string) (string, string, string, error)

	UpdateUserPassword(email *string, hashedPassword *string) error

	GetUserStatForGeneratingAccessToken(userId *string) (*string, error)

	GetUserDataLite(userId *string) (*responsemodels_authSvc.UserProfile, error)

	UpdateUserDetails(editInput *requestmodels_authSvc.EditUserProfile) error

	GetUserProfileURLAndUserName(userId *string) (*responsemodels_authSvc.UserDataLite, error)

	IsUserExistsByID(userID string) (bool, *error)

	SearchUserByNameOrUserName(myId, searchText, limit, offset *string) (*[]responsemodels_authSvc.UserDataForList, error)

	SetUserProfileImg(userId, imageUrl *string) error

	GetFollowersDetails(userIds *[]uint64) (*[]responsemodels_authSvc.UserDataForList, error)
	GetFollowingsDetails(userIds *[]uint64) (*[]responsemodels_authSvc.UserDataForList, error)
}
