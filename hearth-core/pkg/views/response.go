package views

import (
	"github.com/gin-gonic/gin"
)

type Success struct {
	StatusCode int         `json:"-"`
	Message    string      `json:"message"`
	Data       any `json:"data,omitempty"`
}

func (s *Success) SetStatusCode(code int) *Success {
	s.StatusCode = code
	return s
}

func (s *Success) SetMessage(msg string) *Success {
	s.Message = msg
	return s
}

func (s *Success) SetData(data interface{}) *Success {
	s.Data = data
	return s
}

func (s *Success) Send(c *gin.Context) error {
	if s.StatusCode == 0 {
		s.StatusCode = 200
	}
	c.JSON(s.StatusCode, s)
	return nil
}

type Failure struct {
	StatusCode int    `json:"-"`
	Message    string `json:"message"`
}

func (f *Failure) SetStatusCode(code int) *Failure {
	f.StatusCode = code
	return f
}

func (f *Failure) SetMessage(msg string) *Failure {
	f.Message = msg
	return f
}

func (f *Failure) Send(c *gin.Context) error {
	if f.StatusCode == 0 {
		f.StatusCode = 500
	}
	c.JSON(f.StatusCode, f)
	return nil
}
