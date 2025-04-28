package server

import (
	"strings"

	"github.com/labstack/echo/v4"
)

const (
	// UserKey is the key for the user in the context and value is string.
	UserKey = "X-User"
	// BrowserKey is the key for the browser in the context and value is bool.
	BrowserKey = "X-Browser"
)

func IsBrowserWithAgent(userAgent string) bool {
	userAgent = strings.ToLower(userAgent)
	if userAgent != "" {
		// Check for common browser user-agent strings
		return strings.Contains(userAgent, "mozilla") ||
			strings.Contains(userAgent, "chrome") ||
			strings.Contains(userAgent, "safari") ||
			strings.Contains(userAgent, "opera") ||
			strings.Contains(userAgent, "msie") ||
			strings.Contains(userAgent, "edge")
	}

	return false
}

// MiddlewareUserInfo adds the user and browser to echo's store.
//
//	user, _ := c.Get(UserKey).(string)
//	isBrowser, _ := c.Get(BrowserKey).(bool)
func MiddlewareUserInfo(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		user := c.Request().Header.Get(UserKey)
		c.Set(UserKey, user)

		browser := IsBrowserWithAgent(c.Request().Header.Get("User-Agent"))
		c.Set(BrowserKey, browser)

		return next(c)
	}
}

func IsBrowser(c echo.Context) bool {
	isBrowser, ok := c.Get(BrowserKey).(bool)
	if !ok {
		return false
	}

	return isBrowser
}

func GetUser(c echo.Context) string {
	user, ok := c.Get(UserKey).(string)
	if !ok {
		return ""
	}

	return user
}

func GetUserOrDefault(c echo.Context, defaultUser string) string {
	user, ok := c.Get(UserKey).(string)
	if !ok {
		return defaultUser
	}

	return user
}
