# 交叉编译配置
export CGO_ENABLED=0
export GOOS=linux # darwin linux windows
export GOARCH=amd64 # amd64 arm64
go build .
