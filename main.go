package main

import (
	"errors"
	"fmt"
	"os/exec"
	"regexp"
	"strings"
)

func main() {

	var VSSDDevices []string
	VSSDCount := 0
	out, err := exec.Command("ls", "/dev/").Output()
	fmt.Printf("ls OUT:%s\n", strings.Split(string(out), "\n"))
	fmt.Printf("ls Err:%s\n", err)
	if err == nil {
		re := regexp.MustCompile("nvme[0-9].")
		VSSDDevices = (re.FindAllString(string(out), -1))
		VSSDCount = len(VSSDDevices)
		fmt.Println("VSSDDevices:", VSSDDevices)
		fmt.Println("VSSDCount:", VSSDCount)
	}
	if VSSDCount != 0 {
		devName := VSSDDevices[0]
		// run script to find firmvare
		out, err = exec.Command("/mnt/filer/zuari/tools/vssd_tools/nvmeredrive", "-GL", "-d", "/dev/"+devName).Output()
		fmt.Printf("mnt OUT:%s\n", out)
		fmt.Printf("mnt Err:%v\n", err)
		if err == nil {
			activeSlot := ScanOn(out, "Active")
			fmt.Println("active slot:", activeSlot)
			slotValue, _ := ParseAfter(activeSlot, ":")
			fmt.Println("slot value:", slotValue)
			FWLine := ScanOn(out, "Slot "+slotValue)
			fmt.Println("Line :", FWLine)
			FWVer, _ := ParseAfter(FWLine, ":")
			fmt.Println("Firmware:", FWVer)
		}
	}

}

func ScanOn(out []byte, text string) string {
	lines := strings.Split(string(out), "\n")
	for _, line := range lines {
		if strings.Contains(line, text) {
			return strings.TrimSpace(line)
		}
	}
	return ""
}

func ParseAfter(out string, sep string) (string, error) {
	i := strings.IndexAny(out, ":") + 1
	if len(out) < i || i == 0 {
		return "", errors.New("Not possible to parse output")
	}
	stringOut := strings.TrimSpace(out[i:])
	return stringOut, nil

}
