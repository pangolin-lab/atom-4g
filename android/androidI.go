package androidLib

import "C"
import (
	"fmt"
	"github.com/Iuduxras/atom-4g/Service4G"
	"github.com/Iuduxras/atom-4g/ethereum"
	"github.com/Iuduxras/atom-4g/tun2Pipe"
	"github.com/Iuduxras/atom-4g/wallet"
	"github.com/Iuduxras/pangolin-node-4g/account"
	"github.com/Iuduxras/pangolin-node-4g/pbs/pipeProxy"
	"github.com/btcsuite/btcutil/base58"
)

type VpnDelegate interface {
	tun2Pipe.VpnDelegate
	GetBootPath() string
}

const Separator = "@@@"
var proxyConf = &pipeProxy.ProxyConfig{}
var _instance *Service4G.Consumer4G = nil

type ConsumeDelegate interface {
	GetBootPath() string
}

//consumer setup
func InitConsumer(addr, cipher, url, boot, ip,mac,IPs ,dbPath string,d ConsumeDelegate) error{

	proxyConf.WConfig = &wallet.WConfig{
		BCAddr:     addr,
		Cipher:     cipher,
		SettingUrl: url,
		Ip:         ip,
		Mac:        mac,
	}

	proxyConf.BootNodes = boot
	mis := proxyConf.FindBootServers(d.GetBootPath())
	if len(mis) == 0 {
		return fmt.Errorf("no valid boot strap node")
	}

	proxyConf.ServerId = mis[0]
	println(proxyConf.String())

	return nil
}

func SetupConsumer(password,locAddr string) error{
	w, err := wallet.NewWallet(proxyConf.WConfig, password)
	if err != nil {
		return err
	}
	consumer, e := Service4G.NewConsumer(locAddr, w)
	if e != nil {
		panic(err)
	}
	_instance = consumer
	return nil
}

func Consuming(){
	if _instance ==nil{
		return
	}
	_instance.Consuming()
	_instance = nil
}

func StopConsuming(){
	if _instance !=nil {
		_instance.Done <- fmt.Errorf("user close this")
		_instance = nil
	}
}

/*
	returns:
	{
		Accepted bool
		Credit   int64
	}
*/
func Query() string{
	if _instance !=nil{
		return _instance.Query()
	}else{
		return ""
	}
}

func Recharge(no int) bool{
	if _instance !=nil{
		if err:=_instance.Recharge(no);err!=nil{
			return false
		}else{
			return true
		}
	}else{
		return false
	}
}


//accounts and ethereum opts

func VerifyAccount(addr, cipher, password string) bool {
	if _, err := account.AccFromString(addr, cipher, password); err != nil {
		fmt.Println("Valid Account:", err)
		return false
	}
	return true
}

func CreateAccount(password string) string {

	key, err := account.GenerateKey(password)
	if err != nil {
		return ""
	}
	address := key.ToNodeId().String()
	cipherTxt := base58.Encode(key.LockedKey)

	return address + Separator + cipherTxt
}

func IsProtonAddress(address string) bool {
	return account.ID(address).IsValid()
}

func LoadEthAddrByProtonAddr(protonAddr string) string {
	return ethereum.BoundEth(protonAddr)
}

func EthBindings(ETHAddr string) string {
	ethB, no := ethereum.BasicBalance(ETHAddr)
	if ethB == nil {
		return ""
	}

	return fmt.Sprintf("%f"+Separator+"%d",
		ethereum.ConvertByDecimal(ethB),
		no)
}

func CreateEthAccount(password, directory string) string {
	return ethereum.CreateEthAccount2(password, directory)
}

func VerifyEthAccount(cipherTxt, pwd string) bool {
	return ethereum.VerifyEthAccount(cipherTxt, pwd)
}

func BindProtonAddress(protonAddr, cipherKey, password string) string {
	tx, err := ethereum.Bind(protonAddr, cipherKey, password)
	if err != nil {
		fmt.Printf("\nBind proton addr(%s) err:%s", protonAddr, err)
		return err.Error()
	}
	return tx
}
func UnbindProtonAddress(protonAddr, cipherKey, password string) string {
	tx, err := ethereum.Unbind(protonAddr, cipherKey, password)
	if err != nil {
		fmt.Printf("\nBind proton addr(%s) err:%s", protonAddr, err)
		return err.Error()
	}
	return tx
}
