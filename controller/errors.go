package controller

import (
	"encoding/json"
	"fmt"

	"github.com/gin-gonic/gin"
)

//ErrMsg is returned as response in API when an error occurred.
type ErrMsg struct {
	Status int   `json:"status"`
	Error  error `json:"-"`
}

//MarshalJSON is a function of Marshaler interface.
func (err ErrMsg) MarshalJSON() ([]byte, error) {
	type message ErrMsg
	return json.Marshal(&struct {
		Error string `json:"error"`
		message
	}{
		Error:   err.Error.Error(),
		message: (message)(err),
	})
}

func (err ErrMsg) String() string {
	return fmt.Sprintf("Status %d: %s", err.Status, err.Error.Error())
}

//JSON uses context to return error in response of API.
func (err ErrMsg) JSON(c *gin.Context) {
	c.JSON(err.Status, err)
}

//NewError returns and logs an error ocurred.
func NewError(status int, err error) ErrMsg {
	//Log error
	return ErrMsg{
		Status: status,
		Error:  err,
	}
}

//JSONError creates a new error and returns response of API at the same time.
func JSONError(status int, err error, c *gin.Context) {
	er := NewError(status, err)
	er.JSON(c)
}
