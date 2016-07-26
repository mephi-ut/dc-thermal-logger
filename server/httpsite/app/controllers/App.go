package controllers

import (
	//"fmt"
	"devel.mephi.ru/dyokunev/dc-thermal-logger/server/httpsite/app"
	"github.com/revel/revel"
)

type App struct {
	*revel.Controller
}

func (c App) Logout() revel.Result {
	app.CasClient.RedirectToLogout(c.Response.Out, c.Request.Request)
	return nil
}
