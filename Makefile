ENV        ?= dp
include    env.$(ENV)
print_env:
	@echo $(PROJECT)@$(VERSION)

ifeq ($(OS),Windows_NT)
    PLATFORM="windows"
else
    ifeq ($(shell uname),Darwin)
        PLATFORM="mac"
    else
        PLATFORM="unix"
    endif
endif

ifeq ($(PLATFORM),"mac")
	QINIU_DIR=~/work/qiniu/
else
    QINIU_DIR=~/work/qiniu/
endif

QINIU_DIST_DIR=${QINIU_DIR}${PROJECT}/${VERSION}/

SERVER_MAIN_FILE=cmd/server/main.go

BIN_DIR=bin/

BUILD_TIME=`git show -s --format=%cd`
GO_VERSION=`go version`
GIT_HASH=`git show -s --format=%H`
BUILD_CMD_UNIX=go build -ldflags "-X 'main.AppVersion=${VERSION}' -X 'main.BuildTime=${BUILD_TIME}' -X 'main.GoVersion=${GO_VERSION}' -X 'main.GitHash=${GIT_HASH}'"
BUILD_CMD_WIN=go build -ldflags "-s -w -X 'main.AppVersion=${VERSION}' -X 'main.BuildTime=${BUILD_TIME}' -X 'main.GoVersion=${GO_VERSION}' -X 'main.GitHash=${GIT_HASH}'"

default: compile_ui_web win64 win32 linux mac

# 非客户端版本打包，更新客户端需先运行 make compile_ui_client
win64: prepare compile_server_win64 zip_web_win64
win32: prepare compile_server_win32 zip_web_win32
linux: prepare compile_server_linux zip_web_linux
mac:   prepare compile_server_mac   zip_web_mac

prepare: update_version
update_version: gen_version_file

gen_version_file:
	@echo 'gen version'
	@mkdir -p ${QINIU_DIR}/${PROJECT}/
	@echo '{"version": "${VERSION}"}' > ${QINIU_DIR}/${PROJECT}/version.json

compile_ui:
	@cd ui && pnpm build && cd ..

# server
compile_server_win64:
	@echo 'start compile win64'
	@CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 \
		${BUILD_CMD_WIN} -x -v \
		-o ${BIN_DIR}win64/${PROJECT}-server-next.exe ${SERVER_MAIN_FILE}

compile_server_win32:
	@echo 'start compile win32'
	@CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 \
		${BUILD_CMD_WIN} -x -v \
		-o ${BIN_DIR}win32/${PROJECT}-server-next.exe ${SERVER_MAIN_FILE}

compile_server_linux:
	@echo 'start compile linux'
ifeq ($(PLATFORM),"mac")
	@CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-gcc CXX=/usr/local/gcc-4.8.1-for-linux64/bin/x86_64-pc-linux-g++ \
		${BUILD_CMD_UNIX} \
		-o ${BIN_DIR}linux/${PROJECT}-server-next ${SERVER_MAIN_FILE}
else
	CGO_ENABLED=1 GOOS=linux GOARCH=amd64 CC=gcc CXX=g++ \
		${BUILD_CMD_UNIX} \
		-o ${BIN_DIR}linux/${PROJECT}-server-next ${SERVER_MAIN_FILE}
endif

compile_server_mac:
	@echo 'start compile darwin'
	@echo
	@CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 \
		${BUILD_CMD_UNIX} \
		-o ${BIN_DIR}darwin/${PROJECT}-server-next ${SERVER_MAIN_FILE}
