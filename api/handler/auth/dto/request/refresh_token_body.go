package request

type RefreshTokenBody struct {
	RefreshToken string `json:"refreshToken" validate:"required"`
}
