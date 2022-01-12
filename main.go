package main

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

const CPUResolutionMillis = 1000
const MemResolutionMillis = 1000

func main() {
	cpuUtilization := lockedFloat{}
	memUtilization := lockedFloat{}

	go cpuUtilizationRoutine(&cpuUtilization)
	go memUtilizationRoutine(&memUtilization)

	router := gin.Default()

	router.GET("/cpu/info", getCpuInfo)
	router.GET("/cpu/utilization", func(c *gin.Context) {
		getCpuUtilization(&cpuUtilization, c)
	})

	router.GET("/memory/info", getMemInfo)
	router.GET("/memory/utilization", func(c *gin.Context) {
		getMemUtilization(&memUtilization, c)
	})

	router.Run()
}

func cpuUtilizationRoutine(val *lockedFloat) {
	for {
		cpuUtilization, err := CpuUtilizationDelta(CPUResolutionMillis)
		if err != nil {
			fmt.Println(err)
		}
		if cpuUtilization > 0 {
			val.Set(cpuUtilization)
		}
	}
}

func memUtilizationRoutine(val *lockedFloat) {
	for {
		memUtilization, err := MemoryUtilizationDelta(MemResolutionMillis)
		if err != nil {
			fmt.Println(err)
		}
		if memUtilization > 0 {
			val.Set(memUtilization)
		}
	}
}

func handleError(c *gin.Context, err error) {
	c.JSON(500, gin.H{
		"message": err.Error(),
	})
}

func getCpuInfo(c *gin.Context) {
	cpuInfo, err := NewCpuInfo()
	if err == nil {
		c.JSON(200, cpuInfo)
	} else {
		handleError(c, err)
	}
}

func getCpuUtilization(val *lockedFloat, c *gin.Context) {
	c.JSON(200, gin.H{
		"cpuUtilization": val.Get(),
	})
}

func getMemInfo(c *gin.Context) {
	memInfo, err := NewMemInfo()
	if err == nil {
		c.JSON(200, memInfo)
	} else {
		handleError(c, err)
	}
}

func getMemUtilization(val *lockedFloat, c *gin.Context) {
	c.JSON(200, gin.H{
		"memoryUtilization": val.Get(),
	})
}
