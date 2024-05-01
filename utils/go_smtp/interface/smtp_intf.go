package interface_smtp_authSvc

type ISmtp interface {
	SendVerificationEmailWithOtp(otp int, recieverEmail string, recieverName string) error
	SendRestPasswordEmailOtp(otp int, recieverEmail string) error
}
