package models

/*Theses models are used for swagger-documentation purpose*/

// swagger:model LoginRequestBody
type LoginRequestBody struct {
	// The Username of a user
	// example: "pranay"
	Username string `json:"username" validate:"required,min=3,max=15" example:"userABC"`

	// The Password of a user
	// example: "pranay"
	Password string `json:"password" validate:"required,min=4" example:"pass123"`
}

// swagger:model CreateUserBody
type CreateUserBody struct {
	// The Username of a user
	// example: "pranay"
	Username string `json:"username" validate:"required,min=3,max=15" example:"userABC"`

	// The Password of a user
	// example: abcde
	Password string `json:"password" validate:"required,min=4" example:"pass123"`

	// The Usertype of a user
	// example: "USER" || "ADMIN"
	Usertype string `json:"usertype" validate:"required,eq=USER|eq=ADMIN" example:"{'USER'/'ADMIN'}"`
}

// swagger:model RefreshRotateResponse
type RefreshRotateResponse struct {
	// The jwttoken of a user
	JwtToken string `json:"jwttoken"`

	// The refreshtoken of a user
	RefreshToken string `json:"refreshtoken"`
}

// swagger:model CreateUserResponse
type CreateUserResponse struct {
	// The insertion number of new user
	Result_insertion_number string `json:"result_insertion_number"`
}
