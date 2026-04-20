//go:build darwin

package agent

import (
	"fmt"
	"os/exec"
	"strconv"
	"strings"
)

func queryChainLayerOS(pid int) (ChainLayer, error) {
	out, err := exec.Command("ps",
		"-o", "pid=,ppid=,pgid=,sess=,uid=,tty=,user=,lstart=,comm=,command=",
		"-p", strconv.Itoa(pid)).CombinedOutput()
	if err != nil {
		return ChainLayer{
			PID: pid, PPID: -1, PGID: -1, SID: -1, UID: -1,
			RawPsLine: fmt.Sprintf("<ps err: %v; output=%q>", err, strings.TrimSpace(string(out))),
		}, err
	}
	raw := strings.TrimRight(string(out), "\r\n")
	layer := ChainLayer{
		RawPsLine: raw,
		PID:       pid,
		PPID:      -1,
		PGID:      -1,
		SID:       -1,
		UID:       -1,
	}
	fields := strings.Fields(raw)
	if len(fields) < 13 {
		return layer, nil
	}
	if n, e := strconv.Atoi(fields[0]); e == nil {
		layer.PID = n
	}
	if n, e := strconv.Atoi(fields[1]); e == nil {
		layer.PPID = n
	}
	if n, e := strconv.Atoi(fields[2]); e == nil {
		layer.PGID = n
	}
	if n, e := strconv.Atoi(fields[3]); e == nil {
		layer.SID = n
	}
	if n, e := strconv.Atoi(fields[4]); e == nil {
		layer.UID = n
	}
	layer.TTY = fields[5]
	layer.Username = fields[6]
	layer.Lstart = strings.Join(fields[7:12], " ")
	layer.Comm = fields[12]
	layer.Command = rawAfterNTokens(raw, 13)
	return layer, nil
}

func rawAfterNTokens(raw string, n int) string {
	s := raw
	for i := 0; i < n; i++ {
		j := 0
		for j < len(s) && (s[j] == ' ' || s[j] == '\t') {
			j++
		}
		s = s[j:]
		if s == "" {
			return ""
		}
		j = 0
		for j < len(s) && s[j] != ' ' && s[j] != '\t' {
			j++
		}
		s = s[j:]
	}
	j := 0
	for j < len(s) && (s[j] == ' ' || s[j] == '\t') {
		j++
	}
	return s[j:]
}
