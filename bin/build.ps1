# Set-ExecutionPolicy -Scope Process -ExecutionPolicy Bypass

$PROGRAM_NAME="rapidlog"
$VERSION="v0.1.2"
$env:GOARCH="amd64"

$env:GOOS="windows"; go build -o .\bin\${PROGRAM_NAME}_${VERSION}_${env:GOOS}_${env:GOARCH}.exe
$env:GOOS="linux"; go build -o .\bin\${PROGRAM_NAME}_${VERSION}_${env:GOOS}_${env:GOARCH}
$env:GOOS="darwin"; go build -o .\bin\${PROGRAM_NAME}_${VERSION}_${env:GOOS}_${env:GOARCH}
