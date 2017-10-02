package home

import (
	"net/http"
	"../user"
	"github.com/gin-gonic/gin"
)

const HomeURI = "/home"

func Home(context *gin.Context) {
	if context.Request.Method == "GET" {
		str2 := user.GetAllUsers()
		context.JSON(200, str2)
	} else {
		context.String(http.StatusNotImplemented, "Not implemented")
	}
}

