package global

import "fmt"

const errorCodeBase = 0

const (
	InvalidJSONString int64 = errorCodeBase + 2
	InvalidUserToken  int64 = errorCodeBase + 4

	DatabaseError int64 = errorCodeBase + 8
	IncorrectPin  int64 = errorCodeBase + 9
)

var ErrorMessage = map[int64]string{
	InvalidUserToken: "cannot get user_id from token",
	IncorrectPin:     "Incorrect Pin",
}

func GetErrorMessage(code int64, args ...interface{}) string {
	return fmt.Sprintf(ErrorMessage[code], args...)
}
