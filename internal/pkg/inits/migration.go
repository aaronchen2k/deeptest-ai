package inits

import (
	"github.com/deeptest-com/deeptest-next/cmd/server/v1/router"
	"github.com/deeptest-com/deeptest-next/internal/pkg/core/auth"
	"github.com/deeptest-com/deeptest-next/internal/pkg/core/migration"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/database"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/operation"
	"github.com/deeptest-com/deeptest-next/internal/pkg/serve/web/web_iris"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/model"
	"github.com/deeptest-com/deeptest-next/internal/server/moudules/source"
	_logs "github.com/deeptest-com/deeptest-next/pkg/libs/log"
	"go.uber.org/zap"
)

// 加载模块
var PartyFunc = func(wi *web_iris.WebServer) {
	// 初始化驱动
	err := auth.InitDriver(&auth.Config{DriverType: "jwt", HmacSecret: nil})
	if err != nil {
		_logs.Zap.Panic("err")
	}

	indexModule := router.NewIndexModule()

	wi.AddModule(web_iris.Party{
		Perfix:    "/api/v1",
		PartyFunc: indexModule.ApiParty(),
	})
}

// 填充数据
var SeedFunc = func(wi *web_iris.WebServer, mc *migration.MigrationCmd) {
	err := database.GetInstance().AutoMigrate(model.Models...)
	if err != nil {
		_logs.Errorf("初始化数据表错误", zap.String("错误:", err.Error()))
		return
	}

	mc.AddMigration(
		source.GetPermMigration(),
		source.GetRoleMigration(),
		source.GetUserMigration(),
		operation.GetMigration())

	routes, _ := wi.GetSources()

	// 权鉴模块全部为管理员权限
	authorityTypes := map[string]int{}
	for _, route := range routes {
		authorityTypes[route["path"]] = auth.AdminAuthority
	}

	// notice : 注意模块顺序
	mc.AddSeed(source.NewPermSource(routes), source.RoleSrc, source.UserSrc)
}
