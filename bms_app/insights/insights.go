package insights

import (
	"bytes"
	"encoding/json"
	"log"
	"strconv"
	"strings"
	"syscall/js"
	"time"
)

type salesData struct {
	Date, GrandTotal string
}

type timeline struct {
	Today, Week, Month, AllTime int
}

func fetchSalesTotal(sales []salesData) timeline {
	var salesTotal timeline

	for i := 0; i < len(sales); i++ {
		now := time.Now()
		date, err := time.ParseInLocation(time.DateTime, sales[i].Date, now.Location())
		if err != nil {
			log.Println(err)
		}
		grandTotal, _ := strconv.Atoi(sales[i].GrandTotal)

		difference, _ := strconv.Atoi(strings.Split(now.Sub(date).String(), "h")[0])

		nY, nM, nD := now.Date()
		dY, dM, dD := date.Date()

		switch {
		case nY == dY && nM == dM && nD == dD:
			salesTotal.Today += grandTotal
			salesTotal.Week += grandTotal
			salesTotal.Month += grandTotal
		case difference < 168:
			salesTotal.Week += grandTotal
			salesTotal.Month += grandTotal
		case difference < 720:
			salesTotal.Month += grandTotal
		}

		salesTotal.AllTime += grandTotal
	}

	return salesTotal
}

func FetchTotalSales(this js.Value, args []js.Value) interface{} {
	var (
		sales    []salesData
		salesDec salesData
		writer   bytes.Buffer

		dec = json.NewDecoder(strings.NewReader(args[0].String()))
	)

	dec.Token()
	for dec.More() {
		if err := dec.Decode(&salesDec); err != nil {
			log.Println(err)
		}

		sales = append(sales, salesDec)
	}
	dec.Token()

	json.NewEncoder(&writer).Encode(fetchSalesTotal(sales))

	return string(writer.Bytes())
}
