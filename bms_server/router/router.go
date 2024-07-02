package router

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/ronnie-w/bms-server/authentication"
	"github.com/ronnie-w/bms-server/devices"
	"github.com/ronnie-w/bms-server/inventories"
	"github.com/ronnie-w/bms-server/mpesa"
	"github.com/ronnie-w/bms-server/parser"
	"github.com/ronnie-w/bms-server/reports"
	"github.com/ronnie-w/bms-server/resources"
	"github.com/ronnie-w/bms-server/tabs"
	"github.com/ronnie-w/bms-server/utils"
)

func Routes() *mux.Router {
	mux := mux.NewRouter()

	endPoints := map[string]func(http.ResponseWriter, *http.Request){
		// authentication routes
		"/login_ep":          authentication.Login,
		"/signup_ep":         authentication.Signup,
		"/active":            authentication.FetchActiveUser,
		"/reset_password_ep": authentication.ResetPassword,

		// inventory routes
		"/create_inventory":             inventories.CreateInventory,
		"/insert_to_inventory":          inventories.InsertToInventory,
		"/update_in_inventory/{uuid}":   inventories.UpdateInInventory,
		"/delete_from_inventory/{uuid}": inventories.DeleteFromInventory,

		// module routes
		//tab routes
		"/create_new_tab":       authentication.Authenticate(tabs.CreateNewTab),
		"/fetch_tabs":           tabs.FetchTabs,
		"/fetch_records_in_tab": tabs.FetchRecordsInTab,
		//record routes
		"/create_new_record/{tab_name}": tabs.CreateNewRecord,
		"/fetch_record_data":            tabs.FetchRecordData,
		"/insert_to_record":             tabs.InsertToRecord,
		"/update_in_record/{uuid}":      tabs.UpdateInRecord,
		"/delete_from_record/{uuid}":    authentication.Authenticate(tabs.DeleteFromRecord),

		// reports routes
		"/fetch_reports":           authentication.Authenticate(reports.FetchReports),
		"/graph_sales_and_profits": authentication.Authenticate(reports.GraphSalesAndProfits),

		//util functions
		"/get_columns": utils.GetColumns,
		"/drop_table":  utils.DropTable,
		"/exists":      utils.Exists,

		// unique select route
		"/unique_select_query": utils.UniqueSelectQuery,
		"/like_select_query":   utils.LikeSelectQuery,

		"/mpesa_transaction":  mpesa.ConfirmC2BMpesa,
		"/connect_qr_scanner": devices.ConnectQRScanner,
	}

	templates := []string{"", "login", "admin_signup", "reset_password", "sales", "customers/{customers}"}
	for i := 0; i < len(templates); i++ {
		mux.HandleFunc("/"+templates[i], parser.Parser("index"))
	}

	templatesProtected := []string{"inventories/{inventoriesPath}", "tabs/{tabs}", "staff/{staff}", "settings/{settings}", "reports"}
	for i := 0; i < len(templatesProtected); i++ {
		mux.HandleFunc("/"+templatesProtected[i], authentication.Authenticate(parser.Parser("index")))
	}

	for route, httpFunc := range endPoints {
		mux.HandleFunc(route, httpFunc)
	}

	fileServer := http.FileServer(http.FS(resources.Resources))
	mux.PathPrefix("/").Handler(fileServer)

	return mux
}
