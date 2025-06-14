package http

type MoveRequest struct {
	RightDutyCycle *float64 `json:"rightDutyCycle" binding:"required,gte=-1,lte=1"`
	LeftDutyCycle  *float64 `json:"leftDutyCycle" binding:"required,gte=-1,lte=1"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}
