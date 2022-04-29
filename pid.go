package main

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func isInt(s string) bool {
	_, err := strconv.Atoi(s)
	return err == nil
}

// list all pids in /proc
func list_pid() ([]int64, error) {
	entry_ls, err := os.ReadDir("/proc")
	if err != nil {
		return nil, err
	}
	var pid_ls []int64
	for _, entry := range entry_ls {
		if entry.IsDir() && isInt(entry.Name()) {
			pid, err := strconv.ParseInt(entry.Name(), 10, 64)
			if err != nil {
				return nil, err
			}
			pid_ls = append(pid_ls, pid)
		}
	}
	return pid_ls, nil
}

// reads '/proc/*/fd/*' to grab process IDs
func process_id() ([]string, error) {
	files, err := filepath.Glob("/proc/*/fd/*")
	if err != nil {
		return nil, err
	}
	return files, nil
}

func id_from_fd(s string) (int64, error) {
	return strconv.ParseInt(strings.Split(s, "/")[2], 10, 64)
}
