package validator

import (
	"errors"
	"fmt"
	"mime/multipart"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/jinzhu/gorm"
	"github.com/thedevsaddam/govalidator"
)

// Covid19Validator main var
type Covid19Validator struct {
	DB *gorm.DB `json:"db"`
}

// CustomValidatorRules adds custom validator to govalidator
func (a *Covid19Validator) CustomValidatorRules() {
	// unique value on each column. format : []string{"unique:[table],[column],[exclude_column],[excluded_value]"}
	govalidator.AddCustomRule("unique", func(field string, rule string, message string, value interface{}) error {
		var (
			queryRow *gorm.DB
			total    int
		)
		var limit = 1

		query := `SELECT COUNT(*) as total FROM %s WHERE %s = ?`
		params := strings.Split(strings.TrimPrefix(rule, fmt.Sprintf("%s:", "unique")), ",")

		query = fmt.Sprintf(query, params[0], params[1])
		queryRow = a.DB.Raw(query, value)
		if len(params) == 3 {
			limit, _ = strconv.Atoi(params[2])
		}
		// if len(params) == 2 {
		//     query = fmt.Sprintf(query, params[0], params[1])
		//     queryRow = a.DB.Raw(query, value)
		// } else if len(params) == 4 {
		//     query += ` AND %s != ?`
		//     query = fmt.Sprintf(query, params[0], params[1], params[2])
		//     queryRow = a.DB.Raw(query, value, params[3])
		// } else {
		//     return fmt.Errorf("Arguments not enough")
		// }

		queryRow.Row().Scan(&total)

		if total >= limit {
			if message != "" {
				return errors.New(message)
			}

			return fmt.Errorf("The %s has already been taken", field)
		}

		return nil
	})
	// valid_id. must be a listed id of a model.
	govalidator.AddCustomRule("valid_id", func(field string, rule string, message string, value interface{}) error {
		var (
			db    *gorm.DB
			total int
		)

		table := strings.TrimPrefix(rule, fmt.Sprintf("%s:", "valid_id"))
		db = a.DB
		db.Table(table).
			Where("id IN (?)", value).
			Count(&total)

		if total < 1 {
			return fmt.Errorf(fmt.Sprintf("value %v is not found.", value), field)
		}
		return nil
	})

	// status. status of entity.
	govalidator.AddCustomRule("status", func(field string, rule string, message string, value interface{}) error {
		var (
			db    *gorm.DB
			total int
		)

		split := strings.Split(strings.TrimPrefix(rule, fmt.Sprintf("%s:", "status")), ",")
		db = a.DB
		db.Table(split[0]).
			Where("id IN (?)", value).
			Where("status = ?", split[1]).
			Count(&total)

		if total < 1 {
			return fmt.Errorf(fmt.Sprintf("value %v is not found.", value), field)
		}
		return nil
	})

	// bank_service_unique. unique service id per bank. format : unique_distinct:[search_table],[distinct_column],[unique_column],[limit]
	govalidator.AddCustomRule("unique_distinct", func(field string, rule string, message string, value interface{}) error {
		var (
			db    *gorm.DB
			total int
		)

		params := strings.Split(strings.TrimPrefix(rule, fmt.Sprintf("%s:", "unique_distinct")), ",")
		if len(params) == 4 {
			db = a.DB
			db.Table(params[0]).
				Where("? IN (SELECT DISTINCT ? FROM ?)", params[1], params[1], params[0]).
				Where("? IN (?)", params[2], value).
				Count(&total)

			limit, _ := strconv.Atoi(params[3])
			if total > limit {
				return fmt.Errorf(fmt.Sprintf("value %v already used.", value), field)
			}
		}

		return nil
	})

	// active / inactive string only.
	govalidator.AddCustomRule("active_inactive", func(field string, rule string, message string, value interface{}) error {
		val := value.(string)
		if strings.ToLower(val) != "active" && strings.ToLower(val) != "inactive" {
			return fmt.Errorf("The %s field must be contain word: active or inactive", field)
		}
		return nil
	})

	// validator for pagination
	govalidator.AddCustomRule("asc_desc", func(field string, rule string, message string, value interface{}) error {
		val := value.(string)
		if strings.ToUpper(val) != "ASC" && strings.ToUpper(val) != "DESC" {
			return fmt.Errorf("The %s field must be contain word: asc or desc", field)
		}
		return nil
	})

	// validator for loans
	govalidator.AddCustomRule("loan_statuses", func(field string, rule string, message string, value interface{}) error {
		val := value.(string)
		if val != "approved" && val != "rejected" && val != "processing" {
			return fmt.Errorf("The %s field must be contain either: approved, rejected, or processing", field)
		}
		return nil
	})

	// validator for otp entity types
	govalidator.AddCustomRule("otp_entity_types", func(field string, rule string, message string, value interface{}) error {
		val := value.(string)
		if val != "loan" && val != "borrower" {
			return fmt.Errorf("The %s field must be contain either: loan or borrower", field)
		}
		return nil
	})

	// validator for indonesia phone number
	govalidator.AddCustomRule("id_phonenumber", func(field string, rule string, message string, value interface{}) error {
		reg := regexp.MustCompile(`\+?([ -]?\d+)+|\(\d+\)([ -]\d+)`)
		if value == nil {
			return fmt.Errorf("no value")
		}
		val := value.(string)
		if !reg.MatchString(val) {
			return fmt.Errorf("The %s field is not a valid indonesia phone number", field)
		}
		return nil
	})

	// validator agent category
	govalidator.AddCustomRule("agent_categories", func(field string, rule string, message string, value interface{}) error {
		if value != nil {
			val := value.(string)
			if val != "agent" && val != "account_executive" {
				return fmt.Errorf("The %s field must be contain either: agent or account_executive", field)
			}
		}
		return nil
	})

	// must empty
	govalidator.AddCustomRule("unrequired", func(field, rule, message string, value interface{}) error {
		err := fmt.Errorf("The %s field is prohibited", field)
		if value != nil {
			return err
		}
		if _, ok := value.(multipart.File); ok {
			return nil
		}
		rv := reflect.ValueOf(value)
		switch rv.Kind() {
		case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
			if rv.Len() != 0 {
				return err
			}
		case reflect.Int:
			if !isEmpty(value.(int)) {
				return err
			}
		case reflect.Int8:
			if !isEmpty(value.(int8)) {
				return err
			}
		case reflect.Int16:
			if !isEmpty(value.(int16)) {
				return err
			}
		case reflect.Int32:
			if !isEmpty(value.(int32)) {
				return err
			}
		case reflect.Int64:
			if !isEmpty(value.(int64)) {
				return err
			}
		case reflect.Float32:
			if !isEmpty(value.(float32)) {
				return err
			}
		case reflect.Float64:
			if !isEmpty(value.(float64)) {
				return err
			}
		case reflect.Uint:
			if !isEmpty(value.(uint)) {
				return err
			}
		case reflect.Uint8:
			if !isEmpty(value.(uint8)) {
				return err
			}
		case reflect.Uint16:
			if !isEmpty(value.(uint16)) {
				return err
			}
		case reflect.Uint32:
			if !isEmpty(value.(uint32)) {
				return err
			}
		case reflect.Uint64:
			if !isEmpty(value.(uint64)) {
				return err
			}
		case reflect.Uintptr:
			if !isEmpty(value.(uintptr)) {
				return err
			}
		case reflect.Struct:
			switch rv.Type().String() {
			case "govalidator.Int":
				if v, ok := value.(govalidator.Int); ok {
					if v.IsSet {
						return err
					}
				}
			case "govalidator.Int64":
				if v, ok := value.(govalidator.Int64); ok {
					if v.IsSet {
						return err
					}
				}
			case "govalidator.Float32":
				if v, ok := value.(govalidator.Float32); ok {
					if v.IsSet {
						return err
					}
				}
			case "govalidator.Float64":
				if v, ok := value.(govalidator.Float64); ok {
					if v.IsSet {
						return err
					}
				}
			case "govalidator.Bool":
				if v, ok := value.(govalidator.Bool); ok {
					if v.IsSet {
						return err
					}
				}
			default:
				panic("govalidator: invalid custom type for required rule")

			}

		default:
			panic("govalidator: invalid type for required rule")

		}
		return nil
	})

	// validator for intention purpose
	govalidator.AddCustomRule("loan_purposes", func(field string, rule string, message string, value interface{}) error {
		var (
			queryRow *gorm.DB
			total    int
		)

		query := `SELECT COUNT(*) FROM loan_purposes WHERE name = ? AND status = ?`

		queryRow = a.DB.Raw(query, value, "active")

		queryRow.Row().Scan(&total)

		if total < 1 {
			if message != "" {
				return errors.New(message)
			}

			return fmt.Errorf("The %s doesn't match any loan purposes", field)
		}

		return nil
	})

	// validator loan purpose status
	govalidator.AddCustomRule("loan_purpose_status", func(field string, rule string, message string, value interface{}) error {
		val := value.(string)
		if val != "active" && val != "inactive" {
			return fmt.Errorf("The %s field must be contain either: active or inactive", field)
		}
		return nil
	})

	// active / inactive string only.
	govalidator.AddCustomRule("loan_payment_status", func(field string, rule string, message string, value interface{}) error {
		val := value.(string)
		if strings.ToLower(val) != "terbayar" && strings.ToLower(val) != "gagal_bayar" {
			return fmt.Errorf("The %s field must be contain word: terbayar or gagal_bayar", field)
		}
		return nil
	})
}

// isEmpty check a type is Zero
func isEmpty(x interface{}) bool {
	rt := reflect.TypeOf(x)
	if rt == nil {
		return true
	}
	rv := reflect.ValueOf(x)
	switch rv.Kind() {
	case reflect.Array, reflect.Map, reflect.Slice:
		return rv.Len() == 0
	}
	return reflect.DeepEqual(x, reflect.Zero(rt).Interface())
}
