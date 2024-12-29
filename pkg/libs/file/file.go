package _file

import (
	"archive/zip"
	"bytes"
	"errors"
	"fmt"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	_consts "github.com/deeptest-com/deeptest-next/pkg/libs/consts"
	"github.com/mholt/archiver/v3"
	"github.com/oklog/ulid/v2"
	"github.com/snowlyg/helper/str"
	"math/rand"
	"os"
	"os/exec"
	"os/user"
	"path/filepath"
	"runtime"
	"strings"
	"time"
)

var (
	PathSep = string(os.PathSeparator)
)

// GetExecDir 当前执行目录
func GetExecDir() (ret string) {
	if IsDebug() {
		ret, _ = os.Getwd()
	} else {
		exePath, _ := os.Executable()
		ret = filepath.Dir(exePath)
	}

	ret = AddSepIfNeeded(ret)

	return
}

// GetWorkDir 当前工作目录
func GetWorkDir() (dir string) {
	if consts.WorkDir != "" {
		return consts.WorkDir
	}

	if IsDebug() {
		dir, _ = os.Getwd()
	} else {
		home, _ := GetUserHome()

		dir = filepath.Join(home, consts.System, consts.App)
		dir = AddSepIfNeeded(dir)
	}

	InsureDir(dir)
	consts.WorkDir = dir

	return
}

func IsDebug() bool {
	// is debug in ide or not
	exePath, _ := os.Executable()
	return strings.Index(strings.ToLower(exePath), "goland") > -1
}
func IsRelease() bool {
	return !IsDebug()
}

func GetUserHome() (dir string, err error) {
	user, err := user.Current()
	if err == nil {
		dir = user.HomeDir

	} else { // cross compile support

		if "windows" == runtime.GOOS { // windows
			dir, err = homeWindows()
		} else { // Unix-like system, so just assume Unix
			dir, err = homeUnix()
		}
	}

	dir = AddSepIfNeeded(dir)

	return
}
func homeUnix() (string, error) {
	if home := os.Getenv("HOME"); home != "" {
		return home, nil
	}

	// If failed, try the shell
	var stdout bytes.Buffer
	cmd := exec.Command("sh", "-c", "eval echo ~$USER")
	cmd.Stdout = &stdout
	if err := cmd.Run(); err != nil {
		return "", err
	}

	result := strings.TrimSpace(stdout.String())
	if result == "" {
		return "", errors.New("blank output when reading home directory")
	}

	return result, nil
}

func homeWindows() (string, error) {
	drive := os.Getenv("HOMEDRIVE")
	path := os.Getenv("HOMEPATH")
	home := drive + path

	if drive == "" || path == "" {
		home = os.Getenv("USERPROFILE")
	}

	if home == "" {
		return "", errors.New("HOMEDRIVE, HOMEPATH, and USERPROFILE are blank")
	}

	return home, nil
}

func AddSepIfNeeded(pth string) string {
	if strings.LastIndex(pth, _consts.FilePthSep) < len(pth)-1 {
		pth += _consts.FilePthSep
	}
	return pth
}

func ReadFile(filePath string) string {
	buf := ReadFileBuf(filePath)
	str := string(buf)
	return str
}

func ReadFileBuf(filePath string) []byte {
	buf, err := os.ReadFile(filePath)
	if err != nil {
		return []byte(err.Error())
	}

	return buf
}

func GetUploadFileName(name string) (ret string, err error) {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	rand, _ := ulid.New(ms, entropy)

	ret = AddFileNamePostfix(name, strings.ToLower(rand.String()))

	return
}

func AddFileNamePostfix(name, postfix string) (ret string) {
	name = strings.TrimPrefix(name, "./")

	index := strings.LastIndex(name, ".")

	if index <= 0 {
		return
	}

	base := name[:index]
	ext := name[index+1:]

	ret = str.Join(base, "@", postfix, ".", ext)

	return
}

func GetZipSingleDir(path string) string {
	folder := ""
	z := archiver.Zip{}
	err := z.Walk(path, func(f archiver.File) error {
		if f.IsDir() {
			zfh, ok := f.Header.(zip.FileHeader)
			if ok {
				fmt.Println("file: ", zfh.Name)

				if folder == "" && zfh.Name != "__MACOSX" {
					folder = zfh.Name
				} else {
					if strings.Index(zfh.Name, folder) != 0 {
						return errors.New("found more than one folder")
					}
				}
			}
		}
		return nil
	})

	if err != nil {
		return ""
	}

	return folder
}

func GetFileName(pathOrUrl string) string {
	index := strings.LastIndex(pathOrUrl, PathSep)

	name := pathOrUrl[index+1:]
	return name
}

func GetFileNameWithoutExt(pathOrUrl string) string {
	name := GetFileName(pathOrUrl)
	index := strings.LastIndex(name, ".")
	return name[:index]
}

func GetExtName(pathOrUrl string) string {
	index := strings.LastIndex(pathOrUrl, ".")

	if index == -1 {
		return ""
	}
	return pathOrUrl[index:]
}
