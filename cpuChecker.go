package main

import (
	"bytes"
	"fmt"
	"os/exec"
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

	c1 := exec.Command("./cpu.bash")

	var b2 bytes.Buffer
	c1.Stdout = &b2

	c1.Start()
	c1.Wait()

	str := b2.String()
	trimmed := strings.TrimSpace(str)
	retVal, err := strconv.ParseFloat(trimmed, 64)
	return retVal, err
}
