package users

type UserLoginReceiveStruct struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type UserRegisterStruct struct {
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
	Telephone string `json:"telephone" binding:"required"`
	AuthCode  string `json:"auth_code" binding:"required"`
}

type UserGetAuthCodeStruct struct {
	Telephone string `json:"telephone" binding:"required"`
}

type UserUpdatePasswordStruct struct {
	Telephone string `json:"telephone" binding:"required"`
	Password  string `json:"password" binding:"required"`
	AuthCode  string `json:"auth_code" binding:"required"`
}
