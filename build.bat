@echo on

:: 设置环境变量
set GOOS=android
set GOARCH=arm64
set CGO_ENABLED=1

:: 检查 NDK_PATH 环境变量是否已设置
if "%NDK_PATH%"=="" (
    echo Error: NDK_PATH environment variable is not set. Please set it to the NDK home directory.
    exit /b 1
)

:: 构建 CC 路径
set CC=%NDK_PATH%\toolchains\llvm\prebuilt\windows-x86_64\bin\aarch64-linux-android21-clang

:: 编译 Go 项目
echo Building for Android...
go build -o injector-go main.go