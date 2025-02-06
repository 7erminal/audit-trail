package models

import (
	"errors"
	"fmt"
	"reflect"
	"strings"
	"time"

	"github.com/beego/beego/v2/client/orm"
)

type Audit_trail struct {
	AuditTrailId  int64     `orm:"auto"`
	Action        *Actions  `orm:"rel(fk)"`
	TableName     string    `orm:"size(150)"`
	ColumnChanged string    `orm:"size(150)"`
	Description   string    `orm:"size(255)"`
	DateChanged   time.Time `orm:"type(datetime)"`
	DateCreated   time.Time `orm:"type(datetime)"`
	DateModified  time.Time `orm:"type(datetime)"`
	CreatedBy     *Users    `orm:"rel(fk)"`
	ModifiedBy    *Users    `orm:"rel(fk)"`
}

func init() {
	orm.RegisterModel(new(Audit_trail))
}

// AddAudit_trail insert a new Audit_trail into database and returns
// last inserted Id on success.
func AddAudit_trail(m *Audit_trail) (id int64, err error) {
	o := orm.NewOrm()
	id, err = o.Insert(m)
	return
}

// GetAudit_trailById retrieves Audit_trail by Id. Returns error if
// Id doesn't exist
func GetAudit_trailById(id int64) (v *Audit_trail, err error) {
	o := orm.NewOrm()
	v = &Audit_trail{AuditTrailId: id}
	if err = o.QueryTable(new(Audit_trail)).Filter("AuditTrailId", id).RelatedSel().One(v); err == nil {
		return v, nil
	}
	return nil, err
}

// GetAllAudit_trail retrieves all Audit_trail matches certain condition. Returns empty list if
// no records exist
func GetAllAudit_trail(query map[string]string, fields []string, sortby []string, order []string,
	offset int64, limit int64) (ml []interface{}, err error) {
	o := orm.NewOrm()
	qs := o.QueryTable(new(Audit_trail))
	// query k=v
	for k, v := range query {
		// rewrite dot-notation to Object__Attribute
		k = strings.Replace(k, ".", "__", -1)
		qs = qs.Filter(k, v)
	}
	// order by:
	var sortFields []string
	if len(sortby) != 0 {
		if len(sortby) == len(order) {
			// 1) for each sort field, there is an associated order
			for i, v := range sortby {
				orderby := ""
				if order[i] == "desc" {
					orderby = "-" + v
				} else if order[i] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
			qs = qs.OrderBy(sortFields...)
		} else if len(sortby) != len(order) && len(order) == 1 {
			// 2) there is exactly one order, all the sorted fields will be sorted by this order
			for _, v := range sortby {
				orderby := ""
				if order[0] == "desc" {
					orderby = "-" + v
				} else if order[0] == "asc" {
					orderby = v
				} else {
					return nil, errors.New("Error: Invalid order. Must be either [asc|desc]")
				}
				sortFields = append(sortFields, orderby)
			}
		} else if len(sortby) != len(order) && len(order) != 1 {
			return nil, errors.New("Error: 'sortby', 'order' sizes mismatch or 'order' size is not 1")
		}
	} else {
		if len(order) != 0 {
			return nil, errors.New("Error: unused 'order' fields")
		}
	}

	var l []Audit_trail
	qs = qs.OrderBy(sortFields...).RelatedSel()
	if _, err = qs.Limit(limit, offset).All(&l, fields...); err == nil {
		if len(fields) == 0 {
			for _, v := range l {
				ml = append(ml, v)
			}
		} else {
			// trim unused fields
			for _, v := range l {
				m := make(map[string]interface{})
				val := reflect.ValueOf(v)
				for _, fname := range fields {
					m[fname] = val.FieldByName(fname).Interface()
				}
				ml = append(ml, m)
			}
		}
		return ml, nil
	}
	return nil, err
}

// UpdateAudit_trail updates Audit_trail by Id and returns error if
// the record to be updated doesn't exist
func UpdateAudit_trailById(m *Audit_trail) (err error) {
	o := orm.NewOrm()
	v := Audit_trail{AuditTrailId: m.AuditTrailId}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Update(m); err == nil {
			fmt.Println("Number of records updated in database:", num)
		}
	}
	return
}

// DeleteAudit_trail deletes Audit_trail by Id and returns error if
// the record to be deleted doesn't exist
func DeleteAudit_trail(id int64) (err error) {
	o := orm.NewOrm()
	v := Audit_trail{AuditTrailId: id}
	// ascertain id exists in the database
	if err = o.Read(&v); err == nil {
		var num int64
		if num, err = o.Delete(&Audit_trail{AuditTrailId: id}); err == nil {
			fmt.Println("Number of records deleted in database:", num)
		}
	}
	return
}
