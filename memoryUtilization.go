package main

import (
	"os"
	"strconv"
	"strings"
	"time"
)

const FileMemInfo = "/proc/meminfo"

type MemInfo struct {
	TotalKB int
	FreeKB  int
}

func NewMemInfo() (*MemInfo, error) {
	memInfoBytes, err := os.ReadFile(FileMemInfo)
	if err != nil {
		return nil, err
	}

	memInfoStr := strings.TrimSpace(string(memInfoBytes))
	memInfo := MemInfo{}

	for _, line := range strings.Split(memInfoStr, "\n") {
		line = strings.TrimSpace(line)
		parts := reKeyValDelim.Split(line, 2)

		if len(parts) < 2 {
			continue
		}

		key := parts[0]
		value := parts[1]

		switch key {
		case "MemTotal":
			valueParts := reMultiSpace.Split(value, 2)
			total, err := strconv.Atoi(valueParts[0])
			if err != nil {
				return nil, err
			}
			memInfo.TotalKB = total
			break
		case "MemFree":
			valueParts := reMultiSpace.Split(value, 2)
			free, err := strconv.Atoi(valueParts[0])
			if err != nil {
				return nil, err
			}
			memInfo.FreeKB = free
			break
		default:
			break
		}
	}

	return &memInfo, nil
}

func (m *MemInfo) UsedKB() int {
	return m.TotalKB - m.FreeKB
}

func (m *MemInfo) Utilization() float64 {
	return float64(m.UsedKB()) / float64(m.TotalKB)
}

func MemoryUtilizationDelta(deltaMillis float64) (float64, error) {
	mem0, err := NewMemInfo()
	if err != nil {
		return 0.0, err
	}

	time.Sleep(time.Duration(deltaMillis) * time.Millisecond)

	mem1, err := NewMemInfo()
	if err != nil {
		return 0.0, err
	}

	return float64(mem1.UsedKB()-mem0.UsedKB()) / float64(mem0.TotalKB), nil
}
