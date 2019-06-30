package controllers

var personController PersonController
var healthController HealthController
var userController UserController

//registers routes
func Startup() {
	personController.registerRoutes()
	healthController.registerRoutes()
	userController.registerRoutes()
}
