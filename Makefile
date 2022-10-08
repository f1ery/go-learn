.PHONY: protoc
# 新版本生成
protoc:
	protoc -I=./proto --go_out=plugins=grpc:. ./proto/hello.proto
	#make inject-tag

# 执行不成功，安装go get github.com/favadi/protoc-go-inject-tag
inject-tag:
	ls ./*.pb.go |xargs -n1 -I {} protoc-go-inject-tag -input={}