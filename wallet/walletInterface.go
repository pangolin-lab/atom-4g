package wallet

type UserWallet interface {
	Running(done chan error)
	Finish()
}