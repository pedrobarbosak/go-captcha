# go-captcha
Golang Captcha

## Examples

### Basic
```go
package main

import (
	"log"

	"github.com/pedrobarbosak/go-captcha"
)

func main() {
	recaptcha := captcha.Recaptcha("<secret-key>", 0.5, "<optional-override-url>")

	err := recaptcha.Verify("response")
	if err != nil {
		log.Panic("failed:", err)
	}
	
	log.Println("success")
}
```


### As a Middleware
```go
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"

	"github.com/pedrobarbosak/go-captcha"
)

type Request struct {
	CaptchaToken string `form:"captchaResponse" json:"captchaResponse" validate:"required"`
}

func Middleware(service captcha.Service) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var request Request
		if err := ctx.ShouldBindBodyWith(&request, binding.JSON); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid captcha"})
			ctx.Abort()
			return
		}

		if err := service.Verify(request.CaptchaToken); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{"error": "invalid captcha"})
			ctx.Abort()
			return
		}

		ctx.Next()
	}
}
```

Also available on: https://github.com/pedrobarbosak/go-captcha-gin
