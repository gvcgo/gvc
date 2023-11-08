package confs

import (
	"os"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
)

type ProtobufConf struct {
	GithubUrls      map[string]string `koanf:"github_urls"`
	ProtoGenGoUrl   string            `koanf:"proto_gen_go_url"`
	ProtoGenGRPCUrl string            `koanf:"proto_gen_grpc_url"`
	path            string
}

func NewProtobuf() (r *ProtobufConf) {
	r = &ProtobufConf{
		path: ProtobufDir,
	}
	r.setup()
	return
}

func (that *ProtobufConf) setup() {
	if ok, _ := utils.PathIsExist(that.path); !ok {
		if err := os.MkdirAll(that.path, os.ModePerm); err != nil {
			gprint.PrintError("%+v", err)
		}
	}
}

func (that *ProtobufConf) Reset() {
	that.GithubUrls = map[string]string{
		utils.Windows: "https://github.com/protocolbuffers/protobuf/releases/download/v25.0/protoc-25.0-win64.zip",
		"linux_amd64": "https://github.com/protocolbuffers/protobuf/releases/download/v25.0/protoc-25.0-linux-x86_64.zip",
		"linux_arm64": "https://github.com/protocolbuffers/protobuf/releases/download/v25.0/protoc-25.0-linux-aarch_64.zip",
		utils.MacOS:   "https://github.com/protocolbuffers/protobuf/releases/download/v25.0/protoc-25.0-osx-universal_binary.zip",
	}
	that.ProtoGenGoUrl = "google.golang.org/protobuf/cmd/protoc-gen-go@latest"
	that.ProtoGenGRPCUrl = "google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
}
