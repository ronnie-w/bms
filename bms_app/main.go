package main

import (
	"syscall/js"

	"github.com/ronnie-w/bms-app/dialer"
	"github.com/ronnie-w/bms-app/schema"
	"github.com/ronnie-w/bms-app/utils"
)

func main() {
	jsFuncs := map[string]func(this js.Value, args []js.Value) interface{}{
		// authentication
		"Login":           dialer.DefaultDialer("/login_ep", false),
		"Signup":          dialer.DefaultDialer("/signup_ep", false),
		"FetchActiveUser": dialer.DefaultDialer("/active", false),
		"ResetPassword":   dialer.DefaultDialer("/reset_password_ep", false),

		// inventory functions
		"CreateInventory":     dialer.DefaultDialer("/create_inventory", false),
		"InsertToInventory":   dialer.DefaultDialer("/insert_to_inventory", false),
		"UpdateInInventory":   dialer.DefaultDialer("/update_in_inventory", true),
		"DeleteFromInventory": dialer.DefaultDialer("/delete_from_inventory", true),

		// module functions
		//tab functions
		"CreateNewTab":      dialer.DefaultDialer("/create_new_tab", false),
		"FetchTabs":         dialer.DefaultDialer("/fetch_tabs", false),
		"FetchRecordsInTab": dialer.DefaultDialer("/fetch_records_in_tab", false),
		//record functions
		"CreateNewRecord":  dialer.DefaultDialer("/create_new_record", true),
		"FetchRecordData":  dialer.DefaultDialer("/fetch_record_data", false),
		"InsertToRecord":   dialer.DefaultDialer("/insert_to_record", false),
		"UpdateInRecord":   dialer.DefaultDialer("/update_in_record", true),
		"DeleteFromRecord": dialer.DefaultDialer("/delete_from_record", true),
		//util functions
		"DropTable":  dialer.DefaultDialer("/drop_table", false),
		"GetColumns": dialer.DefaultDialer("/get_columns", false),
		"Exists":     dialer.DefaultDialer("/exists", false),

		// reports functions
		"FetchReports":         dialer.DefaultDialer("/fetch_reports", false),
		"GraphSalesAndProfits": dialer.DefaultDialer("/graph_sales_and_profits", false),

		// unique select function
		"UniqueSelectQuery": dialer.DefaultDialer("/unique_select_query", false),
		"LikeSelectQuery":   dialer.DefaultDialer("/like_select_query", false),

		"TransactMpesa":    dialer.DefaultDialer("/mpesa_transaction", false),
		"ConnectQrScanner": dialer.DefaultDialer("/connect_qr_scanner", false),

		"StrSanitize":      schema.StrSanitize,
		"StrParse":         schema.StrParse,
		"GenerateRandomId": utils.GenerateRandomId,
		"GetCurrentTime":   utils.GetCurrentTime,
	}

	for jsFunc, goFunc := range jsFuncs {
		js.Global().Set(jsFunc, js.FuncOf(goFunc))
	}

	<-make(chan bool)
}
