package _GService

import (
	"github.com/pangolin-lab/proton-node-4g/service/rpcMsg"
	"net"
)

type fakeFlowManger struct {

}

func (ffm *fakeFlowManger)RequireService(conn net.Conn) rpcMsg.BucketCheck{

}