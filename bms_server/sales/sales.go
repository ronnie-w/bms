package sales

import (
	"fmt"

	"github.com/ronnie-w/bms-server/schema"
)

func CreateSalesTable() {
	salesTable := []string{
		"sales_id int not null",
		"server_id varchar(36) not null",
		"customer_number varchar(20) not null",
		"item_uuid varchar(36) not null",
		"item_name varchar(255) not null",
		"item_inventory varchar(255) not null",
		"item_price float(2) not null",
		"quantity_purchased int not null",
		"remaining_quantity int not null",
		"purchase_cost float(2) not null",
		"grand_total float(2) not null",
		"transaction_type varchar(50) not null",
		"transaction_id varchar(255) default NA",
		"qr_code_value varchar(255) not null default NA",
	}

	if err := schema.CreateTableQuery("sales", salesTable); err != nil {
		fmt.Println(err)
	}
}
