package middleware

import (
	multi_iris "github.com/deeptest-com/deeptest-next/internal/pkg/core/auth/iris"
	"net/http"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

/**
 * 验证 multi
 * @method MultiHandler
 */
func MultiHandler() iris.Handler {
	verifier := multi_iris.NewVerifier()
	verifier.Extractors = []multi_iris.TokenExtractor{multi_iris.FromHeader} // extract token only from Authorization: Bearer $token
	verifier.ErrorHandler = func(ctx *context.Context, err error) {
		ctx.StopWithError(http.StatusUnauthorized, err)
	}
	return verifier.Verify()
}
