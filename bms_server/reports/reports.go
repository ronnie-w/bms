package reports

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/ronnie-w/bms-server/database"
	"github.com/ronnie-w/bms-server/schema"
	"github.com/ronnie-w/bms-server/utils"
)

var (
	json = jsoniter.ConfigCompatibleWithStandardLibrary
	db   = database.Conn()
)

type TimeLine struct {
	Today, Week, Month, AllTime int
}

type Product struct {
	ProductName                                        string
	ProductPrice                                       string
	ProductSalesQuantity, ProductSales, ProductProfits TimeLine
}

type SalesReport struct {
	Sales, Profits          TimeLine
	ProductSalesAndGrossing map[string]*Product
}

func fetchSalesProfits(sales []map[string]string) SalesReport {
	var (
		salesTotal, profitsTotal TimeLine

		products = make(map[string]*Product, len(sales))
	)

	for i := 0; i < len(sales); i++ {
		products[sales[i]["item_uuid"]] = &Product{
			ProductName:  sales[i]["item_name"],
			ProductPrice: sales[i]["item_price"],
		}
	}

	for i := 0; i < len(sales); i++ {
		now := time.Now()
		date, err := time.ParseInLocation(time.DateTime, sales[i]["date_created"], now.Location())
		if err != nil {
			fmt.Println(err)
		}

		quantityPurchased, _ := strconv.Atoi(sales[i]["quantity_purchased"])
		grandTotal, _ := strconv.Atoi(sales[i]["grand_total"])
		purchaseCost, _ := strconv.Atoi(sales[i]["purchase_cost"])
		profit := grandTotal - purchaseCost

		timeDiff, _ := strconv.Atoi(strings.Split(now.Sub(date).String(), "h")[0])

		nY, nM, nD := now.Date()
		dY, dM, dD := date.Date()

		switch {
		case nY == dY && nM == dM && nD == dD:
			products[sales[i]["item_uuid"]].ProductSales.Today += grandTotal
			products[sales[i]["item_uuid"]].ProductSales.Week += grandTotal
			products[sales[i]["item_uuid"]].ProductSales.Month += grandTotal
			products[sales[i]["item_uuid"]].ProductProfits.Today += profit
			products[sales[i]["item_uuid"]].ProductProfits.Week += profit
			products[sales[i]["item_uuid"]].ProductProfits.Month += profit

			products[sales[i]["item_uuid"]].ProductSalesQuantity.Today += quantityPurchased
			products[sales[i]["item_uuid"]].ProductSalesQuantity.Week += quantityPurchased
			products[sales[i]["item_uuid"]].ProductSalesQuantity.Month += quantityPurchased

			salesTotal.Today += grandTotal
			salesTotal.Week += grandTotal
			salesTotal.Month += grandTotal
			profitsTotal.Today += profit
			profitsTotal.Week += profit
			profitsTotal.Month += profit
		case timeDiff < 168:
			products[sales[i]["item_uuid"]].ProductSales.Week += grandTotal
			products[sales[i]["item_uuid"]].ProductSales.Month += grandTotal
			products[sales[i]["item_uuid"]].ProductProfits.Week += profit
			products[sales[i]["item_uuid"]].ProductProfits.Month += profit

			products[sales[i]["item_uuid"]].ProductSalesQuantity.Week += quantityPurchased
			products[sales[i]["item_uuid"]].ProductSalesQuantity.Month += quantityPurchased

			salesTotal.Week += grandTotal
			salesTotal.Month += grandTotal
			profitsTotal.Week += profit
			profitsTotal.Month += profit
		case timeDiff < 720:
			products[sales[i]["item_uuid"]].ProductSales.Month += grandTotal
			products[sales[i]["item_uuid"]].ProductProfits.Month += profit

			products[sales[i]["item_uuid"]].ProductSalesQuantity.Month += quantityPurchased

			salesTotal.Month += grandTotal
			profitsTotal.Month += profit
		}

		products[sales[i]["item_uuid"]].ProductSales.AllTime += grandTotal
		products[sales[i]["item_uuid"]].ProductProfits.AllTime += profit

		products[sales[i]["item_uuid"]].ProductSalesQuantity.AllTime += quantityPurchased

		salesTotal.AllTime += grandTotal
		profitsTotal.AllTime += profit
	}

	return SalesReport{Sales: salesTotal, Profits: profitsTotal, ProductSalesAndGrossing: products}
}

