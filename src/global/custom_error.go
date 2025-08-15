package global

type SystemError struct {
	Code    int64  `json:"code"`
	Message string `json:"message"`
}

func (system SystemError) Error() string {
	return system.Message
}
