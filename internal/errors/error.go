package errors

type Response struct {
	Code   int    `json:"code"`
	Reason string `json:"reason"`
}

func (r *Response) Error() string {
	return r.Reason
}

func (r *Response) StatusCode() int {
	return r.Code
}
