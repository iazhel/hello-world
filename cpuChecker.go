package main

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"strings"
)

type Process struct {
	pid int
	cpu float64
}

func main() {
	c, err := CpuChecker()
	fmt.Println(c)
	fmt.Println(err)
}

func CpuChecker() (float64, error) {

	pwd, err := os.Getwd()
	if err != nil {
		fmt.Println("error path:", err)
		return 0, err
	}

	dir, _ := path.Split(os.Args[0])
	filePath := path.Join(pwd, dir, "cpu.bash")
	fmt.Println("CPU script filePath:", filePath)

	c1 := exec.Command(filePath)
	var b2 bytes.Buffer
	c1.Stdout = &b2

	c1.Start()
	c1.Wait()

	str := b2.String()
	trimmed := strings.TrimSpace(str)
	retVal, err := strconv.ParseFloat(trimmed, 64)
	return retVal, err
}
