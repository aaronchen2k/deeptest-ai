package _file

import (
	"github.com/mholt/archiver/v3"
	"path/filepath"
)

func Unzip(zipPath, dist string) (targetDir string, err error) {
	scriptFolder := GetZipSingleDir(zipPath)

	if scriptFolder != "" { // single dir in zip
		targetDir = filepath.Join(dist, scriptFolder)
		err = archiver.Unarchive(zipPath, targetDir)

	} else { // more then one dir, unzip to a folder
		fileNameWithoutExt := GetFileNameWithoutExt(zipPath)
		targetDir = filepath.Join(dist, fileNameWithoutExt) + PathSep
		err = archiver.Unarchive(zipPath, targetDir)
	}

	return
}
