package model

type FileMeta struct {
	UserSn   uint64 `json:"user_sn" db:"user_sn"`
	FileSha1 string `json:"file_sha1" db:"file_sha1"`
	FileName string `json:"file_name" db:"file_name"`
	FileSize int64  `json:"file_size" db:"file_size"`
	Location string `json:"file_addr" db:"file_addr"`
	UploadAt string `json:"upload_at" db:"upload_at"`
}
