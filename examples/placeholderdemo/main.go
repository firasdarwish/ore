package main

import (
	"context"
	"fmt"

	"github.com/firasdarwish/ore"
)

type UserRole struct {
}
type SomeService struct {
	someConfig string
}

func main() {
	//register SomeService which depends on "someConfig"
	ore.RegisterFunc[*SomeService](ore.Scoped, func(ctx context.Context) (*SomeService, context.Context) {
		someConfig, ctx := ore.Get[string](ctx, "someConfig")
		return &SomeService{someConfig}, ctx
	})

	//someConfig is unknow at registration time
	//the value of "someConfig" depends on the future user's request
	ore.RegisterPlaceholder[string]("someConfig")

	//Seal registration, no further registration is allowed
	ore.Seal()
	ore.Validate()

	//a request arrive
	ctx := context.Background()
	//suppose that the request is sent by "admin"
	ctx = context.WithValue(ctx, "role", "admin")

	//inject a different config depends on the request,
	userRole := ctx.Value("role").(string)
	if userRole == "admin" {
		ctx = ore.ProvideScopedValue(ctx, "Admin config", "someConfig")
	} else if userRole == "supervisor" {
		ctx = ore.ProvideScopedValue(ctx, "Supervisor config", "someConfig")
	} else if userRole == "user" {
		ctx = ore.ProvideScopedValue(ctx, "Public user config", "someConfig")
	}

	service, _ := ore.Get[*SomeService](ctx)
	fmt.Println(service.someConfig) //"Admin config"
}
