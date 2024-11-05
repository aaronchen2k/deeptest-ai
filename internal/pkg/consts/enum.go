package consts

type RoleType string

const (
	Admin              RoleType = "admin"
	User               RoleType = "user"
	Tester             RoleType = "tester"
	Developer          RoleType = "developer"
	ProductManager     RoleType = "product_manager"
	IntegrationAdmin   RoleType = "api-admin"
	IntegrationGeneral RoleType = "general"
)

func (e RoleType) String() string {
	return string(e)
}

type CategoryDiscriminator string

const (
	EndpointCategory CategoryDiscriminator = "endpoint"
	ScenarioCategory CategoryDiscriminator = "scenario"
	PlanCategory     CategoryDiscriminator = "plan"
	SchemaCategory   CategoryDiscriminator = "schema"
)

func (e CategoryDiscriminator) String() string {
	return string(e)
}
