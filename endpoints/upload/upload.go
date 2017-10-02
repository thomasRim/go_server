package upload

import (
	"net/http"
	"crypto/md5"
	"io"
	"strconv"
	"fmt"
	"html/template"
	"time"
	"github.com/gin-gonic/gin"
	"log"
)

const (
	UploadURI = "/upload"
	uploadDirPath = "./uploadedfiles/"
)
func Upload(context *gin.Context) {
	if context.Request.Method == "GET" {
		crutime := time.Now().Unix()

		h := md5.New()
		io.WriteString(h, strconv.FormatInt(crutime, 10))
		token := fmt.Sprintf("%x", h.Sum(nil))

		t, _ := template.ParseFiles("./addons/uploadForm.html")
		context.Header("Content-Type","text/html")
		t.Execute(context.Writer, token)
	} else {
		file, _ := context.FormFile("file")
		log.Println(file.Filename)

		// Upload the file to specific dst.
		context.SaveUploadedFile(file, uploadDirPath+file.Filename)

		context.String(http.StatusOK, fmt.Sprintf("'%s' uploaded!", file.Filename))
	}
}
