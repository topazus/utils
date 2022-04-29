package main

import (
	"bufio"
	"fmt"

	//"log"
	"os"
	"strings"
	"syscall"

	//"github.com/kr/pretty"
	"golang.org/x/sys/unix"
)

type DiskInfo struct {
	Type   string `json:"partition type"`
	FsType string `json:"filesystem type"`
	Size   uint64 `json:"size"`
	Used   uint64 `json:"used"`
	Free   uint64 `json:"free"`
	Avail  uint64 `json:"avail"`
}

func main() {

	fmt.Println("----")
	print_disk_info("/boot")
	print_disk_info("/boot/efi")
	print_disk_info("/")
	print_disk_info("/dev/shm")
	print_disk_info("/dev")
	print_disk_info("/run")

	ls, _ := list_pid()
	fmt.Println(ls)
	a, _ := process_id()
	fmt.Println(a)

}

// disk usage of path/disk
func DiskUsage(path string) (DiskInfo, error) {
	var fs unix.Statfs_t
	err := unix.Statfs(path, &fs)
	//fs := syscall.Statfs_t{}
	//err := syscall.Statfs(path, &fs)
	if err != nil {
		return DiskInfo{}, err
	}
	var disk DiskInfo
	disk.Type = fsTypeMap[fs.Type]
	disk.Size = fs.Blocks * uint64(fs.Bsize)
	disk.Avail = fs.Bavail * uint64(fs.Bsize)
	disk.Free = fs.Bfree * uint64(fs.Bsize)
	disk.Used = disk.Size - disk.Free
	return disk, nil
}

const (
	B  = 1
	KB = 1024 * B
	MB = 1024 * KB
	GB = 1024 * MB
	TB = 1024 * GB
)

func print_disk_info(path string) {
	disk, _ := DiskUsage(path)

	fmt.Printf("%10s:", path)
	if disk.Size > MB && disk.Size < GB {
		fmt.Printf("%8.2f MB %8.2fMB	%8.2fMB\n", float64(disk.Size)/float64(MB), float64(disk.Avail)/float64(MB), float64(disk.Free)/float64(MB))
	} else if disk.Size > GB {
		fmt.Printf("%8.2f GB %8.2fGB	%8.2fGB\n", float64(disk.Size)/float64(GB), float64(disk.Avail)/float64(GB), float64(disk.Free)/float64(GB))
	} else {
		fmt.Printf("%f %f %f\n", float64(disk.Size), float64(disk.Avail), float64(disk.Free))
	}
}

/*
type Statfs_t struct {
	Type    int64
	Bsize   int64
	Blocks  uint64
	Bfree   uint64
	Bavail  uint64
	Files   uint64
	Ffree   uint64
	Fsid    Fsid
	Namelen int64
	Frsize  int64
	Flags   int64
	Spare   [4]int64
}
__fsword_t f_type;     Type of filesystem (see below)
__fsword_t f_bsize;    Optimal transfer block size
fsblkcnt_t f_blocks;   Total data blocks in filesystem
fsblkcnt_t f_bfree;    Free blocks in filesystem
fsblkcnt_t f_bavail;   Free blocks available to unprivileged user
fsfilcnt_t f_files;    Total inodes in filesystem
fsfilcnt_t f_ffree;    Free inodes in filesystem
fsid_t     f_fsid;     Filesystem ID
__fsword_t f_namelen;  Maximum length of filenames
__fsword_t f_frsize;   Fragment size (since Linux 2.6)
__fsword_t f_flags;    Mount flags of filesystem (since Linux 2.6.36)
__fsword_t f_spare[xxx]; Padding bytes reserved for future use
*/

func parse_mount_line(line string) []string {
	var info []string
	info = append(info, strings.Fields(strings.Split(line, "-")[0])...)
	info = append(info, strings.Fields(strings.Split(line, "-")[1])...)
	return info
}
func read_lines(path string) ([]string, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	// default split function is bufio.ScanLines, which will split on newlines
	scanner.Split(bufio.ScanLines)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
	}
	return lines, scanner.Err()
}
func mounts() ([]MountInfo, error) {
	//file_name:="/proc/mounts"
	file_name := "/proc/self/mountinfo"
	lines, err := read_lines(file_name)
	if err != nil {
		return []MountInfo{}, err
	}
	mount_list := []MountInfo{}
	for _, line := range lines {
		info := parse_mount_line(line)
		if len(info) != 11 {
			continue
		}
		mount_point := info[MountPoint]
		mount_options := info[MountOptions]
		fstype := info[FsType]
		mount_source := info[MountSource]
		var fs syscall.Statfs_t
		err := syscall.Statfs(mount_point, &fs)
		if err != nil {
			continue
		}
		mount_device := MountInfo{
			MountSource:  mount_source,
			MountPoint:   mount_point,
			MountOptions: mount_options,
			Fstype:       fstype,
			Type:         fsTypeMap[fs.Type],
			Total:        fs.Blocks * uint64(fs.Bsize),
			Free:         fs.Bfree * uint64(fs.Bsize),
			Used:         (fs.Blocks - fs.Bfree) * uint64(fs.Bsize),
			TotalInodes:  fs.Files,
			InodesFree:   fs.Ffree,
			InodesUsed:   fs.Files - fs.Ffree,
			Blocks:       fs.Blocks,
			BlockSize:    uint64(fs.Bsize),
			Metadata:     fs,
		}
		mount_list = append(mount_list, mount_device)
	}
	return mount_list, nil
}
