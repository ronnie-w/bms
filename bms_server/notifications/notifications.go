package notifications

import (
	"fmt"

	"github.com/ronnie-w/bms-server/schema"
)

func CreateNotificationsTable() {
	notificationsTable := []string{
		"notification varchar(255) not null",
		"priority varchar(8) not null",
		"viewed bool default false",
		"attended bool default false",
	}

	if err := schema.CreateTableQuery("notifications", notificationsTable); err != nil {
		fmt.Println(err)
	}
}

func NewNotification(notification, priority string) {
	if insertRes := schema.InsertQuery("notifications", []interface{}{notification, priority, false, false}); insertRes.Err != nil {
		fmt.Println(insertRes.Err)
	}
}
