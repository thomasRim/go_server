package user

import (
	"fmt"
	"net/http"
	"encoding/json"
	"errors"
	"../../database"
	"../../utils"
	"../../constants"
	"github.com/gin-gonic/gin"
	"strings"
)

const (
	UsersURI    = "/users"
	pathParam   = "userId"
	UserInfoURI = UsersURI + "/:" + pathParam

	tableName        = "person"
	createTableQuery = "CREATE TABLE IF NOT EXISTS %v (" +
		"id INTEGER PRIMARY KEY AUTOINCREMENT," +
		" userId VARCHAR(20) NOT NULL DEFAULT ''," +
		" first_name VARCHAR(40) NOT NULL DEFAULT ''," +
		" last_name VARCHAR(40) NOT NULL DEFAULT ''," +
		" created_time VARCHAR(40)," +
		" user_name VARCHAR(40) DEFAULT ''," +
		" email VARCHAR(40) DEFAULT ''," +
		" hash VARCHAR(40) " +
		");"
	databaseName = constants.SQLiteDbFile
)

/*
	Types
 */
type User struct {
	UserId       string `json:"user_id"`
	First_name   string `json:"first_name"`
	Last_name    string `json:"last_name"`
	Created_time string `json:"created_time,omitempty"`
	UserName     string `json:"user_name"`
	Email        string `json:"email,omitempty"`
}

type UserSlice struct {
	Users []User `json:"users"`
}

/*
	Funcs
 */

// Public

func Users(context *gin.Context) {
	prepareTable()

	switch context.Request.Method {
	case http.MethodGet:
		if len(context.Params) > 0 {
			param := context.Param(pathParam)
			user, err := fetchUserWithID(param)
			if checkIfErr(err) {
				context.String(http.StatusInternalServerError, "Fetch user error. "+err.Error())
				break
			}
			context.JSON(http.StatusOK, user)
		} else {
			context.JSON(http.StatusOK, GetAllUsers())
		}
	case http.MethodPost:
		user, err := insertUser(context.Request)
		if checkIfErr(err) {
			context.String(http.StatusInternalServerError, "Add user error. "+err.Error())
			break
		}
		context.JSON(http.StatusOK, user)
	case http.MethodPatch:
		if len(context.Params) > 0 {
			param := context.Param(pathParam)
			user, err := updateUser(param, context.Request)
			if checkIfErr(err) {
				context.String(http.StatusInternalServerError, "Update user error. "+err.Error())
				break
			}
			context.JSON(http.StatusOK, user)
		}
		context.String(http.StatusBadRequest, "Unable to find required user ID as path parameter")

	case http.MethodDelete:
		if len(context.Params) > 0 {
			param := context.Param(pathParam)
			user, err := removeUserWithId(param)
			if checkIfErr(err) {
				context.String(http.StatusInternalServerError, "Delete user error "+err.Error())
				break
			}
			context.JSON(http.StatusOK, user)

		}
		context.String(http.StatusBadRequest, "Unable to find required user ID as path parameter")
	default:
		context.String(http.StatusNotImplemented, "Not implemented")
	}

}

func GetAllUsers() UserSlice {
	prepareTable()

	db, err := database.DB(databaseName)
	checkIfErr(err)

	res, err := db.Query("SELECT " +
		"user_id," +
		"first_name," +
		"last_name," +
		"user_name," +
		"email," +
		"created_time" +
		" FROM " + tableName +
		" LIMIT 100")
	checkIfErr(err)

	user := User{}

	var slice UserSlice
	slice.Users = []User{}

	for res.Next() {
		err = res.Scan(&user.UserId,
			&user.First_name,
			&user.Last_name,
			&user.UserName,
			&user.Email,
			&user.Created_time)
		checkIfErr(err)

		slice.Users = append(slice.Users, user)
	}
	res.Close()

	return slice
}

// Private

