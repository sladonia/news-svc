package handler

type Level string

func (e Level) String() string {
	return string(e)
}

const (
	LevelUser   Level = "user"
	LevelSystem Level = "system"
)

func NewApiError(msg string, level Level) ApiError {
	return ApiError{
		Error: Error{
			Message: msg,
			Level:   level,
		},
	}
}

type ApiError struct {
	Error Error `json:"error"`
}

type Error struct {
	Message string `json:"message"`
	Level   Level  `json:"level"`
}
