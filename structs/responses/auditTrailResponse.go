package responses

import "audit_trail_service/models"

type AuditTrailResponse struct {
	IdentificationTypeId int64
	Name                 string
	Code                 string
}

type AuditTrailsResponseDTO struct {
	StatusCode   int
	Audit_trails *[]interface{}
	StatusDesc   string
}

type AuditTrailResponseDTO struct {
	StatusCode  int
	Audit_trail *models.Audit_trail
	StatusDesc  string
}
