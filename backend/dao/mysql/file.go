package dao

import (
	"go.uber.org/zap"
	"yiran/model"
)

type fileDAO struct{}

var DefaultFile = fileDAO{}

// 文件上传文成，保存 meta 信息到数据库
func (f fileDAO) Create(file *model.FileMeta) (err error) {
	// 涉及处理多表，需要使用事务
	tx, err := db.Begin()
	if err != nil {
		zap.L().Error("mysql transaction failed", zap.Error(err))
		return
	}
	defer ClearTransaction(tx) // 如果出现异常情况,导致没有 commit和rollback,可以用来收尾

	// 1. 更新文件表记录
	sqlStr := `INSERT INTO bt_file (file_sha1, file_name, file_size, file_addr, status) VALUES (?, ?, ?, ?, 1)`
	_, err = tx.Exec(sqlStr, file.FileSha1, file.FileName, file.FileSize, file.Location)
	if err != nil {
		zap.L().Error("fileDAO.Create failed", zap.Any("file", file), zap.Error(err))
		return
	}
	// 2. 更新用户文件表
	if err = DefaultUser.SaveFile(file); err != nil {
		return
	}

	// 事务提交
	if errTx := tx.Commit(); errTx != nil {
		zap.L().Error("commit article failed", zap.Any("file", file), zap.Error(err))
		return errTx
	}
	return
}

// 从数据库获取文件元信息
func (f fileDAO) GetFileMeta(fileHash string) (fileMeta *model.FileMeta, err error) {
	fileMeta = new(model.FileMeta)
	sqlStr := `SELECT file_sha1, file_name, file_size, file_addr from bt_file where file_sha1=? and status=1 limit 1`
	err = db.Get(fileMeta, sqlStr, fileHash)
	if err != nil {
		zap.L().Error("get file failed", zap.Error(err))
		return
	}
	return
}

// 从数据库获批量获取文件元信息
func (f fileDAO) GetFileMetaList(offset, limit int) (fileMetaList []*model.FileMeta, err error) {
	sqlStr := `SELECT file_sha1, file_name, file_size, file_addr from bt_file where status=1 order by id desc limit ?, ?`
	err = db.Select(&fileMetaList, sqlStr, offset, limit)
	if err != nil {
		zap.L().Error("get fileMeta list failed", zap.Error(err))
		return
	}
	return
}

// 更新文件的存储地址(如文件被转移了)
func (f fileDAO) UpdateFileLocation(fileHash, fileAddr string) (err error) {
	sqlStr := `UPDATE bt_file SET file_addr = ? where file_sha1 = ? limit 1`
	_, err = db.Exec(sqlStr, fileAddr, fileHash)
	if err != nil {
		zap.L().Error("update fileMeta location failed", zap.Error(err))
		return
	}
	return
}
