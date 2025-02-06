package controllers

import (
	"audit_trail_service/models"
	"audit_trail_service/structs/requests"
	"audit_trail_service/structs/responses"
	"encoding/json"
	"errors"
	"strconv"
	"strings"
	"time"

	"github.com/beego/beego/v2/core/logs"

	beego "github.com/beego/beego/v2/server/web"
)

// Audit_trailController operations for Audit_trail
type Audit_trailController struct {
	beego.Controller
}

// URLMapping ...
func (c *Audit_trailController) URLMapping() {
	c.Mapping("Post", c.Post)
	c.Mapping("GetOne", c.GetOne)
	c.Mapping("GetAll", c.GetAll)
	c.Mapping("Put", c.Put)
	c.Mapping("Delete", c.Delete)
}

// Post ...
// @Title Post
// @Description create Audit_trail
// @Param	body		body 	requests.AuditTrailRequest	true		"body for Audit_trail content"
// @Success 201 {int} models.Audit_trail
// @Failure 403 body is empty
// @router / [post]
func (c *Audit_trailController) Post() {
	var v requests.AuditTrailRequest
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)

	var dateChanged time.Time
	proceed := false
	message := "An error occurred adding this audit request"
	statusCode := 308
	auditTrailResp := models.Audit_trail{}

	var allowedDateList [4]string = [4]string{"2006-01-02", "2006/01/02", "2006-01-02 15:04:05.000", "2006/01/02 15:04:05.000"}

	for _, date_ := range allowedDateList {
		logs.Debug("About to convert ", v.DateChanged)
		logs.Debug("About to convert ", c.Ctx.Input.Query("Dob"))
		// Convert dob string to date
		tdChanged, error := time.Parse(date_, v.DateChanged)

		if error != nil {
			logs.Error("Error parsing date", error)
			message = "Unable to determine date changed"
			proceed = false
		} else {
			logs.Error("Date converted to time successfully", tdChanged)
			dateChanged = tdChanged
			proceed = true

			break
		}
	}

	if proceed {
		if user, err := models.GetUsersById(v.ChangedBy); err == nil {
			if action, err := models.GetActionsById(v.ActionId); err == nil {
				var auditTrail models.Audit_trail = models.Audit_trail{Action: action, CreatedBy: user, DateChanged: dateChanged, ColumnChanged: v.ColumnChanged, TableName: v.TableChanged, ModifiedBy: user}
				if _, err := models.AddAudit_trail(&auditTrail); err == nil {
					message = "Audit trail added successfully"
					auditTrailResp = auditTrail
					c.Ctx.Output.SetStatus(200)
				} else {
					logs.Error("An error occurred adding audit trail ", err.Error())
					message = "An error occurred adding audit trail"
				}
			} else {
				logs.Error("An error occurred getting actions ", err.Error())
				message = "Action specified cannot be used"
			}
		} else {
			logs.Error("An error occurred getting user ", err.Error())
			message = "User does not exist"
		}
	}

	resp := responses.AuditTrailResponseDTO{StatusCode: statusCode, Audit_trail: &auditTrailResp, StatusDesc: message}
	c.Data["json"] = resp

	c.ServeJSON()
}

// GetOne ...
// @Title Get One
// @Description get Audit_trail by id
// @Param	id		path 	string	true		"The key for staticblock"
// @Success 200 {object} models.Audit_trail
// @Failure 403 :id is empty
// @router /:id [get]
func (c *Audit_trailController) GetOne() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v, err := models.GetAudit_trailById(id)

	statusCode := 308
	message := "An error occurred"
	if err != nil {
		logs.Error("An error occurred fetching audit trail ", err.Error())
		message = "An error occurred fetching audit trail"
		statusCode = 608
		resp := responses.AuditTrailResponseDTO{StatusCode: statusCode, Audit_trail: nil, StatusDesc: message}
		c.Data["json"] = resp
	} else {
		statusCode = 200
		resp := responses.AuditTrailResponseDTO{StatusCode: statusCode, Audit_trail: v, StatusDesc: message}
		c.Data["json"] = resp
	}
	c.ServeJSON()
}

