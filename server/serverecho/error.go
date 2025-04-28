package serverecho

import (
	"encoding/json"
	"errors"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/worldline-go/rest"
)

func HTTPErrorHandler(err error, c echo.Context) {
	if c.Response().Committed {
		return
	}

	errStr := ""
	code := http.StatusInternalServerError
	msg := http.StatusText(http.StatusInternalServerError)
	var he *echo.HTTPError
	if errors.As(err, &he) {
		code = he.Code
		switch m := he.Message.(type) {
		case string:
			msg = m
		case json.Marshaler:
			msgByte, err := json.Marshal(m)
			if err != nil {
				c.Logger().Errorf("failed to marshal error message: %v", err)
			} else {
				msg = string(msgByte)
			}
		case error:
			errStr = m.Error()
		}

		if he.Internal != nil {
			errStr = he.Internal.Error()
		}
	}

	// Send response
	if c.Request().Method == http.MethodHead { // Issue #608
		err = c.NoContent(he.Code)
	} else {
		err = c.JSON(code, rest.ResponseMessage{
			Message: &rest.Message{
				Text: msg,
				Err:  errStr,
			},
		})
	}
	if err != nil {
		c.Logger().Error(err.Error())
	}
}
