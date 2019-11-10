package validator

import (
	"errors"
	"github.com/apptut/validator/rules"
	"reflect"
	"strings"
)

// 内置验证器
var validateMap = map[string]interface{}{
	"Required": rules.Required,
	"Min":      rules.Min,
	"Max":      rules.Max,
	"Regex":    rules.Regex,
	"Int":      rules.Int,
	"Numeric":  rules.Numeric,
	"Nullable": rules.Nullable,
	"Email":    rules.Email,
	"Url":      rules.Url,
	"Mobile":   rules.Mobile,
	"In":       rules.In,
	"Lt":       rules.Lt,
	"Lte":      rules.Lte,
	"Gt":       rules.Gt,
	"Gte":      rules.Gte,
}

// 单个验证字段错误提示
type ValidError struct {
	Field  string
	Errors map[string]string
}

type CustomMsgElem map[string]string

type Validator struct {
	data      map[string][]string      // 需要验证的数据
	rules     map[string][]string      // 验证规则
	customMsg map[string]CustomMsgElem // 自定义错误

	ValidErrors []ValidError // 验证错误
}

/**
 * 不带自定义错误验证
 *
 * @param data map[string][]string 验证的值
 * @param rules map[string]string  验证规则
 * @return Validator, error 默认返回验证错误第一项
 */
func New(data map[string][]string, rules interface{}, args ...map[string]string) (*Validator, error) {
	message := make(map[string]string)
	if len(args) > 0 {
		message = args[0]
	}
	fmtRules := formatRules(rules)
	validator := Validator{data: data, rules: fmtRules}
	if ok := validator.missingCheck(data, fmtRules); !ok {
		// 获取错误的第一项作为返回值
		err := validator.ValidErrors[0]
		val, ok := err.Errors["def"]
		if !ok {
			val = "missing valid error"
		}
		return &validator, errors.New(val)
	}
	validator.parseMessage(message)

	return validator.run()
}

func formatRules(rules interface{}) map[string][]string {

	rulesType := reflect.TypeOf(rules).String()
	if rulesType != "map[string][]string" && rulesType != "map[string]string" {
		panic("the rules only support map[string][]string or map[string]string")
	}

	rulesVals := reflect.ValueOf(rules)
	var fmtRules = make(map[string][]string)
	keys := rulesVals.MapKeys()
	for _, key := range keys {
		rulesItem := rulesVals.MapIndex(key)
		keyStr := key.Interface().(string)
		if rulesItem.Kind() == reflect.String {
			itemStr := strings.Split(rulesItem.Interface().(string), "|")
			fmtRules[keyStr] = itemStr
		} else {
			fmtRules[keyStr] = rulesItem.Interface().([]string)
		}
	}
	return fmtRules
}

func (this *Validator) run() (*Validator, error) {
	for key, item := range this.rules {
		this.parse(key, item)
	}

	if this.ValidErrors != nil || len(this.ValidErrors) > 0 {
		err := this.ValidErrors[0]
		val, ok := err.Errors["def"]
		if !ok {
			for _, item := range err.Errors {
				return this, errors.New(item)
			}
		}

		return this, errors.New(val)
	}

	return this, nil
}

func (this *Validator) parse(key string, rules []string) {
	for _, rule := range rules {
		flagIndex := strings.Split(rule, ":")
		param := ""
		ruleName := rule
		if len(flagIndex) > 1 {
			ruleName = flagIndex[0]
			param = flagIndex[1]
		}

		if _, ok := validateMap[ucfirst(ruleName)]; !ok {
			panic(ruleName + "the valid rule not exist")
		}

		if this.isVerifiable(key, rules) {
			dynamicFunc := reflect.ValueOf(validateMap[ucfirst(ruleName)])
			if dynamicFunc.IsValid() {
				value := this.data[key]
				arguments := make([]reflect.Value, 2) // 传递2个固定参数
				arguments[0] = reflect.ValueOf(value)
				arguments[1] = reflect.ValueOf(param)
				result := dynamicFunc.Call(arguments)
				ok := result[0].Interface().(bool);
				if !ok {
					this.addErrors(key, ruleName, value)
				}
			}
		}
	}
}

/**
 * 处理错误数据
 *
 * @param key
 * @param rule
 */
