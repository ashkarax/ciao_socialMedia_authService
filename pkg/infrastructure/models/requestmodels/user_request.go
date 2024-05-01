package requestmodels_authSvc

type UserSignUpReq struct {
	Name            string
	UserName        string
	Email           string
	Password        string
	ConfirmPassword string
}

type UserLoginReq struct {
	Email    string
	Password string
}

type ForgotPasswordData struct {
	Otp             string
	Password        string
	ConfirmPassword string
}

type EditUserProfile struct {
	Name     string
	UserName string
	Bio      string
	Links    string
	UserId   string
}
