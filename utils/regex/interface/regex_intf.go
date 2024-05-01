package interface_regex_authSvc

type IRegexUtil interface {
	IsValidUsername(username string) (bool, string)
	IsValidPassword(password string) (bool, string)
}
