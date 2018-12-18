package main

import (
	"bufio"
	"os/exec"
	"strconv"
	"strings"
)

type mdInfo struct {
	Path  string
	Size  uint64
	Used  uint64
	Avail uint64
	Use   uint64
}

type nodeStat struct {
	Type     string
	Active   uint64
	Total    uint64
	Write    uint64
	Read     uint64
	Remove   uint64
	Flush    uint64
	AllWrite uint64
	AllRead  uint64
}

func getMdInfo() ([]*mdInfo, error) {
	args := []string{
		"node", "md", " info", "-r",
	}

	output, err := dog(args)
	if err != nil {
		return nil, err
	}

	var i []*mdInfo
	for _, line := range output {
		a := strings.Fields(line)
		size, _ := strconv.ParseUint(a[2], 10, 64)
		used, _ := strconv.ParseUint(a[3], 10, 64)
		avail, _ := strconv.ParseUint(a[4], 10, 64)
		use, _ := strconv.ParseUint(strings.TrimSuffix(a[5], "%"), 10, 64)
		i = append(i,
			&mdInfo{
				Path:  a[6],
				Size:  size,
				Used:  used,
				Avail: avail,
				Use:   use,
			})
	}

	return i, nil
}

func getNodeStat() ([]*nodeStat, error) {
	args := []string{
		"node", "stat", "-r",
	}

	output, err := dog(args)
	if err != nil {
		return nil, err
	}

	var ns []*nodeStat
	var types = []string{"client", "peer"}
	for i, line := range output {
		a := strings.Fields(line)
		active, _ := strconv.ParseUint(a[0], 10, 64)
		total, _ := strconv.ParseUint(a[1], 10, 64)
		write, _ := strconv.ParseUint(a[2], 10, 64)
		read, _ := strconv.ParseUint(a[3], 10, 64)
		remove, _ := strconv.ParseUint(a[4], 10, 64)
		flush, _ := strconv.ParseUint(a[5], 10, 64)
		allWrite, _ := strconv.ParseUint(a[6], 10, 64)
		allRead, _ := strconv.ParseUint(a[7], 10, 64)
		ns = append(ns,
			&nodeStat{
				Type:     types[i],
				Active:   active,
				Total:    total,
				Write:    write,
				Read:     read,
				Remove:   remove,
				Flush:    flush,
				AllWrite: allWrite,
				AllRead:  allRead,
			})
	}

	return ns, nil
}

func dog(args []string) ([]string, error) {
	var stdOut []string

	cmd := exec.Command("dog", args...)
	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		return nil, err
	}

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			outPut := scanner.Text()
			stdOut = append(stdOut, outPut)
		}
	}()

	err = cmd.Start()
	if err != nil {
		return nil, err
	}

	err = cmd.Wait()
	if err != nil {
		return nil, err
	}

	return stdOut, nil
}
