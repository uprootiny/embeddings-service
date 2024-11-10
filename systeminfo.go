package main

import (
    "os/exec"
    "strings"
)

// type SystemInfo struct {
//     Hostname     string `json:"hostname"`
//     OS           string `json:"os"`
//     Uptime       string `json:"uptime"`
//     Kernel       string `json:"kernel"`
//     Architecture string `json:"architecture"`
// }

func GetSystemInfo() SystemInfo {
    hostname, _ := exec.Command("hostname").Output()
    osInfo, _ := exec.Command("lsb_release", "-d").Output()
    uptime, _ := exec.Command("uptime", "-p").Output()
    kernel, _ := exec.Command("uname", "-r").Output()
    arch, _ := exec.Command("uname", "-m").Output()

    return SystemInfo{
        Hostname:     strings.TrimSpace(string(hostname)),
        OS:           strings.TrimSpace(strings.Replace(string(osInfo), "Description:\t", "", 1)),
        Uptime:       strings.TrimSpace(string(uptime)),
        Kernel:       strings.TrimSpace(string(kernel)),
        Architecture: strings.TrimSpace(string(arch)),
    }
}
