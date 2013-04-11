@echo off

setlocal

if exist make.bat goto ok
echo make.bat must be run within its container folder
goto end

:ok

set OLDGOPATH=%GOPATH%
set GOPATH=%~dp0

go install code.google.com/p/freetype-go/freetype
go install utility/process
go install image_sign

set GOPATH=%OLDGOPATH%

:end
echo finished.
