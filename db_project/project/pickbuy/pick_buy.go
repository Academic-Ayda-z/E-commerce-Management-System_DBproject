package pickbuy

import (
	"database/sql"
	"db_p/profile"
	"db_p/structs"
	"fmt"
	"time"
)

func GetAllCategories(db *sql.DB) ([]structs.Categories, error) {
	var categories []structs.Categories
	var category structs.Categories

	result, err := db.Query("SELECT distinct cat_id, cat_title FROM categories ")
	if err != nil {
		return categories, err
	}
	for result.Next() {
		err = result.Scan(&category.Cat_id, &category.Cat_title)
		categories = append(categories, category)
		if err != nil {
			panic(err)
		}
	}
	return categories, err
}
func showAndChooseCategory(allCategories []structs.Categories) structs.Categories {
	println("Choose a category:")
	dictCategories := make(map[int]structs.Categories)
	var j int
	i := 0
	for i < len(allCategories) {
		dictCategories[i] = allCategories[i]
		//todo show more info
		fmt.Printf("%d ) %s\n", i+1, dictCategories[i].Cat_title)
		i += 1
	}
	fmt.Scanf("%d", &j)
	//	fmt.Printf("\n%d\n", dictCategories[j].Cat_id)
	return dictCategories[j-1]
}
func GetProductsByCat(category int, db *sql.DB) ([]structs.Product, error) {
	var products []structs.Product

	rows, err := db.Query("SELECT * FROM products WHERE product_cat=?", category)
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var p structs.Product

		err = rows.Scan(&p.Product_id, &p.Product_cat, &p.Product_brand, &p.Product_title,
			&p.Product_price, &p.Product_desc, &p.Product_image, &p.Product_keywords, &p.Product_count)

		if err != nil {

			return products, err
		}
		products = append(products, p)
	}
	return products, err

}
func GetAllOrders(db *sql.DB, id int) ([]structs.Order_products, int, error) {
	rows, err := db.Query("SELECT * FROM orders WHERE user_id=?", id)
	sum := 0
	if err != nil {
		panic(err)
	}
	var orders []structs.Order_products
	for rows.Next() {
		var order structs.Order_products
		err = rows.Scan(&order.Order_id, &order.User_id, &order.Product_id, &order.Amt, &order.Qty, &order.Status)
		if err != nil {
			//			panic(err)
			return orders, sum, err
		}
		orders = append(orders, order)
		sum += order.Amt
	}
	return orders, sum, err
}
func GetAllProducts(db *sql.DB) ([]structs.Product, error) {
	var products []structs.Product

	rows, err := db.Query("SELECT * FROM products")
	if err != nil {
		panic(err)
	}

	for rows.Next() {
		var p structs.Product

		err = rows.Scan(&p.Product_id, &p.Product_cat, &p.Product_brand, &p.Product_title,
			&p.Product_price, &p.Product_desc, &p.Product_image, &p.Product_keywords, &p.Product_count)

		if err != nil {

			return products, err
		}
		products = append(products, p)
	}
	return products, err

}
func showAndChooseProduct(products []structs.Product) structs.Product {

	println("\nchoose a product:")
	println("title\t\t\t price")
	dictProducts := make(map[int]structs.Product)
	var j int

	i := 0
	for i < len(products) {
		dictProducts[i] = products[i]
		//todo show more info
		fmt.Printf("%d ) %s\t%d\t\n", i+1, dictProducts[i].Product_title, dictProducts[i].Product_price)
		i += 1
	}
	fmt.Scanf("%d", &j)
	return dictProducts[j]

}
func InformProducts(db *sql.DB, id int) {
	result, err := db.Query(`SELECT order_id,product_id,total_amt  FROM orders WHERE user_id=? and status=? `, id, "SELECTION FAILED")
	if err != nil {
		panic(err)
	}
	for result.Next() {
		var pID int
		var amount int
		var orderID int
		err = result.Scan(&orderID, &pID, &amount)
		if err != nil {
			panic(err)
		}
		var p structs.Product
		result2, err := db.Query(`SELECT *  FROM products WHERE  product_id=?`, pID)
		if err != nil {
			panic(err)
		}
		err = result2.Scan(&p.Product_id, &p.Product_cat, &p.Product_brand, &p.Product_title,
			&p.Product_price, &p.Product_desc, &p.Product_image, &p.Product_keywords, &p.Product_count)
		if err != nil {
			panic(err)
		}

		if p.Product_count >= amount {
			fmt.Printf("product %s is available now\n", p.Product_title)
			_, err = db.Query(`DELETE FROM orders  WHERE order_id=? `, orderID)
		} else {
			fmt.Printf("product %s is not available yet!\n", p.Product_title)
		}
	}
}
func selectProduct(db *sql.DB, id int, product *structs.Product, chosenProducts *[]Pair, sumPrices *int) {
	//check if the product is available ?
	var amount int
	println("how many?")
	fmt.Scanf("%d", &amount)

	if product.Product_count >= amount {
		*sumPrices += product.Product_price * amount
		*chosenProducts = append(*chosenProducts, Pair{*product, amount})

		//update product_ decrease count of product
		_, err := db.Query(`UPDATE products SET product_count = ? WHERE product_id = ?`,
			product.Product_count-amount, product.Product_id)

		if err != nil {
			panic(err)
		}

	} else {
		fmt.Printf("product %s with count %d\n", product.Product_title, amount)
		_, err := db.Query(`INSERT INTO orders  SET user_id=?,product_id=?,total_amt=?,count=?,status=?  `,
			id, product.Product_id, amount*product.Product_price, amount, "SELECTION FAILED")
		if err != nil {
			panic(err)
		}

		//inform the costumer when it became available func :pickbuy.InformProducts(db, id)

		println("this product is unavailable for now. check your profile for any news\n")
	}

}
func ShowOrders(selectedProducts []Pair, totalAmount int) {
	println("product title\t\tprice\t\tamount\n")
	i := 0
	for i < len(selectedProducts) {
		p := selectedProducts[i]
		fmt.Printf("%d ) %s %d \t\t%d\n", i+1, p.product.Product_title, p.product.Product_price, p.amount)
		i += 1
	}
	fmt.Printf("total sum :\t\t %d\n", totalAmount)

}

