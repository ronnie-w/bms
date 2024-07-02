package utils

import (
	"fmt"

	"github.com/ronnie-w/bms-server/schema"
)

func CreateSettingsTable() {
	settingsTable := []string{
		"mpesa_till_number varchar(15) default NA",
		"qr_code_scanner_vendor_id varchar(20) default NA",
		"qr_code_scanner_device_id varchar(20) default NA",
		"printer_vendor_id varchar(20) default NA",
		"printer_device_id varchar(20) default NA",
	}

	if err := schema.CreateTableQuery("settings", settingsTable); err != nil {
		fmt.Println(err)
	}
}
