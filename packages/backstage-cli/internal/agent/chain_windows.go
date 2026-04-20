//go:build windows

package agent

import (
	"fmt"

	"github.com/shirou/gopsutil/v3/process"
)

func queryChainLayerOS(pid int) (ChainBlock, error) {
	p, err := process.NewProcess(int32(pid))
	if err != nil {
		return ChainBlock{
			PID:       pid,
			PPID:      -1,
			PGID:      -1,
			SID:       -1,
			UID:       -1,
			RawPsLine: fmt.Sprintf("<gopsutil err: %v; pid=%d>", err, pid),
		}, err
	}

	ppid, err := p.Ppid()
	if err != nil {
		ppid = -1
	}

	name, err := p.Name()
	if err != nil {
		name = ""
	}

	cmdline, err := p.Cmdline()
	if err != nil {
		cmdline = ""
	}

	username, err := p.Username()
	if err != nil {
		username = ""
	}

	createTime, err := p.CreateTime()
	if err != nil {
		createTime = 0
	}

	raw := fmt.Sprintf("pid=%d ppid=%d name=%q cmdline=%q username=%q create_time=%d",
		pid, ppid, name, cmdline, username, createTime)

	layer := ChainBlock{
		PID:       pid,
		PPID:      int(ppid),
		PGID:      -1,
		SID:       -1,
		UID:       -1,
		Comm:      name,
		Command:   cmdline,
		Username:  username,
		Lstart:    "",
		RawPsLine: raw,
	}

	return layer, nil
}
