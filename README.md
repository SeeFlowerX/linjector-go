# Ref

- https://github.com/erfur/linjector-rs

# Build

build libs/liblinjector_rs.a from here

- https://github.com/SeeFlowerX/linjector-rs

```bash
cargo ndk --target aarch64-linux-android build --release
```

then copy target/aarch64-linux-android/release/liblinjector_rs.a to current project libs/linjector_rs.a

```bat
set NDK_PATH=path\to\ndk
build.bat
```

# Usage

```bash
adb push injector-go /data/local/tmp/injector-go
adb shell chmod +x /data/local/tmp/injector-go
adb shell "su -c \"/data/local/tmp/injector-go -p `pidof com.smile.gifmaker` -i raw-dlopen -f /data/local/tmp/libeloader.so\""
adb shell "su -c \"/data/local/tmp/injector-go -a com.smile.gifmaker -i raw-dlopen -f /data/local/tmp/libeloader.so\""
```

# Help

```bash
Usage: injector-go [options]
Options:
  -V    Print version
  -a string
        target application's package name, get pid by pidof and do injection
  -d    enable debug logs
  -f string
        path of the library/shellcode to inject
  -func-sym string
        function to hijack for injection, in the form "lib.so!symbol_name"
  -h    Print help (see a summary with '-h')
  -i string
        type of injection. Possible values: raw-dlopen, memfd-dlopen, raw-shellcode (default "raw-dlopen")
  -logcat
        print logs to logcat
  -p int
        pid of the target process
  -var-sym string
        variable to hijack for injection, in the form "lib.so!symbol_name"
```