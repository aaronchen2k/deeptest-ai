package middleware

import (
	multi_iris "github.com/deeptest-com/deeptest-next/internal/pkg/core/auth/iris"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
)

func JwtHandler() iris.Handler {
	verifier := multi_iris.NewVerifier()

	// extract token only from Authorization: Bearer $token
	verifier.Extractors = []multi_iris.TokenExtractor{multi_iris.FromHeader}

	verifier.ErrorHandler = func(ctx *context.Context, err error) {
		//ctx.StatusCode(iris.StatusUnauthorized)
		ctx.JSON(_domain.Response{
			Code: _domain.AuthErr.Code,
			Msg:  ctx.Path(),
		})
		// ctx.StopWithError(http.StatusUnauthorized, err)
	} // extract token only from Authorization: Bearer $token
	return verifier.Verify()
}
