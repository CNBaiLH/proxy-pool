/**
* @Author: Lanhai Bai
* @Date: 2021/8/27 10:22
* @Description:
 */
package utils

type Response struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func JSONData(data interface{}) Response {
	return Response{
		Code:    0,
		Message: "成功",
		Data:    data,
	}
}

func JSONError(code int, message string) Response {
	return Response{
		Code:    code,
		Message: message,
		Data:    nil,
	}
}

type responseOption struct {
	field string
}

func WithResponseField(field string) func(o *responseOption) {
	return func(o *responseOption) {
		if field == "" {
			field = "list"
		}
		o.field = field
	}
}
func JSONDataWithPage(data interface{}, page *Pager, opts ...func(o *responseOption)) Response {
	opt := &responseOption{field: "list"}

	for _, fn := range opts {
		fn(opt)
	}
	return Response{
		Code:    0,
		Message: "成功",
		Data:    map[string]interface{}{opt.field: data, "meta": page},
	}
}
