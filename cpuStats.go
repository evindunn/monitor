package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
)

const FileCpuInfo = "/proc/cpuinfo"

type CpuInfo struct {
	Model      string
	MHz        float64
	Cores      int
	CoreID     int
	PhysicalID int
}

func NewCpuInfo() ([]*CpuInfo, error) {
	reEmptyLine := regexp.MustCompile(`(?m)^\s*$`)

	cpuinfoBytes, err := os.ReadFile(FileCpuInfo)
	if err != nil {
		return nil, err
	}

	cpuInfoStr := strings.TrimSpace(string(cpuinfoBytes))
	cpus := reEmptyLine.Split(cpuInfoStr, -1)

	cpuInfos := make([]*CpuInfo, 0)

	for _, cpu := range cpus {
		cpuInfo := CpuInfo{}

		for _, line := range strings.Split(cpu, "\n") {
			line = strings.TrimSpace(line)
			parts := reKeyValDelim.Split(line, 2)

			if len(parts) < 2 {
				continue
			}

			key := parts[0]
			value := parts[1]

			switch key {
			case "model name":
				cpuInfo.Model = value
				break
			case "cpu MHz":
				mhz, err := strconv.ParseFloat(value, 64)
				if err != nil {
					return nil, err
				}
				cpuInfo.MHz = mhz
				break
			case "cpu cores":
				cores, err := strconv.Atoi(value)
				if err != nil {
					return nil, err
				}
				cpuInfo.Cores = cores
				break
			case "core id":
				coreID, err := strconv.Atoi(value)
				if err != nil {
					return nil, err
				}
				cpuInfo.CoreID = coreID
				break
			case "physical id":
				physID, err := strconv.Atoi(value)
				if err != nil {
					return nil, err
				}
				cpuInfo.PhysicalID = physID
				break
			default:
				break
			}
		}

		cpuInfos = append(cpuInfos, &cpuInfo)
	}

	return cpuInfos, nil
}
