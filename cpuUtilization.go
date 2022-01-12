package main

import (
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const FileProcStat = "/proc/stat"

// https://man7.org/linux/man-pages/man5/proc.5.html
type CpuUtilization struct {
	CpuName   string // The name of the cpu that was queried
	User      int    // Time spent in user mode
	Nice      int    // Time spent in user mode with low priority
	System    int    // Time spent in system mode
	Idle      int    // Time spent in the idle task
	IOWait    int    // Time waiting for I/O to complete
	IRQ       int    // Time servicing interrupts
	SoftIRQ   int    // Time servicing softirqs
	Steal     int    // Time spent in other operating systems when running in a virtualized environment
	Guest     int    // Time spent running a virtual CPU for guest operating systems under the control of the Linux kernel.
	GuestNice int    // Time spent running a niced guest
}

func (stat *CpuUtilization) BusyTime() int {
	return sum(
		stat.User,
		stat.Nice,
		stat.System,
		stat.IRQ,
		stat.SoftIRQ,
		stat.Steal,
		stat.Guest,
		stat.GuestNice,
	)
}

func (stat *CpuUtilization) IdleTime() int {
	return sum(
		stat.Idle,
		stat.IOWait,
	)
}

func (stat *CpuUtilization) TotalTime() int {
	return stat.BusyTime() + stat.IdleTime()
}

func (stat *CpuUtilization) TotalUtilization() float64 {
	return float64(stat.BusyTime()) / float64(stat.TotalTime())
}

func cpuTimes() (float64, float64, error) {
	stats, err := GetCpuStats()
	if err != nil {
		return 0.0, 0.0, err
	}

	busy := 0.0
	total := 0.0

	for _, stat := range stats {
		busy += float64(stat.BusyTime())
		total += float64(stat.TotalTime())
	}

	busy /= float64(len(stats))
	total /= float64(len(stats))

	return busy, total, nil
}

func CpuUtilizationDelta(deltaMillis int) (float64, error) {
	busyT0, totalT0, err := cpuTimes()
	if err != nil {
		return 0.0, err
	}

	time.Sleep(time.Duration(deltaMillis) * time.Millisecond)

	busyT1, totalT1, err := cpuTimes()
	if err != nil {
		return 0.0, err
	}

	return (busyT1 - busyT0) / (totalT1 - totalT0), nil
}

func cpuStatFromLine(line string) (*CpuUtilization, error) {
	procLine := strings.TrimSpace(line)
	lineSpacesCollapsed := reMultiSpace.ReplaceAllString(procLine, " ")
	lineSplit := strings.Split(lineSpacesCollapsed, " ")

	userTime, err := strconv.Atoi(lineSplit[1])
	if err != nil {
		return nil, err
	}

	niceTime, err := strconv.Atoi(lineSplit[2])
	if err != nil {
		return nil, err
	}

	systemTime, err := strconv.Atoi(lineSplit[3])
	if err != nil {
		return nil, err
	}

	idleTime, err := strconv.Atoi(lineSplit[4])
	if err != nil {
		return nil, err
	}

	ioWaitTime, err := strconv.Atoi(lineSplit[5])
	if err != nil {
		return nil, err
	}

	irqTime, err := strconv.Atoi(lineSplit[6])
	if err != nil {
		return nil, err
	}

	softIrqTime, err := strconv.Atoi(lineSplit[7])
	if err != nil {
		return nil, err
	}

	stealTime, err := strconv.Atoi(lineSplit[8])
	if err != nil {
		return nil, err
	}

	guestTime, err := strconv.Atoi(lineSplit[9])
	if err != nil {
		return nil, err
	}

	guestNiceTime, err := strconv.Atoi(lineSplit[10])
	if err != nil {
		return nil, err
	}

	return &CpuUtilization{
		CpuName:   lineSplit[0],
		User:      userTime,
		Nice:      niceTime,
		System:    systemTime,
		Idle:      idleTime,
		IOWait:    ioWaitTime,
		IRQ:       irqTime,
		SoftIRQ:   softIrqTime,
		Steal:     stealTime,
		Guest:     guestTime,
		GuestNice: guestNiceTime,
	}, nil
}

func GetCpuStats() ([]*CpuUtilization, error) {
	reCpuLine := regexp.MustCompile(`^cpu\d`)

	procBytes, err := os.ReadFile(FileProcStat)
	if err != nil {
		return nil, err
	}

	allCpuStat := make([]*CpuUtilization, 0)
	procStr := strings.TrimSpace(string(procBytes))

	for _, line := range strings.Split(procStr, "\n") {
		if reCpuLine.MatchString(line) {
			cpuStat, err := cpuStatFromLine(line)
			if err != nil {
				return nil, err
			}
			allCpuStat = append(allCpuStat, cpuStat)
		}
	}

	return allCpuStat, nil
}
