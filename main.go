package main

/*
#cgo LDFLAGS: -L${SRCDIR}/libs -Wl,-Bstatic -llinjector_rs -Wl,-Bdynamic -ldl
#include <stdlib.h>
int inject(int pid, const char* file_path, int injection_type, const char* func_sym, const char* var_sym, int debug, int logcat);
*/
import "C"
import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"unsafe"
)

func main() {
	// 定义命令行参数
	pid := flag.Int("p", 0, "pid of the target process")
	appPackageName := flag.String("a", "", "target application's package name, get pid by pidof and do injection")
	filePath := flag.String("f", "", "path of the library/shellcode to inject")
	injectionTypeStr := flag.String("i", "raw-dlopen", "type of injection. Possible values: raw-dlopen, memfd-dlopen, raw-shellcode")
	funcSym := flag.String("func-sym", "", "function to hijack for injection, in the form \"lib.so!symbol_name\"")
	varSym := flag.String("var-sym", "", "variable to hijack for injection, in the form \"lib.so!symbol_name\"")
	debug := flag.Bool("d", false, "enable debug logs")
	logcat := flag.Bool("logcat", false, "print logs to logcat")
	help := flag.Bool("h", false, "Print help (see a summary with '-h')")
	version := flag.Bool("V", false, "Print version")

	// 解析命令行参数
	flag.Parse()

	// 显示帮助信息
	if *help {
		fmt.Println("Usage: injector-go [options]")
		fmt.Println("Options:")
		flag.PrintDefaults()
		os.Exit(0)
	}

	// 显示版本信息
	if *version {
		fmt.Println("Version: 1.0")
		os.Exit(0)
	}

	// 如果没有提供 pid，尝试通过包名获取 pid
	if *pid == 0 && *appPackageName != "" {
		output, err := exec.Command("pidof", *appPackageName).Output()
		if err != nil {
			fmt.Printf("Failed to get PID for package %s: %v\n", *appPackageName, err)
			os.Exit(1)
		}
		pidStr := strings.TrimSpace(string(output))
		if pidStr == "" {
			fmt.Printf("No process found for package %s\n", *appPackageName)
			os.Exit(1)
		}
		_, err = fmt.Sscanf(pidStr, "%d", pid)
		if err != nil {
			fmt.Printf("Failed to parse PID: %v\n", err)
			os.Exit(1)
		}
	}
	// 检查必要参数
	if *pid == 0 || *filePath == "" {
		fmt.Println("Both PID and file path are required. Use -h for help.")
		os.Exit(1)
	}
	// 将注入类型字符串转换为整数
	injectionType := 0
	switch *injectionTypeStr {
	case "raw-dlopen":
		injectionType = 0
	case "memfd-dlopen":
		injectionType = 1
	case "raw-shellcode":
		injectionType = 2
	default:
		fmt.Printf("Invalid injection type: %s\n", *injectionTypeStr)
		os.Exit(1)
	}

	// 转换文件路径、函数符号、变量符号为 C 字符串
	cFilePath := C.CString(*filePath)
	defer C.free(unsafe.Pointer(cFilePath))

	var cFuncSym *C.char
	if *funcSym != "" {
		cFuncSym = C.CString(*funcSym)
		defer C.free(unsafe.Pointer(cFuncSym))
	}

	var cVarSym *C.char
	if *varSym != "" {
		cVarSym = C.CString(*varSym)
		defer C.free(unsafe.Pointer(cVarSym))
	}

	// 调用注入函数
	result := C.inject(
		C.int(*pid),
		cFilePath,
		C.int(injectionType),
		cFuncSym,
		cVarSym,
		boolToInt(*debug),
		boolToInt(*logcat),
	)

	if result != 0 {
		fmt.Println("Injection failed")
	} else {
		fmt.Println("Injection succeeded")
	}
}

// boolToInt 将布尔值转换为整数
func boolToInt(b bool) C.int {
	if b {
		return 1
	}
	return 0
}
