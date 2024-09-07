package dto

type AuthDTO struct {
	Token string
}

func NewAuthDTO(token string) *AuthDTO {
	return &AuthDTO{
		Token: token,
	}
}
