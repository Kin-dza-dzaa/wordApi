package apierror

import "encoding/json"

type Response struct {
	Result     string `json:"string"`
	Message    string `json:"message"`
	StatusCode int    `json:"code"`
}

func (response *Response) Error() string {
	return response.Message
}

func (response *Response) Marshal() []byte {
	if bytes, err := json.Marshal(response); err != nil {
		return nil
	} else {
		return bytes
	}
}

func NewResponse(ressult, Message string, statusCode int) *Response {
	return &Response{
		Result:     ressult,
		Message:    Message,
		StatusCode: statusCode,
	}
}
