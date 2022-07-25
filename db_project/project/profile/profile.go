package profile

import (
	"database/sql"
	"db_p/structs"
	//"db_p/pickbuy"
	"errors"
	//"fmt"
	//"github.com/go-sql-driver/mysql"
	//"log"
)

//todo log out

//todo is this an admin or not   profiles are different

//are the orders available or not yet
func ShowNews(db *sql.DB, id int) []structs.Product {
	result, err := db.Query(`SELECT order_id,product_id,total_amt  FROM orders WHERE user_id=? and status=? `, id, "SELECTION FAILED")
	var products []structs.Product
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
			//fmt.Printf("product %s is available now\n", p.Product_title)
			_, err = db.Query(`DELETE FROM orders  WHERE order_id=? `, orderID)
			products = append(products, p)
		} else {
			//fmt.Printf("product %s is not available yet!\n", p.Product_title)
		}
	}
	return products
	//or //use -> pickbuy.InformProducts(db, id)
}
func ShowLogs(db *sql.DB, id int) {
	//todo  use ->pickbuy.ShowLog(id, db)
}
func Modify_email(id int, Email string, db *sql.DB) error {
	_, err := db.Query(`UPDATE user_info SET email = ? WHERE user_id = ?`, Email, id)
	if err != nil {
		errors.New("could not change email")
	}
	return err
}
func Modify_password(id int, Password string, db *sql.DB) error {

	_, err := db.Query(`UPDATE user_info SET password = ? WHERE user_id = ?`, Password, id)
	if err != nil {
		errors.New("could not change password")
	}
	return err

}
func Modify_firstname(id int, fname string, db *sql.DB) error {

	_, err := db.Query(`UPDATE user_info SET first_name = ?	WHERE user_id = ?`, fname, id)
	if err != nil {
		errors.New("could not change first_name")
	}
	return err
}
func Modify_lastname(id int, lname string, db *sql.DB) error {
	_, err := db.Exec(`UPDATE user_info SET last_name = ? WHERE user_id = ?`, lname, id)
	if err != nil {
		errors.New("could not change last_name")
	}
	return err
}
func Modify_mobile(id int, phnoeNum string, db *sql.DB) error {
	_, err := db.Query(`UPDATE user_info SET mobile = ? WHERE user_id = ?;`, phnoeNum, id)
	if err != nil {
		errors.New("could not change email")
	}
	return err
}
func Modify_address1(id int, add1 string, db *sql.DB) error {

	_, err := db.Query(`UPDATE user_info SET address1 = ? WHERE user_id = ?`, add1, id)
	if err != nil {
		errors.New("could not change address1")
	}
	return err
}
func Modify_address2(id int, add2 string, db *sql.DB) error {

	_, err := db.Query(`UPDATE user_info SET address2 = ? WHERE user_id = ?`, add2, id)
	if err != nil {
		errors.New("could not change address2")
	}
	return err
}
