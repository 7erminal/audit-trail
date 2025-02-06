package main

import (
	"github.com/beego/beego/v2/client/orm/migration"
)

// DO NOT MODIFY
type AuditTrail_20250206_100140 struct {
	migration.Migration
}

// DO NOT MODIFY
func init() {
	m := &AuditTrail_20250206_100140{}
	m.Created = "20250206_100140"

	migration.Register("AuditTrail_20250206_100140", m)
}

// Run the migrations
func (m *AuditTrail_20250206_100140) Up() {
	// use m.SQL("CREATE TABLE ...") to make schema update
	m.SQL("CREATE TABLE audit_trail(`audit_trail_id` int(11) NOT NULL AUTO_INCREMENT,`action` int(11) NOT NULL,`table_name` varchar(150) NOT NULL,`column_changed` varchar(150) NOT NULL,`description` varchar(255) DEFAULT NULL,`date_created` datetime DEFAULT CURRENT_TIMESTAMP,`date_modified` datetime ON UPDATE CURRENT_TIMESTAMP,`created_by` int(11) NOT NULL,`modified_by` int(11) DEFAULT NULL,PRIMARY KEY (`audit_trail_id`), FOREIGN KEY (created_by) REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE NO ACTION, FOREIGN KEY (modified_by) REFERENCES users(user_id) ON UPDATE CASCADE ON DELETE NO ACTION, FOREIGN KEY (action) REFERENCES actions(action_id) ON UPDATE CASCADE ON DELETE NO ACTION)")
}

// Reverse the migrations
func (m *AuditTrail_20250206_100140) Down() {
	// use m.SQL("DROP TABLE ...") to reverse schema update
	m.SQL("DROP TABLE `audit_trail`")
}