func FetchReports(rw http.ResponseWriter, r *http.Request) {
	var (
		rows         *sql.Rows
		err          error
		salesMapData []map[string]string

		uniqueVal = schema.GeneralDecoder(r)
	)

	if len(uniqueVal.Args2) > 0 {
		for _, u := range uniqueVal.Args2 {
			rows, err = db.Query("SELECT * FROM sales WHERE "+u.Arg2Val1+"=?", u.Arg2Val2)
			if err != nil {
				fmt.Println(err)
				utils.ErrResponse(rw, err, "Could not fetch data from table sales")
			}
		}
	} else {
		rows, err = db.Query("SELECT * FROM sales")
		if err != nil {
			fmt.Println(err)
			utils.ErrResponse(rw, err, "Could not fetch data from table sales")
		}
	}

	salesData := schema.SelectRowsScan(r, rows)
	for _, sd := range salesData {
		salesMap := make(map[string]string, 15)
		for i := 0; i < len(sd); i++ {
			salesMap[sd[i].ColName] = sd[i].ColValue
		}
		salesMapData = append(salesMapData, salesMap)
	}

	json.NewEncoder(rw).Encode(fetchSalesProfits(salesMapData))
}

type salesAndProfits struct {
	Sales, Profits, QuantitySold int
}

func GraphSalesAndProfits(rw http.ResponseWriter, r *http.Request) {
	var (
		rows      *sql.Rows
		err       error
		csvString string = "Date,Sales,Profits,Quantity Sold\n"
	)
	uniqueVal := schema.GeneralDecoder(r)

	if len(uniqueVal.Args2) == 1 {
		rows, err = db.Query("SELECT * FROM sales WHERE "+uniqueVal.Args2[0].Arg2Val1+"=?", uniqueVal.Args2[0].Arg2Val2)
		if err != nil {
			fmt.Println(err)
			utils.ErrResponse(rw, err, "Could not fetch data from table sales")
		}
	} else if len(uniqueVal.Args2) == 2 {
		rows, err = db.Query(`SELECT * FROM sales WHERE 
					`+uniqueVal.Args2[0].Arg2Val1+`=? AND 
					`+uniqueVal.Args2[1].Arg2Val1+` LIKE '%`+uniqueVal.Args2[1].Arg2Val2+`%'`, uniqueVal.Args2[0].Arg2Val2)
		if err != nil {
			fmt.Println(err)
			utils.ErrResponse(rw, err, "Could not fetch data from table sales")
		}
	} else {
		rows, err = db.Query("SELECT * FROM sales")
		if err != nil {
			fmt.Println(err)
			utils.ErrResponse(rw, err, "Could not fetch data from table sales")
		}
	}

	salesData := schema.SelectRowsScan(r, rows)
	salesMapData := make(map[string]*salesAndProfits, len(salesData))

	for i := 0; i < len(salesData); i++ {
		salesMapData[strings.Split(salesData[i][2].ColValue, " ")[0]] = &salesAndProfits{Sales: 0, Profits: 0, QuantitySold: 0}
	}

	for _, sales := range salesData {
		fmt.Println(sales)
		grandTotal, _ := strconv.Atoi(sales[13].ColValue)
		purchaseCost, _ := strconv.Atoi(sales[12].ColValue)
		quantitySold, _ := strconv.Atoi(sales[10].ColValue)

		salesMapData[strings.Split(sales[2].ColValue, " ")[0]].Sales += grandTotal
		salesMapData[strings.Split(sales[2].ColValue, " ")[0]].Profits += grandTotal - purchaseCost
		salesMapData[strings.Split(sales[2].ColValue, " ")[0]].QuantitySold += quantitySold
	}

	for k, v := range salesMapData {
		csvString += k + "," + fmt.Sprint(v.Sales) + "," + fmt.Sprint(v.Profits) + "," + fmt.Sprint(v.QuantitySold) + "\n"
	}

	json.NewEncoder(rw).Encode(csvString)
}
