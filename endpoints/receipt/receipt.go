package receipt

import (
	"net/http"
	"fmt"
	"../../database"
	"../../constants"
	"../../utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
)

const (
	ReceiptURI = "/receipts"
	pathParam = "receiptId"
	ReceiptInfoURI = ReceiptURI+"/:"+ pathParam
	tableName   = "receipt"
	createQuery = "CREATE TABLE IF NOT EXISTS %v (" +
		"id INTEGER PRIMARY KEY AUTOINCREMENT," +
		" receiptId VARCHAR(40)," +
		" date VARCHAR(40) DEFAULT ''," +
		" place VARCHAR(40) DEFAULT ''," +
		" total INTEGER DEFAULT 0," +
		" items VARCHAR (500) DEFAULT '[]'," +
		" created_time VARCHAR(40)," +
		" updated_time VARCHAR(40) DEFAULT '');"
	databaseName = constants.SQLiteDbFile
)

/*
	Types
 */

type Receipt struct {
	ReceiptId string		`json:"receiptId"`
	Date string				`json:"date,omitempty"`
	Place string			`json:"place,omitempty"`
	Total int				`json:"total"`
	Items string			`json:"items,omitempty"`

	CreatedTime string		`json:"created_time"`
	UpdatedTime string		`json:"updated_time,omitempty"`
}

type ReceipstSlice struct {
	Receipts []Receipt `json:"receipts"`
}

/*
	Funcs
 */

// Public

func Receipts(context *gin.Context) {
	prepareTable()

	switch context.Request.Method {
	case http.MethodGet:
		if len(context.Params) > 0 {
			param := context.Param(pathParam)
			receipt,err := fetchReceiptById(param)
			if checkIfErr(err) {
				context.String(http.StatusInternalServerError, "Fetch receipt error. "+err.Error())
				break
			}
			context.JSON(http.StatusOK, receipt)
			// fetch receipt by id
		} else {
			context.JSON(http.StatusOK,allReceipts())
		}
	case http.MethodPost:
		context.JSON(http.StatusOK,addReceipt(context.Request))
	//case http.MethodPatch:
	case http.MethodDelete:
		if len(context.Params) > 0 {
			param := context.Param(pathParam)
			receipt,err := removeReceiptWithId(param)
			if checkIfErr(err) {
				context.String(http.StatusInternalServerError, "Delete receipt error "+err.Error())
				break
			}
			context.JSON(http.StatusOK,receipt)

		}
		context.String(http.StatusBadRequest,"Unable to find required receipt ID as path parameter")

	default:
		context.String(http.StatusNotImplemented,"Not implemented")
	}
}

// Private

func allReceipts() ReceipstSlice  {

	var slice ReceipstSlice
	slice.Receipts = []Receipt{}

	db,err := database.DB(databaseName)
	if checkIfErr(err) {
		fmt.Println("Load DB error "+err.Error())
		return slice
	}

	res, err := db.Query("SELECT " +
		"receiptId," +
		"date," +
		"place," +
		"total," +
		"items," +
		"created_time," +
		"updated_time" +
		" FROM "+ tableName +
		" LIMIT 100")
	if checkIfErr(err) {
		fmt.Println("Query error "+err.Error())
		return slice
	}

	dbObj := Receipt{}

	for res.Next() {
		res.Scan(
			&dbObj.ReceiptId,
			&dbObj.Date,
			&dbObj.Place,
			&dbObj.Total,
			&dbObj.Items,
			&dbObj.CreatedTime,
			&dbObj.UpdatedTime)

		slice.Receipts = append(slice.Receipts, dbObj)
	}
	defer res.Close()

	return slice
}

func fetchReceiptById(receiptId string) (Receipt, error)  {
	receipt := Receipt{}
	db, err := database.DB(databaseName)
	if checkIfErr(err) {
		return receipt, err
	}

	chck, err := db.Query("SELECT " +
		"receiptId," +
		"date," +
		"place," +
		"total," +
		"items," +
		"created_time," +
		"updated_time" +
		" FROM "+ tableName +
		" WHERE receiptId=? ", receiptId)
	if checkIfErr(err) {
		return receipt, err
	}

	for chck.Next() {
		chck.Scan(
			&receipt.ReceiptId,
			&receipt.Date,
			&receipt.Place,
			&receipt.Total,
			&receipt.Items,
			&receipt.CreatedTime,
			&receipt.UpdatedTime)
	}

	return receipt, nil
}

func addReceipt(r *http.Request) []byte {
	// get instance
	object,err := parseJsonFromBody(r)
	if checkIfErr(err) {
		return []byte(err.Error())
	}

	// fill required fields
	object.ReceiptId, err = utils.NewUUID()
	object.ReceiptId = "Receipt_"+object.ReceiptId
	object.CreatedTime = utils.TimeInUTC()

	db, err := database.DB(databaseName)
	if checkIfErr(err) {
		return []byte(err.Error())
	}

	stmt, err := db.Prepare("INSERT INTO receipt(" +
		"receiptId," +
		"date," +
		"place," +
		"total," +
		"items," +
		"created_time)" +
		" VALUES(?,?,?,?,?,?)")
	if checkIfErr(err) {
		return []byte(err.Error())
	}

	_, err = stmt.Exec(
		object.ReceiptId,
		object.Date,
		object.Place,
		object.Total,
		object.Items,
		object.CreatedTime)
	if checkIfErr(err) {
		return []byte(err.Error())
	}

	// return created instance
	str, err := json.Marshal(object)
	if checkIfErr(err) {
		return []byte(err.Error())
	}
	return str
}

func removeReceiptWithId(receiptId string) (Receipt, error) {
	receipt, err := fetchReceiptById(receiptId)
	if checkIfErr(err) {
		return receipt,err
	}

	db,err := database.DB(databaseName)
	if checkIfErr(err) {
		return receipt,err
	}

	_,err = db.Exec("DELETE FROM "+tableName+" WHERE receiptId=?",receiptId)
	if checkIfErr(err) {
		return receipt,err
	}
	return receipt,nil
}

// Helpers

func parseJsonFromBody(r *http.Request) (*Receipt, error) {
	decoder := json.NewDecoder(r.Body)
	var t *Receipt
	err := decoder.Decode(&t)
	if err != nil {
		return nil,err
	}
	defer r.Body.Close()

	return t,nil
}

func prepareTable() {
	database.Prepare(
		tableName,
		fmt.Sprintf(createQuery,tableName),
		databaseName)
}

func checkIfErr(err error) bool {
	if err != nil {
		fmt.Println("Error "+err.Error())
		return true
	}
	return false
}
