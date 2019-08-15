package main

import (
	"database/sql"
	"fmt"
	pb "github.com/GoogleCloudPlatform/microservices-demo/src/productcatalogservice/genproto"
	_ "github.com/go-sql-driver/mysql"
)

const DRIVER = "mysql"

type MySQL struct {
	db *sql.DB
}
type ProductMysql struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	Picture      string `json:"picture"`
	CurrencyCode string `json:"price_usd_currency_code"`
	Units        int64  `json:"price_usd_units"`
	Nanos        int32  `json:"price_usd_nanos"`
	Categories   string `json:"categories"`
}

func NewSQLConnection() (*MySQL, error) {
	this := MySQL{}
	err := this.establishConnection()
	if err != nil {
		return nil, err
	}
	return &this, nil
}
func (m *MySQL) establishConnection() (err error) {
	connectionString := fmt.Sprintf("%s:%s@tcp(%s)/%s", UserName, Password, Endpoint, DbName)
	db, err := sql.Open(DRIVER, connectionString)
	if err != nil {
		log.Errorf("enable to connect to %s database. error: %s", DbName, err.Error())
		return err
	}
	m.db = db
	return nil
}

func (m *MySQL) ReadProducts() (*pb.ListProductsResponse, error) {

	// Execute the query
	results, err := m.db.Query("SELECT * FROM productcatalog")
	if err != nil {
		log.Errorf("unable to read data from mysql. Erorr: %s", err.Error())
		return nil, err
	}
	productsList := new(pb.ListProductsResponse)
	for results.Next() {
		var product = new(pb.Product)
		product.PriceUsd = new(pb.Money)
		// for each row, scan the result into our tag composite object
		category := ""
		err = results.Scan(&product.Id, &product.Name, &product.Description, &product.Picture,
			&product.PriceUsd.CurrencyCode, &product.PriceUsd.Units, &product.PriceUsd.Nanos,
			&category)
		if err != nil {
			log.Errorf("unable to read data from mysql. Erorr: %s", err.Error())
			continue
		}
		product.Categories = append(product.Categories, category)
		// and then print out the tag's Name attribute
		log.Printf("%s", product.Id)
		productsList.Products = append(productsList.Products, product)
	}
	return productsList, err
}
