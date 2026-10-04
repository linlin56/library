package controllers

import (
	"library/i18n"

	"github.com/gin-gonic/gin"
)

// tr translates a message key into the language requested by the client (Accept-Language).
func tr(c *gin.Context, key string, args ...any) string {
	return i18n.T(i18n.FromHeader(c.GetHeader("Accept-Language")), key, args...)
}

// trErr translates a business error, or returns the generic fallbackKey message.
func trErr(c *gin.Context, err error, fallbackKey string) string {
	return i18n.ErrorMessage(i18n.FromHeader(c.GetHeader("Accept-Language")), err, fallbackKey)
}
