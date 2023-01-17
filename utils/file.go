package utils

import (
	"encoding/json"
	"golang.org/x/xerrors"
	"io/fs"
	"os"
	"path/filepath"
)

// DefaultCacheDir 设置缓存目录（cnnvd-list-update的目录）
func DefaultCacheDir() string {
	//根据用户的操作系统获取缓存目录，若无法获取缓存目录，则获取临时目录
	cacheDir, err := os.UserCacheDir()
	if err != nil {
		cacheDir = os.TempDir()
	}
	dir := filepath.Join(cacheDir, "cnnvd-list-update")
	return dir
}

// CNNVDListDir 获取cnnvd-list目录
func CNNVDListDir() string {
	return filepath.Join(DefaultCacheDir(), "cnnvd-list")
}

// Write 写入漏洞数据
func Write(filePath string, data interface{}) error {
	// 返回filePath的路径
	dir := filepath.Dir(filePath)
	// 创建路径
	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		return xerrors.Errorf("failed to create %s: %w", dir, err)
	}

	// 创建文件
	f, err := os.Create(filePath)
	if err != nil {
		return xerrors.Errorf("file create error: %w", err)
	}
	defer f.Close()

	// 将data 序列化
	b, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return xerrors.Errorf("JSON marshal error: %w", err)
	}

	// 写入文件
	_, err = f.Write(b)
	if err != nil {
		return xerrors.Errorf("file write error: %w", err)
	}
	return nil
}

func Read(filePath string) (any, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, xerrors.Errorf("fail to read :%w", err)
	}
	var result interface{}
	err = json.Unmarshal(data, result)
	if err != nil {
		return nil, xerrors.Errorf("fail to unmarshal:%w", err)
	}
	return result, nil
}

// PathExists 判断所给路径是否存在
func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

// IsDir 判断所给路径是否为文件夹
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

func GetFile(path string) ([]string, error) {
	var files []string
	err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
		if !IsDir(path) {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return files, nil
}
