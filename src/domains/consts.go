package domains

import (
	"path/filepath"
	"runtime"
)

const (
	Success             = "success"
	Failure             = "failure"
	Warning             = "warning"
	ForgotPasswordToken = "forgot_password_token"
	ConfirmProfileToken = "confirm_profile_token"
)
const (
	ExpiredInForgotPasswordToken uint64 = 14400  //seconds = 4 hours
	ExpiredInConfirmProfileToken uint64 = 172800 //seconds = 48 hours = 2 days
	ExpiredInAccessToken         uint64 = 3600   //seconds = 1 hour
	ExpiredInRefreshToken        uint64 = 604800 //seconds = 168 hours = 7 days
)

var (
	_, b, _, _ = runtime.Caller(0)
	RootPath = filepath.Join(filepath.Dir(b), "../..")
)

