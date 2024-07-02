package sales

import (
	"fmt"

	"github.com/ronnie-w/bms-server/schema"
)

func CreateCustomerTable() {
	customersTable := []string{
		"customer_number varchar(20) not null",
		"customer_name varchar(255) not null",
	}

	if err := schema.CreateTableQuery("customers", customersTable); err != nil {
		fmt.Println(err)
	}
}
