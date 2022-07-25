package API

import (
	"database/sql"
	"db_p/pickbuy"
	"db_p/profile"
	"db_p/signUP_IN"
	"db_p/structs"
	"github.com/gin-gonic/gin"
	"github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"strconv"
)

func checkErr(c *gin.Context, err error, x gin.H) {
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
	} else {
		c.IndentedJSON(http.StatusOK, x)
	}
}
func setupDB() *sql.DB {
	cfg := mysql.Config{
		User:                 "root",
		Passwd:               "ayda",
		Net:                  "tcp",
		Addr:                 "127.0.0.1:3306",
		DBName:               "mySql",
		AllowNativePasswords: true,
	}

	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Print(err.Error())
	}

	return db
}

type user struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func getID(c *gin.Context) int {
	Id := c.Param("id")
	userId, err := strconv.Atoi(Id)
	checkErr(c, err, gin.H{"data": userId})
	return userId
}
func CreatUser(c *gin.Context) {
	Email := c.Param("email")
	Password := c.Param("password")
	userId, err := signUP_IN.SignUp(Email, Password, setupDB())
	if err != nil {
		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": userId, "email": Email, "password": Password, "message": "you signed up successfully!"})
}
func SignInUser(c *gin.Context) {
	Email := c.Param("email")
	Password := c.Param("password")
	userId, err := signUP_IN.SignIn(Email, Password, setupDB())
	if (err != nil) || (userId == -1) {
		c.String(http.StatusInternalServerError, err.Error())

		c.AbortWithError(http.StatusInternalServerError, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"data": userId, "email": Email, "password": Password, "message": "you signed in successfully!"})
}

func News(c *gin.Context) {
	Id := getID(c)
	NewProducts := profile.ShowNews(setupDB(), Id)
	c.JSON(http.StatusOK, gin.H{"newProducts": NewProducts, "id": Id, "message": "products that were unavailable for you are available now !"})
}
func ModifyMobile(c *gin.Context) {
	userId := getID(c)
	temp := c.Param("mobile")
	err := profile.Modify_mobile(userId, temp, setupDB())
	checkErr(c, err, gin.H{"id": userId, "new_phone_num": temp, "message": "you changed your phone number!"})

}
func ModifyEmail(c *gin.Context) {
	Id := getID(c)
	newEmail := c.Param("email")
	err := profile.Modify_email(Id, newEmail, setupDB())
	checkErr(c, err, gin.H{"id": Id, "new_email": newEmail, "message": "you changed your Email!"})
}
func ModifyPassword(c *gin.Context) {
	Id := getID(c)
	newPassword := c.Param("password")
	err := profile.Modify_password(Id, newPassword, setupDB())
	checkErr(c, err, gin.H{"id": Id, "newpassword": newPassword, "message": "you changed your password!"})

}
func ModifyAdd2(c *gin.Context) {
	Id := getID(c)
	add2 := c.Param("add2")
	err := profile.Modify_address2(Id, add2, setupDB())
	checkErr(c, err, gin.H{"id": Id, "new address2": add2, "message": "you changed your address2!"})
}
func ModifyAdd1(c *gin.Context) {
	Id := getID(c)
	add1 := c.Param("add1")
	err := profile.Modify_address1(Id, add1, setupDB())
	checkErr(c, err, gin.H{"data": Id, "message": "you changed your address1!"})

}
func ModifyLastName(c *gin.Context) {
	Id := getID(c)
	lName := c.Param("lastName")
	err := profile.Modify_lastname(Id, lName, setupDB())
	checkErr(c, err, gin.H{"id": Id, "new_lastname": lName, "message": "you changed your last name!"})
}
func ModifyFirstName(c *gin.Context) {
	userId := getID(c)
	temp := c.Param("firstName")
	err := profile.Modify_firstname(userId, temp, setupDB())
	checkErr(c, err, gin.H{"id": userId, "new_firstName": temp, "message": "you changed your first name!"})
}
func ProductsByCategory(c *gin.Context) {
	catId := c.Param("catId")

	id, err := strconv.Atoi(catId)
	checkErr(c, err, gin.H{"data": id})
	db := setupDB()
	products, err := pickbuy.GetProductsByCat(id, db)
	checkErr(c, err, gin.H{"products": products, "category": catId, "message": "all products by category"})
}
func Logs(c *gin.Context) {
	Id := getID(c)
	var logs []structs.Logs
	logs, err := pickbuy.GetLogs(setupDB(), Id)

	checkErr(c, err, gin.H{"logs": logs, "userId": Id, "message": "your log history!"})
}

func AllCategories(c *gin.Context) {
	db := setupDB()
	allCat, err := pickbuy.GetAllCategories(db)
	checkErr(c, err, gin.H{"categories": allCat, "message": "All categories"})

}
func GetAllProductsByCategory(c *gin.Context) {
	db := setupDB()
	allCat, err := pickbuy.GetAllCategories(db)
	checkErr(c, err, gin.H{"categories": allCat, "message": "All categories"})

	var cat structs.Categories
	for i := range allCat {
		cat = allCat[i]
		id := cat.Cat_id
		db := setupDB()
		products, err := pickbuy.GetProductsByCat(id, db)
		checkErr(c, err, gin.H{"products": products, "category": id, "message": "all products by all category"})

	}

}
func AllProducts(c *gin.Context) {
	products, err := pickbuy.GetAllProducts(setupDB())
	checkErr(c, err, gin.H{"products": products, "message": "All products"})
}
func GetAllOrders(c *gin.Context) {
	id := getID(c)
	orders, sum, err := pickbuy.GetAllOrders(setupDB(), id)
	checkErr(c, err, gin.H{"orders": orders, "userId": id, "total_sum": sum, "message": "All orders"})
}
func GetRecommend(c *gin.Context) {
	productId := c.Param("productId")
	id, err := strconv.Atoi(productId)
	checkErr(c, err, gin.H{"data": id})
	db := setupDB()
	var p structs.Product
	result, err := db.Query("SELECT * FROM products WHERE product_id=?", productId)
	for result.Next() {

		err = result.Scan(&p.Product_id, &p.Product_cat, &p.Product_brand, &p.Product_title,
			&p.Product_price, &p.Product_desc, &p.Product_image, &p.Product_keywords, &p.Product_count)
	}
	products, err := pickbuy.RecommendProducts(db, p)

	checkErr(c, err, gin.H{"products": products, "userId": id, "message": "we recommend you these products"})
}
func DeleteOrder(c *gin.Context) {
	orderId := c.Param("orderId")
	db := setupDB()
	_, err := db.Query(`DELETE FROM orders WHERE order_id=?`, orderId)
	checkErr(c, err, gin.H{"orderId": orderId, "message": "order deleted"})

}
func AddOrder(c *gin.Context) {
	userId := getID(c)

	pId, err := strconv.Atoi(c.Param("productId"))
	checkErr(c, err, gin.H{})

	amount, err := strconv.Atoi(c.Param("amount"))
	checkErr(c, err, gin.H{})

	result, err := setupDB().Query("SELECT * FROM products WHERE product_id=?", pId)
	checkErr(c, err, gin.H{})

	var p structs.Product
	for result.Next() {

		err = result.Scan(&p.Product_id, &p.Product_cat, &p.Product_brand, &p.Product_title,
			&p.Product_price, &p.Product_desc, &p.Product_image, &p.Product_keywords, &p.Product_count)
	}

	err = pickbuy.AddOrder(amount, p, setupDB(), userId)
	checkErr(c, err, gin.H{"product": p, "userId": userId, "massage": "order added successfully"})

}
func BuyOrders(c *gin.Context) {
	id := getID(c)

	cvv, err := strconv.Atoi(c.Param("cvv"))
	checkErr(c, err, gin.H{})

	cardNum, err := strconv.Atoi(c.Param("cardNumber"))
	checkErr(c, err, gin.H{})

	address := c.Param("address")

	err = pickbuy.Buy(setupDB(), cardNum, cvv, address, id)
	checkErr(c, err, gin.H{"userId": id, "cvv": cvv, "card num": cardNum, "address": address, "message": "orders bought"})
}
