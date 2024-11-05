package _file

import (
	"os"
	"path"
)

func InsureDir(fp string) error {
	if IsExist(fp) {
		return nil
	}
	return os.MkdirAll(fp, os.ModePerm)
}

func IsExist(fp string) bool {
	_, err := os.Stat(fp)
	return err == nil || os.IsExist(err)
}

func IsFile(fp string) bool {
	f, e := os.Stat(fp)
	if e != nil {
		return false
	}
	return !f.IsDir()
}

func Remove(filename string) error {
	if IsFile(filename) && IsExist(filename) {
		return os.Remove(filename)
	}
	return nil
}

func WriteBytes(filePath string, b []byte) (int, error) {
	os.MkdirAll(path.Dir(filePath), os.ModePerm)
	fw, err := os.Create(filePath)
	if err != nil {
		return 0, err
	}
	defer fw.Close()
	return fw.Write(b)
}
