package utils

import "fmt"

func PrintRed(str string) {
	fmt.Printf("\033[31m%s\033[0m", str)
}

func PrintGreen(str string) {
	fmt.Printf("\033[32m%s\033[0m", str)
}

func PrintBlue(str string) {
	fmt.Printf("\033[34m%s\033[0m", str)
}