echo "Start compiling..."

set version=.\versions\

@REM MacOS Amd64
set mac_amd=.\versions\macos-amd64
if exist %mac_amd% (
    echo ""
) else (
    md %mac_amd%
)
set GOOS=darwin
set GOARCH=amd64
go build -ldflags "-s -w" -o %mac_amd% .
echo "compilation for macos-amd64 is done."
set name=gvc-macos-amd64.tar.gz

if exist %version%%name% (
    del %version%%name%
)

tar -cvzf %version%%name% %mac_amd%

@REM MacOS Arm64
set mac_arm=.\versions\macos-arm64
if exist %mac_arm% (
    echo ""
) else (
    md %mac_arm%
)
set GOOS=darwin
set GOARCH=arm64
go build -ldflags "-s -w" -o %mac_arm% .
echo "compilation for macos-arm64 is done."
set name=gvc-macos-arm64.tar.gz

if exist %version%%name% (
    del %version%%name%
)

tar -cvzf %version%%name% %mac_arm%

@REM Linux Amd64
set lin_amd=.\versions\linux-amd64
if exist %lin_amd% (
    echo ""
) else (
    md %lin_amd%
)
set GOOS=linux
set GOARCH=amd64
go build -ldflags "-s -w" -o %lin_amd% .
echo "compilation for linux-amd64 is done."
set name=gvc-linux-amd64.tar.gz

if exist %version%%name% (
    del %version%%name%
)

tar -cvzf %version%%name% %lin_amd%

@REM Linux Arm64
set lin_arm=.\versions\linux-arm64
if exist %lin_arm% (
    echo ""
) else (
    md %lin_arm%
)
set GOOS=linux
set GOARCH=arm64
go build -ldflags "-s -w" -o %lin_arm% .
echo "compilation for linux-arm64 is done."
set name=gvc-linux-arm64.tar.gz

if exist %version%%name% (
    del %version%%name%
)

tar -cvzf %version%%name% %lin_arm%

@REM Windows Arm64
set win_arm=.\versions\windows-arm64
if exist %win_arm% (
    echo ""
) else (
    md %win_arm%
)
set GOOS=windows
set GOARCH=arm64
go build -ldflags "-s -w" -o %win_arm% .
echo "compilation for win-arm64 is done."
set name=gvc-windows-arm64.tar.gz

if exist %version%%name% (
    del %version%%name%
)

tar -cvzf %version%%name% %win_arm%

@REM Windows Amd64
set win_amd=.\versions\windows-amd64
if exist %win_amd% (
    echo ""
) else (
    md %win_amd%
)
set GOOS=windows
set GOARCH=amd64
go build -ldflags "-s -w" -o %win_amd% .
echo "compilation for win-amd64 is done."
set name=gvc-windows-amd64.tar.gz

if exist %version%%name% (
    del %version%%name%
)

tar -cvzf %version%%name% %win_amd%
