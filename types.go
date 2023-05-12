package main

type userLoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type otpResponse struct {
	OTP string `json:"otp"`
}