type Pair struct {
	product structs.Product
	amount  int
}

func RecommendProducts(db *sql.DB, product structs.Product) ([]structs.Product, error) {
	result, err := db.Query(`SELECT * FROM products WHERE product_cat=? and product_title!=?`,
		product.Product_cat, product.Product_title)
	if err != nil {
		panic(err)
	}

	//var products []structs.Product
	println("recommended products:")
	println("product title\t\tamount\t\tbrand\n")
	var products []structs.Product
	i := 0
	for result.Next() {
		if i >= 4 {
			break
		}
		var p structs.Product

		err = result.Scan(&p.Product_id, &p.Product_cat, &p.Product_brand, &p.Product_title,
			&p.Product_price, &p.Product_desc, &p.Product_image, &p.Product_keywords, &p.Product_count)

		if err == nil {
			fmt.Printf("%d ) %s\t %d \t\t%d\n", i+1, p.Product_title, p.Product_price, p.Product_brand)
			products = append(products, p)
		}

		i += 1
	}
	return products, err
}
func Pick(db *sql.DB, id int) ([]Pair, int) {
	//show categories
	var sumPrices = 0
	var chosenProducts []Pair
	allCategories, _ := GetAllCategories(db)
	category := showAndChooseCategory(allCategories)
	fmt.Printf("%d", category.Cat_id)
	products, _ := GetProductsByCat(category.Cat_id, db)
	var exit string
	for true {
		product := showAndChooseProduct(products)
		println("wanna select this product?\n1)Yse\n2)No\n")
		var ans int
		fmt.Scanf("%d", &ans)
		if ans == 1 {
			fmt.Printf("product %s selected\n", product.Product_title)
			//show all info about it
			//todo
			//			getProductss()

			selectProduct(db, id, &product, &chosenProducts, &sumPrices)
			RecommendProducts(db, product)
		}
		println("to exit enter EXIT.") //yes of no
		fmt.Scanf("%s", &exit)
		if exit == "EXIT" {
			break
		}
	}
	return chosenProducts, sumPrices
}
func AddOrder(amount int, p structs.Product, db *sql.DB, id int) error {

	_, err := db.Query(`INSERT INTO orders  SET user_id=?,product_id=?,total_amt=?,count=?,status=?  `,
		id, p.Product_id, amount*p.Product_price, amount, "SELECTED")

	if err != nil {
		return err
		//			panic(err.Error()) // proper error handling instead of panic in your app
	}
	return err
}
func Order(db *sql.DB, id int) {
	products, _ := Pick(db, id)
	//showOrders(products, sumAmount)
	var p Pair
	i := 0
	for i < len(products) {
		p = products[i]
		err := AddOrder(p.amount, p.product, db, id)
		if err != nil {
			panic(err)
		}
		i += 1
	}
}
func chooseAddress(db *sql.DB, id int) string {
	var address [3]string
	result, err := db.Query(`SELECT address1,address2 FROM user_info WHERE user_id=?`, id)
	for result.Next() {
		err = result.Scan(&address[1], &address[2])
		if err != nil {
			panic(err)

		}
	}
	//choose address
	add := -1
	if address[1] != "X" {
		fmt.Printf("%d)address1:%s\n", 1, address[1])
		if address[2] != "X" {
			fmt.Printf("%d)address2:%s\n", 2, address[2])
			fmt.Scanf("%d", &add)
		}

	} else {
		println("please enter your address:")
		fmt.Scanf("%s", &address[0])
		profile.Modify_address1(id, address[0], db)

	}
	if add != -1 {
		address[0] = address[add]
	}
	return address[0]
}

