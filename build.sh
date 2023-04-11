#!/bin/zsh
echo "Start compiling..."

version="./versions/"

win_amd="./versions/windows-amd64"
if [ ! -d $win_amd ];then
  mkdir -p $win_amd
fi
export GOOS="windows"
export GOARCH="amd64"
go build -ldflags "-s -w" -o $win_amd .
echo "compilation for win-amd64 is done."
name="gvc-windows-amd64.zip"
if [  -d "./versions/$name" ];then
  rm "./versions/$name"
fi
zip -rq $name $win_amd
mv $name $version

win_arm="./versions/windows-arm64"
if [ ! -d $win_arm ];then
  mkdir -p $win_arm
fi
export GOOS="windows"
export GOARCH="arm64"
go build -ldflags "-s -w" -o $win_arm .
echo "compilation for win-arm64 is done."
name="gvc-windows-arm64.zip"
if [  -d "./versions/$name" ];then
  rm "./versions/$name"
fi
zip -rq $name $win_arm
mv $name $version

linux_amd="./versions/linux-amd64"
if [ ! -d $linux_amd ];then
  mkdir -p $linux_amd
fi
export GOOS="linux"
export GOARCH="amd64"
go build -ldflags "-s -w" -o $linux_amd .
echo "compilation for linux-amd64 is done."
name="gvc-linux-amd64.zip"
if [  -d "./versions/$name" ];then
  rm "./versions/$name"
fi
zip -rq $name $linux_amd
mv $name $version

linux_arm="./versions/linux-arm64"
if [ ! -d $linux_arm ];then
  mkdir -p $linux_arm
fi
export GOOS="linux"
export GOARCH="arm64"
go build -ldflags "-s -w" -o $linux_arm .
echo "compilation for linux-arm64 is done."
name="gvc-linux-arm64.zip"
if [  -d "./versions/$name" ];then
  rm "./versions/$name"
fi
zip -rq $name $linux_arm
mv $name $version

mac_arm="./versions/macos-arm64"
if [ ! -d $mac_arm ];then
  mkdir -p $mac_arm
fi
export GOOS="darwin"
export GOARCH="arm64"
go build -ldflags "-s -w" -o $mac_arm .
echo "compilation for darwin-arm64 is done."
name="gvc-macos-arm64.zip"
if [  -d "./versions/$name" ];then
  rm "./versions/$name"
fi
zip -rq $name $mac_arm
mv $name $version

mac_amd="./versions/macos-amd64"
if [ ! -d $mac_amd ];then
  mkdir -p $mac_amd
fi
export GOOS="darwin"
export GOARCH="amd64"
go build -ldflags "-s -w" -o $mac_amd .
echo "compilation for darwin-amd64 is done."
name="gvc-macos-amd64.zip"
if [  -d "./versions/$name" ];then
  rm "./versions/$name"
fi
zip -rq $name $mac_amd
mv $name $version
