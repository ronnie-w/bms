package inventories

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/ronnie-w/bms-server/notifications"
	"github.com/ronnie-w/bms-server/sales"
	"github.com/ronnie-w/bms-server/schema"
	"github.com/ronnie-w/bms-server/utils"
)

func CreateInventoryTable() {
	inventoriesTable := []string{"inventory_name varchar(255) not null", "inventory_value float(2) not null"}
	if err := schema.CreateTableQuery("inventories", inventoriesTable); err != nil {
		fmt.Println(err)
	}

	sales.CreateSalesTable()
}

func CreateInventory(rw http.ResponseWriter, r *http.Request) {
	var (
		err              error
		inventory        = schema.GeneralDecoder(r)
		inventoryColumns = []string{
			"item_name varchar(255) not null",
			"description varchar(255) not null",
			"unit_cost float(2) not null",
			"quantity int not null",
			"total float(2) default 0",
			"recommended_selling_price float(2) not null",
			"reorder_quantity int not null",
			"reorder bool default false",
			"supplier_or_manufacturer varchar(255) default unknown",
			"supplier_email varchar(255) default none",
			"supplier_contact varchar(20) default none",
			"" + schema.QueryParser(inventory) + "",
		}
	)

	if insertErr := schema.InsertQuery("inventories", []interface{}{inventory.Args1.Arg1Val, 0}); insertErr.Err != nil {
		utils.ErrResponse(rw, nil, insertErr)
	} else {
		err = schema.CreateTableQuery(inventory.Args1.Arg1Val, inventoryColumns)
		if err != nil {
			fmt.Println(err)
			utils.ErrResponse(rw, err, "Could not create inventory")
		} else {
			utils.SuccessResponse(rw, "Inventory: "+inventory.Args1.Arg1Val+" created successfully")
		}
	}

}

func InsertToInventory(rw http.ResponseWriter, r *http.Request) {
	table, insertRes := utils.DefaultInsert(r)
	if insertRes.Err != nil {
		utils.ErrResponse(rw, nil, insertRes.Err)
	} else {
		updateInventoryVal(table, insertRes.Uuid)
		utils.SuccessResponse(rw, "Insert operation complete without errors")
	}
}

func UpdateInInventory(rw http.ResponseWriter, r *http.Request) {
	table := schema.GeneralDecoder(r).Args1.Arg1Val
	uuid := strings.Split(r.URL.Path, "/")[2]
	var (
		unitCost, quantity = fetchCostQ(table, uuid)
		inventoryValue     = fetchInventoryVal(table)
		updateReturn       interface{}
	)

	inventoryValue -= unitCost * float32(quantity)
	inventoryValUpdate(inventoryValue, table)

	updateReturn = schema.UpdateQuery(uuid, schema.GeneralDecoder(r))
	updateInventoryVal(table, uuid)

	utils.ErrOrWarn(rw, updateReturn, "Update operation complete without errors", "Could not update item "+uuid+" in record")
}

func DeleteFromInventory(rw http.ResponseWriter, r *http.Request) {
	uuid := strings.Split(r.URL.Path, "/")[2]
	var (
		total          float32
		deleteData     = schema.GeneralDecoder(r)
		inventoryValue = fetchInventoryVal(deleteData.Args1.Arg1Val)
	)

	err := schema.DevSelect([]string{"total"}, deleteData.Args1.Arg1Val, "uuid", uuid, []interface{}{&total})
	if err != nil {
		fmt.Println(err)
	}
	inventoryValUpdate(inventoryValue-total, deleteData.Args1.Arg1Val)

	if err := schema.DeleteQuery(uuid, deleteData); err != nil {
		utils.ErrResponse(rw, err, "Could not delete item "+uuid+" from record")
	} else {
		utils.SuccessResponse(rw, "Delete operation complete without errors")
	}
}

func fetchInventoryVal(table string) float32 {
	var inventoryValue float32
	if err := schema.DevSelect([]string{"inventory_value"}, "inventories", "inventory_name", table, []interface{}{&inventoryValue}); err != nil {
		fmt.Println(err)
	}

	return inventoryValue
}

func inventoryValUpdate(inventoryValue, table interface{}) {
	err := schema.DevUpdate("inventories", "inventory_name", []string{"inventory_value"}, []interface{}{inventoryValue, table})
	if err != nil {
		fmt.Println(err)
	}
}

func fetchCostQ(table, uuid string) (float32, int) {
	var (
		unitCost float32
		quantity int
	)

	if err := schema.DevSelect([]string{"unit_cost", "quantity"}, table, "uuid", uuid, []interface{}{&unitCost, &quantity}); err != nil {
		fmt.Println(err)
	}
	return unitCost, quantity
}

func updateInventoryVal(table, uuid string) {
	var (
		unitCost, quantity                                     = fetchCostQ(table, uuid)
		inventoryValue                                         = fetchInventoryVal(table)
		reorderQuantity                                        int
		reorder                                                bool
		itemName, supplierName, supplierContact, supplierEmail string
	)

	inventoryValUpdate(inventoryValue+unitCost*float32(quantity), table)
	err := schema.DevUpdate(table, "uuid", []string{"total"}, []interface{}{unitCost * float32(quantity), uuid})
	if err != nil {
		fmt.Println(err)
	}

	err = schema.DevSelect(
		[]string{"reorder", "reorder_quantity", "supplier_or_manufacturer", "supplier_contact", "supplier_email", "item_name"},
		table, "uuid", uuid,
		[]interface{}{&reorder, &reorderQuantity, &supplierName, &supplierContact, &supplierEmail, &itemName},
	)
	if err != nil {
		fmt.Println(err)
	}

	if quantity < reorderQuantity {
		if reorder {
			InventoryReorder(supplierName, supplierContact, supplierEmail)
		}
		notifications.NewNotification(""+itemName+" in "+table+" inventory is running low, "+fmt.Sprintf("%v", quantity)+" units remaining in stock", "high")
	}
}
