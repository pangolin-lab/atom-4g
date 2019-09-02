package Service4G

import (
	"fmt"
	"github.com/Iuduxras/atom-4g/wallet"
)

type Consumer4G struct {
	Done   chan error
	Wallet wallet.UserWallet
}

func NewConsumer(addr string, w wallet.UserWallet) (*Consumer4G, error) {
	ap := &Consumer4G{
		Wallet:      w,
		Done:        make(chan error),
	}
	return ap, nil
}

func (pp *Consumer4G) Consuming() {

	go pp.Wallet.Running(pp.Done)

	select {
	case err := <-pp.Done:
		fmt.Printf("Consumer4G exit for:%s", err.Error())
	}

	pp.Finish()
}


func (pp *Consumer4G) Finish() {

	if pp.Wallet != nil {
		pp.Wallet.Finish()
	}
}
