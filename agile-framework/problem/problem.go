package problem

import (
	"encoding/json"
	"net/http"
)

const (
	ContentTypeJSON = "application/problem+json"
)

type Option interface {
	apply(*Problem)
}

var (
	ErrConflict                    = Of(http.StatusConflict)
	ErrNotFound                    = Of(http.StatusNotFound)
	ErrUnauthorized                = Of(http.StatusUnauthorized)
	ErrForbidden                   = Of(http.StatusForbidden)
	ErrMethodNotAllowed            = Of(http.StatusMethodNotAllowed)
	ErrStatusRequestEntityTooLarge = Of(http.StatusRequestEntityTooLarge)
	ErrTooManyRequests             = Of(http.StatusTooManyRequests)
	ErrBadRequest                  = Of(http.StatusBadRequest)
	ErrBadGateway                  = Of(http.StatusBadGateway)
	ErrInternalServerError         = Of(http.StatusInternalServerError)
	ErrRequestTimeout              = Of(http.StatusRequestTimeout)
	ErrServiceUnavailable          = Of(http.StatusServiceUnavailable)
	ErrStatusUnprocessableEntity   = Of(http.StatusUnprocessableEntity)
)

type optionFunc func(*Problem)

func (f optionFunc) apply(problem *Problem) { f(problem) }

type Problem struct {
	data   map[string]interface{}
	reason error
}

func (p Problem) JSON() []byte {
	b, _ := p.MarshalJSON()
	return b
}

func (p Problem) UnmarshalJSON(b []byte) error {
	return json.Unmarshal(b, &p.data)
}

func (p Problem) MarshalJSON() ([]byte, error) {
	return json.Marshal(&p.data)
}

func (p Problem) JSONString() string {
	return string(p.JSON())
}

func (p Problem) Error() string {
	return p.JSONString()
}

func (p Problem) Is(err error) bool {
	return p.Error() == err.Error()
}

func (p Problem) Unwrap() error {
	return p.reason
}

func (p Problem) WriteTo(w http.ResponseWriter) (int, error) {
	w.Header().Set("Content-Type", ContentTypeJSON)
	if statuscode, ok := p.data["status"]; ok {
		if statusint, ok := statuscode.(int); ok {
			w.WriteHeader(statusint)
		}
	}
	return w.Write(p.JSON())
}

func New(opts ...Option) *Problem {
	problem := &Problem{}
	problem.data = make(map[string]any)
	for _, opt := range opts {
		opt.apply(problem)
	}
	return problem
}

func Of(statusCode int) *Problem {
	return New(Status(statusCode), Title(http.StatusText(statusCode)))
}

func (p *Problem) Append(opts ...Option) *Problem {
	for _, opt := range opts {
		opt.apply(p)
	}
	return p
}

func Wrap(err error) Option {
	return optionFunc(func(problem *Problem) {
		problem.reason = err
		problem.data["reason"] = err.Error()
	})
}

func Type(uri string) Option {
	return optionFunc(func(problem *Problem) {
		problem.data["type"] = uri
	})
}

func Title(title string) Option {
	return optionFunc(func(problem *Problem) {
		problem.data["title"] = title
	})
}

func Status(status int) Option {
	return optionFunc(func(problem *Problem) {
		problem.data["status"] = status
	})
}

func Detail(detail string) Option {
	return optionFunc(func(problem *Problem) {
		problem.data["detail"] = detail
	})
}

func Instance(uri string) Option {
	return optionFunc(func(problem *Problem) {
		problem.data["instance"] = uri
	})
}

func Custom(key string, value any) Option {
	return optionFunc(func(problem *Problem) {
		problem.data[key] = value
	})
}
