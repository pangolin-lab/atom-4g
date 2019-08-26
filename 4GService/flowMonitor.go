package _GService

import (
	"github.com/proton-lab/proton-node-4g/service/rpcMsg"
	"net"
)

type flowManager interface {
	//check bucket level
	RequireService(conn net.Conn) rpcMsg.BucketCheck
	//calculate usage by self
	CalculateUsage()(tempUsageLocal uint64,bucketLvlLocal uint64)
}