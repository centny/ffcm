@echo off
del /Q /S build
mkdir build
mkdir build\cswf.ffcm
msbuild io.vty.cswf.ffcm.sln /p:Configuration="Release" /p:Platform="x64" /t:clean /t:build
xcopy io.vty.cswf.ffcm.console\bin\x64\Release\cswf-ffcm.exe*  build\cswf.ffcm
xcopy io.vty.cswf.ffcm.console\bin\x64\Release\*.dll build\cswf.ffcm
go build -o build\cswf.ffcm\ffcm.exe github.com/Centny/ffcm/ffcm
xcopy *.properties build\cswf.ffcm
xcopy *.sh build\cswf.ffcm
cd build
zip -r cswf.ffcm.zip cswf.ffcm
if not "%1"=="" (
 echo Upload package to fvm server %1
 fvm -u %1 cswf.ffcm 0.0.1 cswf.ffcm.zip
)
cd ..\