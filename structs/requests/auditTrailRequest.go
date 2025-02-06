package requests

type AuditTrailRequest struct {
	ActionId      int64
	ChangedBy     int64
	ColumnChanged string
	TableChanged  string
	DateChanged   string
}