// GetAll ...
// @Title Get All
// @Description get Audit_trail
// @Param	query	query	string	false	"Filter. e.g. col1:v1,col2:v2 ..."
// @Param	fields	query	string	false	"Fields returned. e.g. col1,col2 ..."
// @Param	sortby	query	string	false	"Sorted-by fields. e.g. col1,col2 ..."
// @Param	order	query	string	false	"Order corresponding to each sortby field, if single value, apply to all sortby fields. e.g. desc,asc ..."
// @Param	limit	query	string	false	"Limit the size of result set. Must be an integer"
// @Param	offset	query	string	false	"Start position of result set. Must be an integer"
// @Success 200 {object} models.Audit_trail
// @Failure 403
// @router / [get]
func (c *Audit_trailController) GetAll() {
	var fields []string
	var sortby []string
	var order []string
	var query = make(map[string]string)
	var limit int64 = 10
	var offset int64

	// fields: col1,col2,entity.col3
	if v := c.GetString("fields"); v != "" {
		fields = strings.Split(v, ",")
	}
	// limit: 10 (default is 10)
	if v, err := c.GetInt64("limit"); err == nil {
		limit = v
	}
	// offset: 0 (default is 0)
	if v, err := c.GetInt64("offset"); err == nil {
		offset = v
	}
	// sortby: col1,col2
	if v := c.GetString("sortby"); v != "" {
		sortby = strings.Split(v, ",")
	}
	// order: desc,asc
	if v := c.GetString("order"); v != "" {
		order = strings.Split(v, ",")
	}
	// query: k:v,k:v
	if v := c.GetString("query"); v != "" {
		for _, cond := range strings.Split(v, ",") {
			kv := strings.SplitN(cond, ":", 2)
			if len(kv) != 2 {
				c.Data["json"] = errors.New("Error: invalid query key/value pair")
				c.ServeJSON()
				return
			}
			k, v := kv[0], kv[1]
			query[k] = v
		}
	}

	statusCode := 308
	message := "An error occurred"

	l, err := models.GetAllAudit_trail(query, fields, sortby, order, offset, limit)
	if err != nil {
		logs.Error("An error occurred fetching audit trails ", err.Error())
		message = "An error occurred fetching audit trails"
		statusCode = 608
		resp := responses.AuditTrailsResponseDTO{StatusCode: statusCode, Audit_trails: nil, StatusDesc: message}
		c.Data["json"] = resp
		c.Data["json"] = err.Error()
	} else {
		statusCode = 200
		resp := responses.AuditTrailsResponseDTO{StatusCode: statusCode, Audit_trails: &l, StatusDesc: message}
		c.Data["json"] = resp
		c.Data["json"] = l
	}
	c.ServeJSON()
}

// Put ...
// @Title Put
// @Description update the Audit_trail
// @Param	id		path 	string	true		"The id you want to update"
// @Param	body		body 	models.Audit_trail	true		"body for Audit_trail content"
// @Success 200 {object} models.Audit_trail
// @Failure 403 :id is not int
// @router /:id [put]
func (c *Audit_trailController) Put() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	v := models.Audit_trail{AuditTrailId: id}
	json.Unmarshal(c.Ctx.Input.RequestBody, &v)
	if err := models.UpdateAudit_trailById(&v); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}

// Delete ...
// @Title Delete
// @Description delete the Audit_trail
// @Param	id		path 	string	true		"The id you want to delete"
// @Success 200 {string} delete success!
// @Failure 403 id is empty
// @router /:id [delete]
func (c *Audit_trailController) Delete() {
	idStr := c.Ctx.Input.Param(":id")
	id, _ := strconv.ParseInt(idStr, 0, 64)
	if err := models.DeleteAudit_trail(id); err == nil {
		c.Data["json"] = "OK"
	} else {
		c.Data["json"] = err.Error()
	}
	c.ServeJSON()
}
