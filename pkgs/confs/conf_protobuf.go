package confs

import (
	"os"

	"github.com/moqsien/goutils/pkgs/gtea/gprint"
	"github.com/moqsien/gvc/pkgs/utils"
)

type ProtobufConf struct {
	ProtocUrl       string `koanf,json:"protoc_url"`
	ProtoGenGoUrl   string `koanf:"proto_gen_go_url"`
	ProtoGenGRPCUrl string `koanf:"proto_gen_grpc_url"`
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
	that.ProtocUrl = "https://github.com/protocolbuffers/protobuf/releases/latest/"
	that.ProtoGenGoUrl = "google.golang.org/protobuf/cmd/protoc-gen-go@latest"
	that.ProtoGenGRPCUrl = "google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest"
}
