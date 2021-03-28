package apierr

type APIError string

func (a APIError) Error() string {
	return string(a)
}

const (
	ErrEmailAlreadyExists APIError = "そのメールアドレスはすでに登録されています"
)
