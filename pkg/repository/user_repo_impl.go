package repository_authSvc

import (
	"errors"
	"fmt"
	"time"

	requestmodels_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/models/requestmodels"
	responsemodels_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/infrastructure/models/responsemodels"
	interfaceRepository_authSvc "github.com/ashkarax/ciao_socialMedia_authService/pkg/repository/interface"
	"gorm.io/gorm"
)

type UserRepo struct {
	DB *gorm.DB
}

func NewUserRepo(db *gorm.DB) interfaceRepository_authSvc.IUserRepo {
	return &UserRepo{DB: db}
}

func (d *UserRepo) IsUserExist(email string) bool {
	var userCount int

	delUncompletedUser := "DELETE FROM users WHERE email =$1 AND status =$2"
	result := d.DB.Exec(delUncompletedUser, email, "pending")
	if result.Error != nil {
		fmt.Println("Error in deleting already existing user with this email and status pending")
	}

	query := "SELECT COUNT(*) FROM users WHERE email=$1 AND status!=$2"
	err := d.DB.Raw(query, email, "deleted").Row().Scan(&userCount)
	if err != nil {
		fmt.Println("error in usercount query")
	}
	if userCount >= 1 {
		return true
	}

	return false
}

func (d *UserRepo) IsUserExistWithSameUserName(username string) bool {
	var userCount int

	query := "SELECT COUNT(*) FROM users WHERE user_name=$1 AND status!=$2"
	err := d.DB.Raw(query, username, "deleted").Row().Scan(&userCount)
	if err != nil {
		fmt.Println("error in usercount query")
	}
	if userCount >= 1 {
		return true
	}

	return false
}

func (d *UserRepo) DeleteRecentOtpRequestsBefore5min() error {
	query := "DELETE FROM otp_infos WHERE expiration < CURRENT_TIMESTAMP - INTERVAL '5 minutes';"
	err := d.DB.Exec(query).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *UserRepo) TemporarySavingUserOtp(otp int, userEmail string, expiration time.Time) error {

	query := `INSERT INTO otp_infos (email, otp, expiration) VALUES ($1, $2, $3)`
	err := d.DB.Exec(query, userEmail, otp, expiration).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *UserRepo) CreateUser(userData *requestmodels_authSvc.UserSignUpReq) error {
	query := "INSERT INTO users (name,user_name, email, password) VALUES($1, $2, $3, $4)"
	err := d.DB.Exec(query, userData.Name, userData.UserName, userData.Email, userData.Password).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *UserRepo) GetOtpInfo(email string) (string, time.Time, error) {
	var expiration time.Time
	type OTPInfo struct {
		OTP        string    `gorm:"column:otp"`
		Expiration time.Time `gorm:"column:expiration"`
	}
	var otpInfo OTPInfo
	if err := d.DB.Raw("SELECT otp, expiration FROM otp_infos WHERE email = ? ORDER BY expiration DESC LIMIT 1;", email).Scan(&otpInfo).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return "", expiration, errors.New("otp verification failed, no data found for this user-email")
		}
		return "", expiration, errors.New("otp verification failed, error finding user data: " + err.Error())
	}

	return otpInfo.OTP, otpInfo.Expiration, nil
}

func (d *UserRepo) ChangeUserStatusActive(email string) error {
	query := "UPDATE users SET status = 'active' WHERE email = $1"
	result := d.DB.Exec(query, email)
	if result.Error != nil {
		fmt.Println("", result.Error)

		return result.Error
	}

	return nil
}

func (d *UserRepo) GetUserId(email string) (string, error) {
	var userId string
	query := "SELECT id FROM users WHERE email=$1 AND status=$2"
	err := d.DB.Raw(query, email, "active").Row().Scan(&userId)
	if err != nil {
		fmt.Println("", err)
		return "", err
	}
	return userId, nil

}

func (d *UserRepo) GetHashPassAndStatus(email string) (string, string, string, error) {
	var hashedPassword, status, userid string

	query := "SELECT password, id, status FROM users WHERE email=? AND status!='delete'"
	err := d.DB.Raw(query, email).Row().Scan(&hashedPassword, &userid, &status)
	if err != nil {
		return "", "", "", errors.New("no user exist with the specified email,signup first")
	}

	return hashedPassword, userid, status, nil
}

func (d *UserRepo) UpdateUserPassword(email *string, hashedPassword *string) error {
	query := `UPDATE users SET password=$1 WHERE email=$2`
	err := d.DB.Exec(query, *hashedPassword, *email).Error
	if err != nil {
		return err
	}
	return nil
}

