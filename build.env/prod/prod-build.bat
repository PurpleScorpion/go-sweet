@echo off
set CUR_PATH=%cd%
set VER=latest

cd /d %CUR_PATH%

docker login -u xxxxx -p xxxxx registry.cn-beijing.aliyuncs.com
if %errorlevel% NEQ 0 GOTO ERROR

docker build -f .\prod.dockerfile -t registry.cn-beijing.aliyuncs.com/demo/sweet-go:%VER% .
if %errorlevel% NEQ 0 GOTO ERROR

docker push registry.cn-beijing.aliyuncs.com/demo/sweet-go:%VER%
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