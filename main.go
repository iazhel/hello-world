package main

import (
	"errors"
	"fmt"
	"net"
	"os/exec"
	"regexp"
	"runtime"
	"strconv"
	"strings"
	//	ac "agentclient"
	//	fm "filemanager"
)

// Application config structure
type HardwareConfig struct {
	AgentVer      string  `json:"agent_version"`       // returned by agent
	FirmwareVer   string  `json:"firmware_version"`    // returned by agent
	VSSDCount     int     `json:"number_of_vssd"`      // returned by agent
	OsName        string  `json:"os_name"`             // returned by agent
	OsVer         string  `json:"os_version"`          // returned by agent
	ZuariBuildVer string  `json:"zuari_build_version"` // returned by agent
	Kernel        string  `json:"kernel"`              // returned by agent
	MAC           string  `json:"mac_address"`         // returned by agent
	CPU           string  `json:"cpu_model"`           // returned by agent
	CPUCores      int     `json:"cpu_cores"`           // returned by agent
	RAM           float64 `json:"ram"`                 // returned by agent
	DriveSize     int     `json:"drives_size"`         // returned by agent
}

func main() {
	hw := GetHardwareConfig()
	fmt.Printf("%#v", hw)
}

func GetHardwareConfig() HardwareConfig {
	sysConfig := HardwareConfig{
		CPUCores: runtime.NumCPU(),
	}

	switch sysConfig.OsName {
	//		case "WIN64":
	case "LINUX64":
		// get OS version
		out, err := exec.Command("hostnamectl").Output()
		line := ScanOn(out, "Operating System")
		if err == nil {
			OsVer, _ := ParseAfter(string(line), ":")
			sysConfig.OsVer = OsVer
		}
		// get Kernel
		out, err = exec.Command("uname", "-r").Output()
		if err == nil {
			sysConfig.Kernel = strings.TrimSpace(string(out))
		}
		// get CPU model
		out, err = exec.Command("cat", "/proc/cpuinfo").Output()
		if err == nil {
			cpuModel := ScanOn(out, "model name")
			cpuModel, _ = ParseAfter(cpuModel, ":")
			sysConfig.CPU = cpuModel
		}

		// get total RAM, GB
		out, err = exec.Command("cat", "/proc/meminfo").Output()
		if err == nil {
			line := ScanOn(out, "MemTotal")
			// find all numbers
			re := regexp.MustCompile("[0-9]+")
			numbers := (re.FindAllString(line, -1))
			ram, _ := strconv.ParseFloat(numbers[0], 64)
			sysConfig.RAM = float64(int(0.5 + ram/(1024*1024)))
		}
		// get zuari build version
		out, err = exec.Command("cat", "/opt/stellus/stellus-release").Output()
		if err == nil {
			sysConfig.ZuariBuildVer = strings.TrimSpace(string(out))
		}

		// get SSD count
		var VSSDDevices []string
		out, err = exec.Command("ls", "/dev/").Output()
		if err == nil {
			re := regexp.MustCompile("nvme[0-9][a-z].")
			VSSDDevices = (re.FindAllString(string(out), -1))
			sysConfig.VSSDCount = len(VSSDDevices)
		}
		// find firmvare version
		if sysConfig.VSSDCount != 0 {
			var slotValue string
			devName := VSSDDevices[0]
			out, err = exec.Command("/mnt/filer/zuari/tools/vssd_tools/nvmeredrive", "-GL", "-d", "/dev/"+devName, "-f").Output()
			if err == nil {
				activeSlot := ScanOn(out, "Active")
				slotValue, _ = ParseAfter(activeSlot, ":")
				FWLine := ScanOn(out, "Slot "+slotValue)
				FWVer, _ := ParseAfter(FWLine, ":")
				sysConfig.FirmwareVer = FWVer
			}

			// determine drivers capasity
			out, err := exec.Command("/opt/stellus/svm/hms/output/vssdutil/debug/vssdutil", "-i", slotValue).Output()
			if err != nil {
				line := ScanOn(out, "Virtual Block Count =")
				sizeStr, _ := ParseAfter(line, "=")
				size, _ := strconv.Atoi(sizeStr)
				sysConfig.DriveSize = int(size / 750)
			}

		}

		//  get default mac (hardware) addres
		out, err = exec.Command("route").Output()
		if err == nil {
			defRoute := ScanOn(out, "default")
			interfaces, _ := net.Interfaces()
			for _, inter := range interfaces {
				if strings.Contains(defRoute, inter.Name) {
					sysConfig.MAC = inter.HardwareAddr.String()
				}
			}
		}

	}
	return sysConfig
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
