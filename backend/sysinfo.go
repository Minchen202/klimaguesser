package main

import (
	"bufio"
	"os"
	"strconv"
	"strings"
	"time"
)




func cpuPercent() float64 {
	idle0, total0, ok0 := readCPUSample()
	if !ok0 {
		return 0
	}
	time.Sleep(200 * time.Millisecond)
	idle1, total1, ok1 := readCPUSample()
	if !ok1 {
		return 0
	}

	idleDelta := idle1 - idle0
	totalDelta := total1 - total0
	if totalDelta <= 0 {
		return 0
	}
	usage := (1.0 - float64(idleDelta)/float64(totalDelta)) * 100
	if usage < 0 {
		usage = 0
	}
	return usage
}

func readCPUSample() (idle, total uint64, ok bool) {
	f, err := os.Open("/proc/stat")
	if err != nil {
		return 0, 0, false
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	if !scanner.Scan() {
		return 0, 0, false
	}
	fields := strings.Fields(scanner.Text())
	if len(fields) < 5 || fields[0] != "cpu" {
		return 0, 0, false
	}

	var sum uint64
	for _, v := range fields[1:] {
		n, err := strconv.ParseUint(v, 10, 64)
		if err != nil {
			continue
		}
		sum += n
	}
	idleTime, _ := strconv.ParseUint(fields[4], 10, 64)
	return idleTime, sum, true
}


func ramPercent() float64 {
	f, err := os.Open("/proc/meminfo")
	if err != nil {
		return 0
	}
	defer f.Close()

	values := map[string]uint64{}
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		parts := strings.Fields(line)
		if len(parts) < 2 {
			continue
		}
		key := strings.TrimSuffix(parts[0], ":")
		n, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			continue
		}
		values[key] = n
	}

	total := values["MemTotal"]
	available := values["MemAvailable"]
	if total == 0 {
		return 0
	}
	used := total - available
	return float64(used) / float64(total) * 100
}
