package main

import (
	"os"
	"fmt"
	"exec"
	"flag"
	"time"
)


const (
	ANSI_ESCAPE = string(byte(27))
	SECOND      = 1e9
)

func clearScreen() {
	fmt.Printf("%s[2J%s[H", ANSI_ESCAPE, ANSI_ESCAPE)
}


func printHead(interval int64, arguments []string) {
	fmt.Printf("Every %ds: %s %40s\n", interval, arguments,
		time.LocalTime().String())
}


func every(interval int64, program string, arguments []string) {
	clearScreen()
	printHead(interval, arguments)
	proc, err := exec.Run(program, arguments, nil, "",
		exec.PassThrough,
		exec.PassThrough,
		exec.PassThrough)
	if err != nil {
		fmt.Printf("error spawning process: %s\n", err)
		os.Exit(1)
	}
	proc.Wait(os.WSTOPPED)
	time.Sleep(interval * SECOND)
}


func main() {
	var interval int64
	flag.Int64Var(&interval, "n", 1,
		"seconds to wait between updates")
	flag.Parse()
	program, _ := exec.LookPath(flag.Arg(0))
	var arguments = flag.Args()
	for { // loop forever
		every(interval, program, arguments)
	}

}
