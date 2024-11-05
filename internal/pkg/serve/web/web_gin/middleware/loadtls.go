package middleware

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/unrolled/secure"
)

// LoadTls
func LoadTls() gin.HandlerFunc {
	return func(c *gin.Context) {
		middleware := secure.New(secure.Options{
			SSLRedirect: true,
			SSLHost:     "0.0.0.0:443",
		})
		err := middleware.Process(c.Writer, c.Request)
		if err != nil {
			fmt.Println(err)
			return
		}
		c.Next()
	}
}
