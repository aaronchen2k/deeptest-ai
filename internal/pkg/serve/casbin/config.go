package casbin

import (
	"fmt"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	_file "github.com/deeptest-com/deeptest-next/pkg/libs/file"
	"path/filepath"
)

// Remove del config file
func Remove() error {
	casbinPath := getCasbinPath()
	if _file.IsExist(casbinPath) && _file.IsFile(casbinPath) {
		return _file.Remove(casbinPath)
	}
	return nil
}

func getCasbinPath() string {
	return filepath.Join(consts.ExecDir, consts.ConfDir, consts.CasbinFileName)
}

// init initialize config file
// - initialize casbin's config file as rbac_model.conf name
func init() {
	casbinPath := getCasbinPath()
	fmt.Printf("casbin rbac_model.conf's path: %s\n\n", casbinPath)
	if !_file.IsExist(casbinPath) { // casbin rbac_model.conf file
		var rbacModelConf = []byte(`[request_definition]
	r = sub, obj, act

	[policy_definition]
	p = sub, obj, act

	[role_definition]
	g = _, _

	[policy_effect]
	e = some(where (p.eft == allow))

	[matchers]
	m = g(r.sub, p.sub) && keyMatch2(r.obj, p.obj) && (r.act == p.act || p.act == "*")`)
		_, err := _file.WriteBytes(casbinPath, rbacModelConf)
		if err != nil {
			panic(fmt.Errorf("initialize casbin rbac_model.conf file return error: %w ", err))
		}
	}
}
