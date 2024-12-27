package _file

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/deeptest-com/deeptest-next/internal/pkg/consts"
	_consts "github.com/deeptest-com/deeptest-next/pkg/libs/consts"
	_logs "github.com/deeptest-com/deeptest-next/pkg/libs/log"
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
	name = strings.TrimPrefix(name, "./")

	index := strings.LastIndex(name, ".")

	if index <= 0 {
		msg := fmt.Sprintf("文件名错误 %s", name)
		_logs.Info(msg)
		err = errors.New(msg)
		return
	}

	base := name[:index]
	ext := name[index+1:]

	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	rand, _ := ulid.New(ms, entropy)

	ret = str.Join(base, "-", strings.ToLower(rand.String()), ".", ext)

	return
}
