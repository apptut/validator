package rules

import (
	"regexp"
	"unicode/utf8"
	"strconv"
	"unicode"
	"net/url"
)

/**
 * 验证是否必填字符串
 * @param value 需要验证的值
 * @param param 自定义参数
 * @return bool
 */
func Required(value []string, _ string) bool {
	if len(value) <= 0 || len(value[0]) <= 0 {
		return false
	}
	return true
}

/**
 * 正则表达式验证
 *
 * @param value 需要验证的值
 * @param param 自定义参数
 * @return bool
 */
func Regex(value []string, pattern string) bool{
	if m, _ := regexp.MatchString(pattern, value[0]); !m {
		return false
	}
	return true
}


/**
 * 字符串最大长度判断， 包括最大值本身
 *
 * @param value 需要验证的值
 * @param param 自定义参数
 * @return bool
 */
func Max(value []string, param string) bool{
	maxInt, err := strconv.Atoi(param)
	if err != nil {
		return false
	}
	// 获取字符串的真实长度
	valueLen := utf8.RuneCountInString(value[0])

	return valueLen <= maxInt
}


/**
 * 字符串最小长度判断， 包括最小值
 *
 * @param value 需要验证的值
 * @param param 自定义参数
 * @return bool
 */
func Min(value []string, param string) bool{
	minInt, err := strconv.Atoi(param)
	if err != nil {
		return false
	}
	// 获取字符串的真实长度
	valueLen := utf8.RuneCountInString(value[0])
	return valueLen >= minInt
}


/**
 * 判断传入职是否是整数
 *
 * @param value 需要验证的值
 * @param param 自定义参数
 * @return bool
 */
func Int(value []string, _ string) bool{
	if len(value) <= 0 || len(value[0]) <= 0 {
		return false
	}
	_, err := strconv.Atoi(value[0])
	if err != nil {
		return false
	}
	return true
}

/**
 * 验证字符串是否所有字符都是数字
 *
 * @param value 需要验证的值
 * @param param 自定义参数
 * @return bool
 */
func Numeric(value []string, _ string) bool{
	if len(value) <= 0 || len(value[0]) <= 0 {
		return false
	}
	for _, c := range value[0] {
		if !unicode.IsDigit(c) {
			return false
		}
	}
	return true
}

/**
 * 可选项不需要任何验证
 */
func Nullable(value []string, _ string) bool{
	return true
}



/**
 * 验证邮箱地址是否正确
 */
func Email(value []string, _ string) bool{
	pattern := "^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$"
	if m, _ := regexp.MatchString(pattern, value[0]); !m {
		return false
	}
	return true
}


/**
 * 检测当前数据是否是有效Url地址
 */
func Url(value []string, _ string) bool {
	if _, err := url.Parse(value[0]); err != nil {
		return false
	}
	return true
}

/**
 * 检测当前数据是否是有效手机号码
 * 仅支持大陆11位手机号，不支持座机号码
 */
func Mobile(value []string, _ string) bool{
	pattern := "^1[3|5|7|8|9][0-9]{9}$"
	if m, _ := regexp.MatchString(pattern, value[0]); !m {
		return false
	}
	return true
}