func GetProductById(db *sql.DB, product_id int, p *structs.Product) {
	rows, err := db.Query("SELECT * FROM products WHERE product_id=?", product_id)
	if err != nil {
		panic(err)
	}

	for rows.Next() {

		err = rows.Scan(&p.Product_id, &p.Product_cat, &p.Product_brand, &p.Product_title,
			&p.Product_price, &p.Product_desc, &p.Product_image, &p.Product_keywords, &p.Product_count)

		if err != nil {
			panic(err)
		}
	}

}

func updateLog(db *sql.DB, id int, orderID int, totoalAmount int, cardNum int, cvv int, address string) {

	//make a log
	date := time.Now().Format("2006.01.02 15:04:05")
	_, err := db.Query(`INSERT INTO logs  SET order_id=?,user_id=?,action=?,address=?,total_amt=?,card_number=?,cvv=?,date=?`,
		orderID, id, "COMPLETED", address, totoalAmount, cardNum, cvv, date)
	if err != nil {
		panic(err.Error())
	}

	//update status

	_, err = db.Query(`UPDATE orders SET   status= ?	WHERE  order_id=? and user_id=?`,
		"COMPLETED", orderID, id)

	if err != nil {
		panic(err.Error())
	}

}
func Buy(db *sql.DB, cardNum int, cvv int, address string, id int) error {
	var count int
	var orderTotalCost int
	var orderID int
	var product_id int
	var arrOrderIDs []int
	sum := 0
	var products []Pair
	var product structs.Product
	//find selected products
	result, err := db.Query(`SELECT order_id,product_id,count,total_amt FROM orders WHERE status=? and user_id=? `, "SELECTED", id)

	if err != nil {
		panic(err.Error()) // proper error handling instead of panic in your app
	}

	for result.Next() {
		err = result.Scan(&orderID, &product_id, &count, &orderTotalCost)
		if err != nil {
			panic(err.Error()) // proper error handling instead of panic in your app
		}
		arrOrderIDs = append(arrOrderIDs, orderID)
		sum += orderTotalCost
		GetProductById(db, product_id, &product)
		products = append(products, Pair{product, count})
	}
	//show products
	println("\nall your orders:")
	ShowOrders(products, sum)
	//var ans int
	//	println("\nAre you sure you want to buy?\n0)NO\n1)YES\n")
	//fmt.Scanf("%d", &ans)
	//if ans == 1 {
	//log
	//	var cardNum int
	//	var cvv int

	//	println("enter your card number:")
	//	fmt.Scanf("%d", &cardNum)

	//	println("enter your cvv:")
	//	fmt.Scanf("%d", &cvv)

	//	address := chooseAddress(db, id)

	i := 0
	for i < len(arrOrderIDs) {
		//i is order id
		updateLog(db, id, arrOrderIDs[i], sum, cardNum, cvv, address)
		_, err = db.Query("DELETE FROM orders WHERE order_id=?", arrOrderIDs[i])
		if err != nil {
			panic(err)
		}
		i += 1
	}
	//}
	//	ShowLogs(id, db)
	//	GetLogs(db, id)
	return err
}
func GetLogs(db *sql.DB, id int) ([]structs.Logs, error) {
	var logs []structs.Logs
	result, err := db.Query(`SELECT * FROM logs WHERE user_id=?`, id)
	if err != nil {
		panic(err.Error())
	}
	for result.Next() {
		var log structs.Logs
		err = result.Scan(&log.Id, &log.Order_id, &log.User_id, &log.Action, &log.Address,
			&log.Total_amt, &log.Cardnumber, &log.Cvv, &log.Date)
		if err != nil {
			return logs, err
		}
		logs = append(logs, log)
	}
	return logs, err
}
func ShowLogs(id int, db *sql.DB) {
	logs, _ := GetLogs(db, id)
	i := 0
	for i < len(logs) {
		log := logs[i]
		fmt.Printf(" date:%s\naddress:%s\ncard%s\ntotal cost:%d\n",
			log.Date, log.Address, log.Cardnumber, log.Total_amt)
		i += 1
	}
}
