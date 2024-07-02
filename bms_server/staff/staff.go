package staff

import (
	"fmt"

	"github.com/ronnie-w/bms-server/schema"
)

func CreateStaffTable() {
	staffTable := []string{
		"name varchar(255) not null",
		"password varchar(255) not null",

		"national_id_or_passport_number varchar(20)",
		"contact varchar(20)",
		"email varchar(255)",
		"residence_or_address varchar(255)",
		"job_title varchar(255)",
		"current_status varchar(20) default offline",
		"last_login varchar(19) not null",

		// misc
		"kra_pin varchar(20)",
		"medical_details varchar(20)",
		"social_security_details varchar(20)",
		"salary varchar(10)",
	}

	if err := schema.CreateTableQuery("staff", staffTable); err != nil {
		fmt.Println(err)
	}

	if err := schema.CreateTableQuery("admin", []string{
		"name varchar(5) not null",
		"email vachar(255) not null",
		"password varchar(255) not null",
	}); err != nil {
		fmt.Println(err)
	}

	StaffTimeManagementTable()
}

func StaffTimeManagementTable() {
	staffTimeTable := []string{
		"staff_id varchar(36)",
		"staff_name varchar(255)",
		"clock_in varchar(19)",
		"clock_out varchar(19)",
		"overtime_clock_in varchar(19)",
		"overtime_clock_out varchar(19)",
		"leave_days int",
	}

	if err := schema.CreateTableQuery("staff_time_table", staffTimeTable); err != nil {
		fmt.Println(err)
	}
}
