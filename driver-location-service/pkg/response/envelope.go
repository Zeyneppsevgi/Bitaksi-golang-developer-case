package response

type ErrorBody struct {
	Code    string      `json:"code"`
	Message string      `json:"message"`
	Details interface{} `json:"details,omitempty"`
}

type Envelope struct {
	OK    bool        `json:"ok"`
	Data  interface{} `json:"data,omitempty"`
	Error *ErrorBody  `json:"error,omitempty"`
}

func Success(data interface{}) Envelope {
	return Envelope{OK: true, Data: data}
}

func Failure(code, message string, details interface{}) Envelope {
	return Envelope{OK: false, Error: &ErrorBody{Code: code, Message: message, Details: details}}
}
