package _GService

import (
	"fmt"
	"github.com/proton-lab/atom-4g/wallet"
	"net"
	"strconv"
)

type _4GProxy struct {
	Done   chan error
	Wallet *wallet.Wallet
	TunSrc Tun2Pipe
}

func NewProxy(addr string, w *wallet.Wallet, t Tun2Pipe) (*_4GProxy, error) {
	ap := &_4GProxy{
		Wallet:      w,
		TunSrc:      t,
		Done:        make(chan error),
	}
	return ap, nil
}

func (pp *_4GProxy) Proxying() {

	go pp.TunSrc.Proxying(pp.Done)

	go pp.Wallet.Running(pp.Done)

	go pp.Accepting(pp.Done)

	select {
	case err := <-pp.Done:
		fmt.Printf("_4GProxy exit for:%s", err.Error())
	}

	pp.Finish()
}

func (pp *_4GProxy) consume(conn net.Conn) {
	defer conn.Close()

	tgtAddr := pp.TunSrc.GetTarget(conn)

	if len(tgtAddr) < 10 {
		fmt.Println("\nNo such connection's target address:->", conn.RemoteAddr().String())
		return
	}
	fmt.Println("\n Proxying target address:", tgtAddr)

	//TODO::match PAC file in ios or android logic
	pipe := pp.Wallet.SetupPipe(conn, tgtAddr)
	if nil == pipe {
		fmt.Println("Create pipe failed:", tgtAddr)
		return
	}


	rAddr := conn.RemoteAddr().String()
	_, port, _ := net.SplitHostPort(rAddr)
	keyPort, _ := strconv.Atoi(port)
	pp.TunSrc.RemoveFromSession(keyPort)

	//TODO::need to make sure is this ok
	fmt.Printf("\n\nPipe(%s) for(%s) is closing", rAddr, tgtAddr)
}

func (pp *_4GProxy) Finish() {

	if pp.TCPListener != nil {
		pp.TCPListener.Close()
		pp.TCPListener = nil
	}

	if pp.Wallet != nil {
		pp.Wallet.Finish()
	}

	if pp.TunSrc != nil {
		pp.TunSrc.Finish()
	}
}
