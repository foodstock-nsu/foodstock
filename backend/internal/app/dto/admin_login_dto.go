package dto

type AdminLoginInput struct {
	Login    string
	Password string
}

type AdminLoginOutput struct {
	Token string
}