func (this *Validator) addErrors(field string, rule string, value []string) {
	customMsg, exist := this.customMsg[field] // 获取是否对验证字段存在自定义错误提示
	if exist {
		// 检测是否存在默认值, 字段优先级高于其他优先级
		msg, ok := customMsg["def"]
		if ok {
			this.insertError("def", field, msg, rule)
		}
		// 检测是否存在具体匹配错误内容
		fieldMsg, fieldOk := customMsg[rule]
		if !fieldOk {
			this.notExistCustomInsert(field, rule)
		} else {
			key := rule
			this.insertError(key, field, fieldMsg, rule)
		}
	} else {
		this.notExistCustomInsert(field, rule)
	}
}

/**
 * 添加默认错误提示
 *
 * @param field {string} 需要验证的字段
 * @param rule {string} 验证规则
 */
func (this *Validator) notExistCustomInsert(field string, rule string) {
	msg := "the field " + field + " not valid in " + rule
	key := rule
	this.insertError(key, field, msg, rule)
}

/**
 * 验证不通过添加相应的错误提示
 *
 * @param field {string} 需要验证的字段
 * @param rule {string} 验证规则
 */
func (this *Validator) insertError(key string, field string, msg string, rule string) {
	if this.ValidErrors == nil {
		itemErrors := make(map[string]string)
		itemErrors[key] = msg
		validErrItem := ValidError{Field: field, Errors: itemErrors}
		this.ValidErrors = []ValidError{validErrItem}
	} else {
		index := this.existError(field)
		if index >= 0 {
			this.ValidErrors[index].Errors[key] = msg
		} else {
			itemErrors := make(map[string]string)
			itemErrors[key] = msg
			newValidErr := ValidError{Field: field, Errors: itemErrors}
			this.ValidErrors = append(this.ValidErrors, newValidErr)
		}
	}
}

/**
 * 获取错误数组中索引
 *
 * @param field
 * @param rule
 * @return int
 */
func (this *Validator) existError(field string) int {
	for key, item := range this.ValidErrors {
		if item.Field == field {
			return key
		}
	}
	return -1
}

/**
 * 检测是否需要验证
 *
 * @param key
 * @param rules
 * @return bool
 */
func (this *Validator) isVerifiable(key string, rules []string) bool {
	rule, ok := this.data[key]
	if inArray(rules, "nullable") {
		if !ok {
			return false
		} else if rule != nil {
			for _, ruleItem := range rule {

				// 如果发现其中某一项值不为空，则需要验证
				if len(ruleItem) > 0 {
					return true
				}
			}

			return false
		}
	}

	return true
}

func (this *Validator) parseMessage(message map[string]string) {
	if len(message) == 0 {
		return
	}
	for key, item := range message {
		if strings.Contains(key, ".") {
			itemArr := strings.Split(key, ".")
			field := itemArr[0]
			rule := itemArr[1]
			_, ok := this.data[field]
			if _, exist := validateMap[ucfirst(rule)]; exist && ok {
				this.addMessage(field, rule, item)
			}
		} else {
			_, ok := this.data[key]
			if ok {
				this.addMessage(key, "", item)
			}
		}
	}
}

/**
 * 添加自定义错误到提示集合中(ValidErrors)
 *
 * @param str
 * @return string
 */
func (this *Validator) addMessage(field string, rule string, message string) {
	newMsg := make(map[string]string)
	if len(rule) > 0 { // 带具体条件错误提示
		newMsg[rule] = message
	} else {
		newMsg["def"] = message
	}

	if this.customMsg == nil {
		this.customMsg = make(map[string]CustomMsgElem)
		this.customMsg[field] = newMsg
	} else {
		for key, item := range newMsg {
			this.customMsg[field][key] = item
		}
	}

}

/**
 * 字符串首字母大写转换
 *
 * @param str
 * @return string
 */
func ucfirst(str string) string {
	return strings.ToUpper(str[0:1]) + str[1:]
}

/**
 * 检测元素是否存在数组中
 *
 * @param arr
 * @param elem
 * @return bool
 */
func inArray(arr []string, elem string) (bool) {
	for _, val := range arr {
		if elem == val {
			return true
		}
	}
	return false
}

/**
 * 检测message 字段在验证数据中是否存在
 *
 * @param data map[string]string 验证的值
 * @param rules map[string]string 验证规则
 */
func (this *Validator) missingCheck(data map[string][]string, rules map[string][]string) bool {
	if len(rules) == 0 {
		panic("验证规则不存在")
		return false
	}
	for key, item := range rules {
		_, ok := data[key]
		if !inArray(item, "nullable") && !ok {
			msg := "the param " + key + " not valid!"
			this.insertError("def", key, msg, "no")
		}
	}

	return this.ValidErrors == nil || len(this.ValidErrors) <= 0
}
