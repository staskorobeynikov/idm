package main

import (
	"fmt"
	"github.com/jmoiron/sqlx"
	"idm/inner/common"
	"idm/inner/database"
	"idm/inner/employee"
	"idm/inner/info"
	"idm/inner/role"
	"idm/inner/validator"
	"idm/inner/web"
)

func main() {
	cfg := common.GetConfig(".env")
	db := database.ConnectDbWithCfg(cfg)
	defer func() {
		if err := db.Close(); err != nil {
			fmt.Printf("error closing db: %v", err)
		}
	}()
	var server = build(cfg, db)
	var err = server.App.Listen(":8080")
	if err != nil {
		panic(fmt.Sprintf("http server error: %s", err))
	}
}

func build(cfg common.Config, database *sqlx.DB) *web.Server {
	var server = web.NewServer()
	var employeeRepo = employee.NewRepository(database)
	var roleRepo = role.NewRepository(database)
	var vld = validator.New()
	var employeeService = employee.NewService(employeeRepo, vld)
	var employeeController = employee.NewController(server, employeeService)
	employeeController.RegisterRoutes()
	var roleService = role.NewService(roleRepo, vld)
	var roleController = role.NewController(server, &roleService)
	roleController.RegisterRoutes()
	var infoController = info.NewController(server, cfg, database)
	infoController.RegisterRoutes()
	return server
}
