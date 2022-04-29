package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

func get_env() []string {
	return strings.Split(os.Getenv("PATH"), ":")
}
func uptime() string {
	duration := time.Since(boot_time())
	h := int(duration.Hours())
	m := int(duration.Minutes())
	s := int(duration.Seconds())
	res := fmt.Sprintf("%d days, %d hours, %d minutes, %d seconds", h/24, h%24, m, s)
	return res
}
func uptime2() string {
	duration := time.Since(boot_time())

	return duration.String()
}

// return the last reboot time
func boot_time() time.Time {
	// open a file
	f, err := os.Open("/proc/stat")
	if err != nil {
		log.Fatal(err)
	}
	defer f.Close()
	// create a new scanner to read the file
	scanner := bufio.NewScanner(f)
	var btime int
	// default split function is bufio.ScanLines, which will split on newlines
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.HasPrefix(line, "btime") {
			fmt.Sscanf(line, "btime %d", &btime)
		}
	}
	return time.Unix(int64(btime), 0)
}
