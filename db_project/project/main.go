package main

import (
	API "db_p/Api"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	"log"
)

func main() {

	r := gin.Default()
	r.GET("/signUp/:email/:password", API.CreatUser)
	r.GET("/signIn/:email/:password", API.SignInUser)

	r.GET("/profile/:id/logs", API.Logs)
	r.GET("/profile/:id/newProduct", API.News)
	r.GET("profile/modifyFName/:id/:firstName", API.ModifyFirstName)
	r.GET("profile/modifyLName/:id/:lastName", API.ModifyLastName)
	r.GET("profile/modifyEmail/:id/:email", API.ModifyEmail)
	r.GET("profile/modifyPass/:id/:password", API.ModifyPassword)
	r.GET("profile/modifyPhone/:id/:mobile", API.ModifyMobile)
	r.GET("profile/modifyAdd1/:id/:add1", API.ModifyAdd1)
	r.GET("profile/modifyAdd2/:id/:add2", API.ModifyAdd2)

	r.GET("/allCategories", API.AllCategories)
	r.GET("allCategories/:catId", API.ProductsByCategory)
	r.GET("/products", API.AllProducts)
	r.GET("/products/:productId", API.GetAllProductsByCategory)
	r.GET("/products/recommends/:productId", API.GetRecommend)

	r.GET("/orders/:id", API.GetAllOrders)
	r.GET("/order/:id/:productId/:amount", API.AddOrder)
	r.GET("/order/delete/:orderId", API.DeleteOrder)

	r.GET("/buyOrders/:id/:cvv/:cardNumber/:address", API.BuyOrders)

	log.Println("Listening to port 8080...")
	log.Fatal(r.Run(":8080"))

	//id, _ := signUP_IN.SignIn(db)
	//pickbuy.Order(db, id)

	//pickbuy.Buy(db, id)
}
