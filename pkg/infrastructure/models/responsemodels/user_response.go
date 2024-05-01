package responsemodels_authSvc

type SignupData struct {
	Name            string
	UserName        string
	Email           string
	Password        string
	OTP             string
	Token           string
	ConfirmPassword string
	IsUserExist     string
}

type OtpVerifResult struct {
	Otp          string
	AccessToken  string
	RefreshToken string
}

type UserLoginRes struct {
	AccessToken  string
	RefreshToken string
}

type UserProfile struct {
	UserId uint ` gorm:"column:id"`

	Name              string
	UserName          string
	Bio               string
	Links             string
	UserProfileImgURL string
	PostsCount        uint
	FollowersCount    uint
	FollowingCount    uint
	//for userB only
	FollowedBy      string
	FollowingStatus bool
}
