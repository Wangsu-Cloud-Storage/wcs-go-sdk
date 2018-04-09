package core

// 文件上传的额外可选设置

type PutExtra struct {
	Params   map[string]string // 上传可选参数字典，参数名次以 x: 开头
	MimeType string            // 指定文件的 MimeType
	Deadline int               // 文件保存期限。超过保存天数文件自动删除,单位：天。例如：1、2、3……注：0表示尽快删除，上传文件时建议不配置为0
}

func NewDefaultPutExtra() (put_extra *PutExtra) {
	return &PutExtra{make(map[string]string, 0), "", -1}
}

func NewPutExtra(deadline int) (put_extra *PutExtra) {
	return &PutExtra{make(map[string]string, 0), "", deadline}
}
