package auth

type RegisterCommand struct {
	Email    string
	Password string
	Username string
}

type LoginCommand struct {
	Email    string
	Password string
}
