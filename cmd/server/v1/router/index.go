package router

import (
	"github.com/kataras/iris/v12"
)

type IndexModule struct {
	AccountModule *AccountModule `inject:""`
	OptLogModule  *OptLogModule  `inject:""`
	PermModule    *PermModule    `inject:""`
	RoleModule    *RoleModule    `inject:""`
	UserModule    *UserModule    `inject:""`
	FileModule    *FileModule    `inject:""`

	DataModule    *DataModule    `inject:""`
	ProjectModule *ProjectModule `inject:""`
	CaseModule    *CaseModule    `inject:""`
	SetModule     *SetModule     `inject:""`
	PlanModule    *PlanModule    `inject:""`
	ChatbotModule *ChatbotModule `inject:""`
}

func NewIndexModule() *IndexModule {
	return &IndexModule{}
}

func (m *IndexModule) ApiParty() func(rbac iris.Party) {
	return func(rbac iris.Party) {
		rbac.PartyFunc("/accounts", m.AccountModule.Party())
		rbac.PartyFunc("/optlogs", m.OptLogModule.Party())
		rbac.PartyFunc("/roles", m.RoleModule.Party())
		rbac.PartyFunc("/perms", m.PermModule.Party())
		rbac.PartyFunc("/users", m.UserModule.Party())

		rbac.PartyFunc("/init", m.DataModule.Party())
		rbac.PartyFunc("/cases", m.CaseModule.Party())
		rbac.PartyFunc("/sets", m.SetModule.Party())
		rbac.PartyFunc("/plans", m.PlanModule.Party())
		rbac.PartyFunc("/projects", m.ProjectModule.Party())
		rbac.PartyFunc("/file", m.FileModule.Party())

		rbac.PartyFunc("/chatbot", m.ChatbotModule.Party())
	}
}