func insertUser(r *http.Request) (User, error) {

	// getting obj from body json
	user, err := parseJsonFromBody(r)
	if checkIfErr(err) {
		return *user, err
	}

	_, err = findUserExisting(user)
	if err != nil {
		return *user, err
	}

	//check if fields not empty
	user.UserId, err = utils.NewUUID()
	user.Created_time = utils.TimeInUTC()

	db, err := database.DB(databaseName)
	if checkIfErr(err) {
		return *user, err
	}

	stmt, err := db.Prepare("INSERT INTO " + tableName + "(" +
		"user_id," +
		"first_name," +
		"last_name," +
		"email," +
		"user_name," +
		"created_time)" +
		" values(?,?,?,?,?,?)")
	if checkIfErr(err) {
		return *user, err
	}

	_, err = stmt.Exec(
		user.UserId,
		user.First_name,
		user.Last_name,
		user.Email,
		user.UserName,
		user.Created_time)
	if checkIfErr(err) {
		return *user, err
	}

	return *user, err

}

func updateUser(userId string, req *http.Request) (User, error) {
	user, err := fetchUserWithID(userId)
	if checkIfErr(err) {
		return user, err
	}

	obj, err := parseJsonFromBody(req)
	if checkIfErr(err) {
		return user, err
	}
	parsedUser := *obj

	db, err := database.DB(databaseName)
	if checkIfErr(err) {
		return user, err
	}

	reqStr := "UPDATE "+tableName+" SET "

	arr := []string{}
	if len(parsedUser.UserName)>0 {
		str := "user_name="+parsedUser.UserName
		_=append(arr, str)
	} else if len(parsedUser.Last_name)>0 {
		str:="last_name="+parsedUser.Last_name
		_=append(arr,str)
	} else if len(parsedUser.First_name)>0 {
		str:="first_name="+parsedUser.First_name
		_=append(arr,str)
	}
	if len(arr)>0 {
		reqStr = reqStr+strings.Join(arr,",")+"WHERE user_id="+userId

		res,err := db.Query(reqStr)
		fmt.Println(res)

		if checkIfErr(err) {
			return user,err
		}
	}

	user, err = fetchUserWithID(userId)
	if checkIfErr(err) {
		return user, err
	}

	return user,nil
}

func removeUserWithId(userId string) (User, error) {
	user, err := fetchUserWithID(userId)
	if checkIfErr(err) {
		return user, err
	}

	db, err := database.DB(databaseName)
	if checkIfErr(err) {
		return user, err
	}

	_, err = db.Exec("DELETE FROM "+tableName+" WHERE user_id=?", userId)
	if checkIfErr(err) {
		return user, err
	}
	return user, nil
}

/*
	Helpers
 */

func fetchUserWithID(id string) (User, error) {
	user := User{}
	db, err := database.DB(databaseName)
	if checkIfErr(err) {
		return user, err
	}

	chck, err := db.Query("SELECT "+
		"user_id,"+
		"first_name,"+
		"last_name,"+
		"email,"+
		"user_name"+
		" FROM "+ tableName+
		" WHERE user_id=? ", id)
	if checkIfErr(err) {
		return user, err
	}

	for chck.Next() {
		chck.Scan(&user.UserId, &user.First_name, &user.Last_name, &user.Email, &user.UserName)
	}

	return user, nil
}

func findUserExisting(user *User) (bool, error) {
	db, err := database.DB(databaseName)
	if checkIfErr(err) {
		return false, err
	}

	chck, err := db.Query("SELECT "+
		"id,"+
		"user_id,"+
		"first_name"+
		" FROM "+ tableName+
		" WHERE first_name=? AND last_name=?", user.First_name, user.Last_name)
	if checkIfErr(err) {
		return false, err
	}

	var id uint32
	var uid string
	var name string
	for chck.Next() {
		chck.Scan(&id, &uid, &name)
	}
	if len(name) != 0 {
		return true, errors.New("{\"error\":\"such user is existing\"}")
	} else {
		return false, nil
	}
}

func parseJsonFromBody(r *http.Request) (*User, error) {
	decoder := json.NewDecoder(r.Body)
	var t *User
	err := decoder.Decode(&t)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	return t, nil
}
func prepareTable() {
	database.Prepare(
		tableName,
		fmt.Sprintf(createTableQuery, tableName),
		databaseName)
}

func checkIfErr(err error) bool {
	if err != nil {
		fmt.Println(err)
		return true
	}
	return false
}
