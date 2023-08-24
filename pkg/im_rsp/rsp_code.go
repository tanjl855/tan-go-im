package im_rsp

const (
	SUCCESS = 0
	ERROR   = -1

	ParameterValidationError = 1
	UserEmailExist           = 2
	UserEmailRegistering     = 3
	ProcessError             = 4
	VerifyCodeError          = 5
	UserNotExist             = 6
	UserPasswordError        = 7
	VerifyPassword           = 8
	AuthError                = 9
	GroupNotExist            = 10
	PasswordError            = 11
)