func (d *UserRepo) GetUserStatForGeneratingAccessToken(userId *string) (*string, error) {
	var userCurrentStatus string
	query := "SELECT status from users WHERE id=?"
	result := d.DB.Raw(query, userId).Scan(&userCurrentStatus)

	if result.RowsAffected == 0 {
		errMessage := fmt.Sprintf("No results found,No user with this id=%s found in db", *userId)
		return &userCurrentStatus, errors.New(errMessage)
	}
	if result.Error != nil {
		return &userCurrentStatus, result.Error
	}

	return &userCurrentStatus, nil
}

func (d *UserRepo) GetUserDataLite(userId *string) (*responsemodels_authSvc.UserProfile, error) {
	var resp responsemodels_authSvc.UserProfile
	query := "SELECT id,name,user_name,bio,links,profile_img_url FROM users WHERE id=$1"
	err := d.DB.Raw(query, userId).Scan(&resp)
	if err.Error != nil {
		return &resp, err.Error
	} else if err.RowsAffected == 0 {
		return nil, errors.New("no user with the specified user id found")
	}
	return &resp, nil
}

func (d *UserRepo) UpdateUserDetails(editInput *requestmodels_authSvc.EditUserProfile) error {
	query := "UPDATE users SET name=$1,user_name=$2,bio=$3,links=$4 WHERE id=$5"
	err := d.DB.Exec(query, editInput.Name, editInput.UserName, editInput.Bio, editInput.Links, editInput.UserId).Error
	if err != nil {
		return err
	}

	return nil
}

func (d *UserRepo) GetUserProfileURLAndUserName(userId *string) (*responsemodels_authSvc.UserDataLite, error) {

	var resp responsemodels_authSvc.UserDataLite
	query := "SELECT user_name,profile_img_url FROM users WHERE id=$1"
	err := d.DB.Raw(query, userId).Scan(&resp)
	if err.Error != nil {
		return &resp, err.Error
	} else if err.RowsAffected == 0 {
		return nil, errors.New("no user with the specified user id found")
	}
	return &resp, nil
}

func (d *UserRepo) IsUserExistsByID(userID string) (bool, *error) {
	var userCount int
	query := "SELECT COUNT(*) FROM users WHERE id=$1 AND status!=$2"
	err := d.DB.Raw(query, userID, "deleted").Row().Scan(&userCount)
	if err != nil {
		return false, &err
	}
	if userCount >= 1 {
		return true, nil
	}
	return false, nil
}

func (d *UserRepo) GetFollowersDetails(userIds *[]uint64) (*[]responsemodels_authSvc.UserDataForList, error) {
	var userDatas []responsemodels_authSvc.UserDataForList

	interfaceIds := make([]interface{}, len(*userIds))
	for i, id := range *userIds {
		interfaceIds[i] = id
	}

	query := "SELECT id,name, user_name,profile_img_url FROM users WHERE id IN ("
	for i := range *userIds {
		// Add placeholder for each user ID
		query += "?"
		// Add comma if not the last user ID
		if i < len(*userIds)-1 {
			query += ","
		}
	}
	query += ")"

	err := d.DB.Raw(query, interfaceIds...).Scan(&userDatas).Error
	if err != nil {
		fmt.Println("------------", err)
		return nil, err
	}

	return &userDatas, nil
}

func (d *UserRepo) GetFollowingsDetails(userIds *[]uint64) (*[]responsemodels_authSvc.UserDataForList, error) {
	var userDatas []responsemodels_authSvc.UserDataForList

	interfaceIds := make([]interface{}, len(*userIds))
	for i, id := range *userIds {
		interfaceIds[i] = id
	}

	query := "SELECT id,name, user_name,profile_img_url FROM users WHERE id IN ("
	for i := range *userIds {
		// Add placeholder for each user ID
		query += "?"
		// Add comma if not the last user ID
		if i < len(*userIds)-1 {
			query += ","
		}
	}
	query += ")"

	err := d.DB.Raw(query, interfaceIds...).Scan(&userDatas).Error
	if err != nil {
		fmt.Println("------------", err)
		return nil, err
	}

	return &userDatas, nil
}

func (d *UserRepo) SearchUserByNameOrUserName(myId, searchText, limit, offset *string) (*[]responsemodels_authSvc.UserDataForList, error) {

	var resp []responsemodels_authSvc.UserDataForList

	query := "SELECT id,name,user_name,profile_img_url FROM users WHERE (name ILIKE $1 OR user_name ILIKE $1) AND status='active' AND id != $2 LIMIT $3 OFFSET $4"
	err := d.DB.Raw(query, "%"+*searchText+"%", myId, limit, offset).Scan(&resp).Error
	if err != nil {
		return nil, err
	}
	return &resp, nil
}

func (d *UserRepo) SetUserProfileImg(userId, imageUrl *string) error {
	query := "UPDATE users SET profile_img_url=$1 WHERE id=$2"
	err := d.DB.Exec(query, imageUrl, userId).Error
	if err != nil {
		return err
	}
	return nil
}
