package mapper

import (
	httpdto "backend/internal/adapter/in/http/dto"
	appdto "backend/internal/app/dto"
)

func MapRequestToAdminAuth(req httpdto.AdminAuthRequest) appdto.AdminAuthInput {
	return appdto.AdminAuthInput{
		Login:    req.Login,
		Password: req.Password,
	}
}

func MapAdminAuthToResponse(out appdto.AdminAuthOutput) httpdto.AdminAuthResponse {
	return httpdto.AdminAuthResponse{Token: out.Token}
}
