package entity

// User ...
type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
	ID       int    `json:"id"`
}

// ---

// EmptyRequest ...
type EmptyRequest struct{}

// IDRequest ...
type IDRequest struct {
	ID int `json:"id"`
}

// UsernamePasswordRequest ...
type UsernamePasswordRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// UsernameRequest ...
type UsernameRequest struct {
	Username string `json:"username"`
}

// UsernamePasswordEmailRequest ...
type UsernamePasswordEmailRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Email    string `json:"email"`
}

// ---

// UsersErrorResponse ...
type UsersErrorResponse struct {
	Err   string `json:"err,omitempty"`
	Users []User `json:"users"`
}

// UserErrorResponse ...
type UserErrorResponse struct {
	Err  string `json:"err,omitempty"`
	User User   `json:"user"`
}

// IDErrorResponse ...
type IDErrorResponse struct {
	Err string `json:"err,omitempty"`
	ID  int    `json:"id"`
}

/* // ErrorResponse ...
type ErrorResponse struct {
	Err string `json:"err,omitempty"`
} */

// RowsErrorResponse ...
type RowsErrorResponse struct {
	Err          string `json:"err,omitempty"`
	RowsAffected int    `json:"rowsAffected"`
}
