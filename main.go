package main

//noinspection ALL
import (
	"fmt"
	"os"
	"github.com/gin-gonic/gin"

	"./database"
	"./constants"
	// endpoints
	"./endpoints/login"
	"./endpoints/upload"
	"./endpoints/home"
	"./endpoints/user"
	"./endpoints/receipt"
)

func root(context *gin.Context) {
	str := constants.RootURI+"\t- list (this)\n"+
		home.HomeURI+"\t- home items\n"+
		login.LoginURI+"\t- login\n"+
		upload.UploadURI+"\t- upload files\n"+
		user.UsersURI+"\t- user service\n"+
		receipt.ReceiptURI+"\t- receipts service\n"
	context.String(200,str)
}

// - Main
func main() {

	_,_ = os.OpenFile(constants.SQLiteDbFile,os.O_CREATE,0666)
	database.DB(constants.SQLiteDbFile)

	router := gin.Default()

	// Root and Home
	router.GET(constants.RootURI, root)
	router.GET(home.HomeURI, home.Home)

	// Upload
	router.GET(upload.UploadURI, upload.Upload)
	router.POST(upload.UploadURI, upload.Upload)

	// Login
	router.GET(login.LoginURI, login.Login)
	router.POST(login.LoginURI, login.Login)

	// User
	router.GET(user.UsersURI, user.Users)
	router.POST(user.UsersURI, user.Users)

	router.GET(user.UserInfoURI, user.Users)
	router.DELETE(user.UserInfoURI, user.Users)

	//Receipt
	router.GET(receipt.ReceiptURI, receipt.Receipts)
	router.POST(receipt.ReceiptURI, receipt.Receipts)

	router.GET(receipt.ReceiptInfoURI, receipt.Receipts)
	router.DELETE(receipt.ReceiptInfoURI, receipt.Receipts)

	//Exec
	router.Run(constants.Port)

	fmt.Println("Start listening on port "+constants.Port)

}

