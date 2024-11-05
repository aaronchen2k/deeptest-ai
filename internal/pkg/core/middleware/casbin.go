package middleware

import (
	"errors"
	"fmt"
	multi_iris "github.com/deeptest-com/deeptest-next/internal/pkg/core/auth/iris"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/casbin"
	_domain "github.com/deeptest-com/deeptest-next/pkg/domain"
	_logUtils "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"net/http"
	"strconv"

	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"go.uber.org/zap"
)

// Casbin Casbin 权鉴中间件
func Casbin() iris.Handler {
	return func(ctx *context.Context) {
		check, err := Check(ctx.Request(), strconv.FormatUint(uint64(multi_iris.GetUserId(ctx)), 10))
		if err != nil || !check {
			ctx.JSON(_domain.Response{Code: _domain.AuthActionErr.Code, Msg: err.Error()})

			ctx.StopExecution()
			return
		}

		ctx.Next()
	}
}

// Check checks the username, request's method and path and
// returns true if permission grandted otherwise false.
func Check(r *http.Request, userId string) (bool, error) {
	method := r.Method
	path := r.URL.Path
	ok, err := casbin.Instance().Enforce(userId, path, method)
	if err != nil {
		_logUtils.Zap.Error(fmt.Sprintf("验证权限报错：%s-%s-%s", userId, path, method), zap.String("casbinServer.Instance().Enforce()", err.Error()))
		return false, err
	}

	_logUtils.Debug(fmt.Sprintf("权限：%s-%s-%s", userId, path, method))

	if !ok {
		return ok, errors.New("你未拥有当前操作权限，请联系管理员")
	}
	return ok, nil
}
