@echo off
set CUR_PATH=%cd%
set VER=local

cd /d %CUR_PATH%
docker build -f .\local.dockerfile -t 192.168.2.1:5000/sweet-go:%VER% ..\..
if %errorlevel% NEQ 0 GOTO ERROR

@REM docker login -u xxxxx -p xxxxx host_addr
if %errorlevel% NEQ 0 GOTO ERROR

docker push 192.168.2.1:5000/sweet-go:%VER%
if %errorlevel% NEQ 0 GOTO ERROR else GOTO OK

:OK
ECHO ========================================
ECHO command success
ECHO ========================================
GOTO END

:ERROR
ECHO ========================================
ECHO command failed
ECHO ========================================
GOTO END

:END
cd /d %CUR_PATH%