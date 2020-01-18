package main

import "fmt"

//Println 调试打印
func Println(a ...interface{}) {
	if g_bPrintLog {
		fmt.Println(a...)
	}
}

//Println 打印错误
func PrintlnWarning(a ...interface{}) {
	//if g_bPrintLog {
	s := fmt.Sprint(a...)
	fmt.Printf("%c[1;40;33m%s%c[0m\n", 0x1B, s, 0x1B)
	//}
}

//Printf 调试打印
func Printf(format string, a ...interface{}) {
	if g_bPrintLog {
		fmt.Printf(format, a...)
	}
}

//Printf 打印错误
func PrintfWarning(format string, a ...interface{}) {
	//if g_bPrintLog {
	s := fmt.Sprintf(format, a...)
	fmt.Printf("%c[1;40;33m%s%c[0m\n", 0x1B, s, 0x1B)
	//}
}
