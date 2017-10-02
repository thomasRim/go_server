package login

import (
	"net/http"
	"html/template"
	"fmt"
	"encoding/base64"
	"github.com/gin-gonic/gin"
	"time"
	"crypto/md5"
	"io"
	"strconv"
)

const LoginURI = "/login"

func Login(context *gin.Context)  {
	switch context.Request.Method {
	case "GET":
		crutime := time.Now().Unix()

		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("./addons/loginForm.html")
		context.Header("Content-Type","text/html")
		t.Execute(context.Writer, token)
	case "POST":
		// parse
		username := context.PostForm("username")
		password := context.PostForm("password")

		fmt.Println("username:", username)
		fmt.Println("password:", password)

		// search DB , generate and return token
		str := username+"/"+password
		l := base64.StdEncoding.EncodeToString([]byte(str))
		context.String(200,"token: "+l)
	default:
		context.String(http.StatusNotImplemented, "Not implemented")
	}
}
