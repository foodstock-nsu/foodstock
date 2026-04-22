package dto

type AdminAuthInput struct {
	Login    string
	Password string
}

type AdminAuthOutput struct {
	Token string
}
