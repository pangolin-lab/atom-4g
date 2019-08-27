package Service4G

import (
	"github.com/pangolin-lab/go-node-4g/service/rpcMsg"
	"net"
)

type flowManager interface {
	//check bucket level
	RequireService(conn net.Conn) rpcMsg.CreditOnNode
	//calculate usage by self
	CalculateUsage()(tempUsageLocal uint64,bucketLvlLocal uint64)
}