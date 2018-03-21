package core

type CommonResult struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

type SliceUploadResult struct {
	Context  string `json:"ctx"`
	Checksum string `json:"checksum"`
	CRC32    uint   `json:"crc32"`
	Offset   uint64 `json:"offset"`
}
