package entity

// ErrorResponse ...
type ErrorResponse struct {
	Err string `json:"err,omitempty"`
}

/*
// UsernamePasswordEmailRequest (string, string, string) (string, error).
type UsernamePasswordEmailRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// UsernamePasswordRequest (string, string) (string, error).
type UsernamePasswordRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// TokenRequest (string) error.
type TokenRequest struct {
	Token string `json:"token"`
}

// EmptyRequest () ([]User, error).
type EmptyRequest struct{}

// ---

// TokenErrorResponse (string, string, string) (string, error).
type TokenErrorResponse struct {
	Token string `json:"token"`
	Err   string `json:"err,omitempty"`
}

// UsersErrorResponse () ([]User, error).
type UsersErrorResponse struct {
	Err   string `json:"err,omitempty"`
	Users []User `json:"users"`
}

// UserErrorResponse () (User, error).
type UserErrorResponse struct {
	Err  string `json:"err,omitempty"`
	User User   `json:"user"`
}

// ErrorResponse (string, string, string) (string, error).
type ErrorResponse struct {
	Err string `json:"err,omitempty"`
}

// User ...
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	ID       int    `json:"id"`
} */
