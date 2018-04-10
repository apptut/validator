# validator-go
go form validator like php laravel framework validator
go版表单验证器，类似php框架laravel中的validator, 用于对web常用表单数据验证。

## 使用方法

### 1. 简单使用 （simple use）

```go
// 需要验证的数据, go语言net包获取的表单数据为map[string][]string类型，因此验证器保持一致，方便直接验证form表单数据
data := map[string][]string{
		"name" : {"banana"},
		"mobile" : {"13800138000"},
}
rules := map[string]string{
    "name" : "min:1|max:10",
    "mobile" : "mobile",
}
valid, err := validator.New(data, rules)
if err != nil {
    // err 仅仅返回验证器中第一条错误验证信息
    fmt.Println(err.Error())
    // ValidErrors() 可以获取到本次验证所有验证错误字段信息
    fmt.Println(valid.ValidErrors())
}
```

### 2. 自定义错误提示信息 (custom valid msg)

该验证器支持自定义错误信息，方便大家再具体场景自定义错误描述内容，只需要在使用时传入第三个参数即可，错误提示的格式为`map[string]string`类型：

```go
// 验证数据
data := map[string][]string{
		"name" : {"banana"},
		"mobile" : {"13800138000"},
}

// 验证规则
rules := map[string]string{
    "name" : "min:1|max:10",
    "mobile" : "mobile",
}

// 自定义错误提示
msg := map[string]string {
    "name" : "用户名格式不正确",
    "mobile" : "手机号格式不正确"
}

valid, err := validator.New(data, rules, msg)
if err != nil {
    // err 仅仅返回验证器中第一条错误验证信息
    fmt.Println(err.Error())
    // ValidErrors() 可以获取到本次验证所有验证错误字段信息
    fmt.Println(valid.ValidErrors())
}
```

错误消息自定可以针对单个验证规则，也可以针对整个需要验证的字段，例如：

```go
msg := map[string]string {
    "name.min"  : "用户名至少一个字符",
    "name.max"  : "用户名不超过10个字符长度"
    "name"      : "用户名不正确"
}
```

### 3. 内置可以验证规则如下（rules）

默认情况下，传入验证器的所有数据验证都是`required`类型数据，如果需要对某个字段做可选项验证，那么可以添加`nullable`验证，即：

```go
rules := map[string]string{
    "mobile" : "nullable|mobile", // 表示，如果验证数据中存在mobile且值不为空，那么就验证，否则跳过
}
```

| 名称             | 描述                                                                 |
|:---------------- |:---------------------------------------------------------------------|
| required         | 验证默认即required, 一般不需要配置                                    |
| min              | 验证字符串最小长度，支持多字节字符，例如中文                           |
| max              | 验证字符串最大长度，支持多字节字符，例如中文                           |
| regex            | 正则表达式验证，如果正则表达式中包含"|"符号，请参考正则验证试验注意部分 |
| int              | 验证数据是否为整数                                                   |
| numeric          | 验证数据是否为数字串                                                 |
| nullable         | 验证数据可选，如果验证数据不存在或为空值，则跳过后续验证               |
| email            | 验证数据是否为合法邮箱                                               |
| url              | 验证数据是否为合法url地址                                            |
| mobile           | 大陆11位手机号验证                                                   |

#### 3.1 正则验证规则使用注意

一般来说正则验证规则和其他验证规则类似，例如下面验证mobile字段为有效手机号的正则：

```go
rules := map[string]string{
    "mobile" : "min:1|regex:^1[0-9]{10}$",
}
```

但是，如果一个正则验证中包含`|`符号时，则比较特殊，因为验证器本身是基于`|`来切分验证规则的，如果正则表达式中包含改符号会引起验证器错误拆分验证规则问题，那么此时，我们需要把所有验证规则放在一个`slice`切片中即可：

```go
rules := map[string][]string{
    "mobile" : {"min:1", "regex:^1[0-9]{10}$"}
}
```



