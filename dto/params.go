package dto

type ResponseParam struct {
	StatusCode int
	Message    string
	Data       any
	Pagination *Pagination
}

type ParamRequest struct {
	Search     string
	UserID     int
	Pagination Pagination
}
