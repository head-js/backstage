//go:build !darwin && !windows

package agent

import (
	"fmt"
	"runtime"

	"com.lisitede.backstage/framework"
)

func queryChainLayerOS(pid int) (ChainLayer, error) {
	msg := fmt.Sprintf("chain query on %s", runtime.GOOS)
	return ChainLayer{}, framework.NotImplementedException(msg)
}
