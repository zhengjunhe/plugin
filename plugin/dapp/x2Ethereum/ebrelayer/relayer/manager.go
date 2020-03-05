package relayer

import (
	"errors"
	dbm "github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/common/log/log15"
	chain33Types "github.com/33cn/chain33/types"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/relayer/chain33"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/relayer/ethereum"
	relayerTypes "github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/types"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/utils"
	"sync"
	"sync/atomic"
)

var (
	mlog = log15.New("relayer manager", "manager")
)

const (
	Locked   = int32(1)
	Unlocked = int32(99)
	EncryptEnable = int64(1)
)

type RelayerManager struct {
	chain33Relayer *chain33.Chain33Relayer
	ethRelayer *ethereum.EthereumRelayer
	store *Store
	isLocked     int32
	mtx          sync.Mutex
	encryptFlag   int64
	passphase     string
}
//实现记录
//1.验证人的私钥需要通过cli命令行进行导入，且chain33和ethereum两种不同的验证人需要分别导入
//2.显示或者重新替换原有的私钥首先需要通过passpin进行unlock的操作

func NewRelayerManager(chain33Relayer *chain33.Chain33Relayer, ethRelayer *ethereum.EthereumRelayer, db dbm.DB) *RelayerManager {
	manager := &RelayerManager{
		chain33Relayer: chain33Relayer,
		ethRelayer:     ethRelayer,
		store:          NewStore(db),
		isLocked:       Locked,
		mtx:            sync.Mutex{},
		encryptFlag:    0,
		passphase:      "",
	}
	manager.encryptFlag = manager.store.GetEncryptionFlag()
	return manager
}

func (manager *RelayerManager) SetPassphase(setPasswdReq relayerTypes.ReqWalletSetPasswd, result *interface{}) error {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()
	// 新密码合法性校验
	if !utils.IsValidPassWord(setPasswdReq.NewPassphase) {
		return chain33Types.ErrInvalidPassWord
	}
	//保存钱包的锁状态，需要暂时的解锁，函数退出时再恢复回去
	tempislock := atomic.LoadInt32(&manager.isLocked)
	atomic.CompareAndSwapInt32(&manager.isLocked, Locked, Unlocked)

	defer func() {
		//wallet.isWalletLocked = tempislock
		atomic.CompareAndSwapInt32(&manager.isLocked, Unlocked, tempislock)
	}()

	// 钱包已经加密需要验证oldpass的正确性
	if len(manager.passphase) == 0 && manager.encryptFlag == 1 {
		isok := manager.store.VerifyPasswordHash(setPasswdReq.OldPassphase)
		if !isok {
			mlog.Error("SetPassphase Verify Oldpasswd fail!")
			return chain33Types.ErrVerifyOldpasswdFail
		}
	}

	if len(manager.passphase) != 0 && setPasswdReq.OldPassphase != manager.passphase {
		mlog.Error("SetPassphase Oldpass err!")
		return chain33Types.ErrVerifyOldpasswdFail
	}

	//使用新的密码生成passwdhash用于下次密码的验证
	newBatch := manager.store.NewBatch(true)
	err := manager.store.SetPasswordHash(setPasswdReq.NewPassphase, newBatch)
	if err != nil {
		mlog.Error("SetPassphase", "SetPasswordHash err", err)
		return err
	}
	//设置钱包加密标志位
	err = manager.store.SetEncryptionFlag(newBatch)
	if err != nil {
		mlog.Error("SetPassphase", "SetEncryptionFlag err", err)
		return err
	}

	err = manager.ethRelayer.StoreAccountWithNewPassphase(setPasswdReq.NewPassphase, setPasswdReq.OldPassphase)
	if err != nil {
		mlog.Error("SetPassphase", "StoreAccountWithNewPassphase err", err)
		return err
	}

	err = newBatch.Write()
	if err != nil {
		mlog.Error("ProcWalletSetPasswd newBatch.Write", "err", err)
		return err
	}
	manager.passphase = setPasswdReq.NewPassphase
	atomic.StoreInt64(&manager.encryptFlag, EncryptEnable)
	return nil
}

//进行unlok操作
func (manager *RelayerManager) Unlock(passphase string, result *interface{}) error {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()
	if EncryptEnable != manager.encryptFlag {
		return errors.New("pls set passphase first")
	}
	if Unlocked == manager.isLocked {
		return errors.New("unlock already")
	}

	if !manager.store.VerifyPasswordHash(passphase) {
		return errors.New("wrong passphase")
	}

	manager.isLocked = Unlocked
	manager.passphase = passphase
	*result = "Succeed to unlock"

	return nil
}

//锁定操作，该操作一旦执行，就不能替换验证人的私钥，需要重新unlock之后才能修改
func (manager *RelayerManager) Lock(param interface{}, result *interface{}) error {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()
	if err := manager.checkPermission(); nil != err {
		return err
	}
	manager.isLocked = Locked
	return nil
}

//导入chain33relayer验证人的私钥,该私钥实际用于向ethereum提交验证交易时签名使用
func (manager *RelayerManager) ImportChain33RelayerPrivateKey(privateKey string, result *interface{}) error {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()
	if err := manager.checkPermission(); nil != err {
		return err
	}
	_, err := manager.chain33Relayer.ImportPrivateKey(manager.passphase, privateKey)
	if err != nil {
		mlog.Error("ImportChain33ValidatorPrivateKey", "Failed due to casue:", err.Error())
		return err
	}

	*result = "Succeed to import private key for chain33 relayer"
	return nil
}

//生成以太坊私钥
func (manager *RelayerManager) GenerateEthereumPrivateKey(param interface{}, result *interface{}) error {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()
	if err := manager.checkPermission(); nil != err {
		return err
	}
	account4Show := relayerTypes.Account4Show{}
	var err error
	account4Show.Privkey, account4Show.Addr, err = manager.ethRelayer.NewAccount(manager.passphase)
	if nil != err {
		return err
	}
	*result = account4Show
	return nil
}

//为ethrelayer导入chain33私钥，为向chain33发送交易时进行签名使用
func (manager *RelayerManager) ImportChain33PrivateKey4EthRelayer(privateKey string, result *interface{}) error {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()
	if err := manager.checkPermission(); nil != err {
		return err
	}
	if err := manager.ethRelayer.ImportChain33PrivateKey(manager.passphase, privateKey); nil != err {
		return err
	}
	*result = "Succeed to import chain33 private key for ethereum relayer"
	return nil
}

//显示在chain33中以验证人validator身份进行登录的地址
func (manager *RelayerManager) ShowChain33RelayerValidator(param interface{}, result *interface{}) error {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()
	var err error
	*result, err = manager.chain33Relayer.GetAccountAddr()
	if nil != err {
		return err
	}

	return nil
}
//显示在Ethereum中以验证人validator身份进行登录的地址
func (manager *RelayerManager) ShowEthRelayerValidator(param interface{}, result *interface{}) error {
	manager.mtx.Lock()
	defer manager.mtx.Unlock()
	var err error
	*result, err = manager.ethRelayer.GetValidatorAddr()
	if nil != err {
		return err
	}
	return nil
}

func (manager *RelayerManager) checkPermission() error {
	if EncryptEnable != manager.encryptFlag {
		return errors.New("pls set passphase first")
	}
	if Locked == manager.isLocked {
		return errors.New("pls unlock this relay-manager first")
	}
	return nil
}

func (manager *RelayerManager) ShowEthRelayer2EthTxs(param interface{}, result *interface{}) error {
	*result = manager.ethRelayer.QueryTxhashRelay2Eth()
	return nil
}

func (manager *RelayerManager) ShowEthRelayer2Chain33Txs(param interface{}, result *interface{}) error {
	*result = manager.ethRelayer.QueryTxhashRelay2Chain33()
	return nil
}

func (manager *RelayerManager) ShowChain33Relayer2EthTxs(param interface{}, result *interface{}) error {
	*result = manager.chain33Relayer.QueryTxhashRelay2Eth()
	return nil
}

func (manager *RelayerManager) ShowTxsEth2chain33TxLock(param interface{}, result *interface{}) error {
	return nil
}

func (manager *RelayerManager) ShowTxsEth2chain33TxBurn(param interface{}, result *interface{}) error {
	return nil
}

func (manager *RelayerManager) ShowTxsChain33ToEthTxLock(param interface{}, result *interface{}) error {

	return nil
}

func (manager *RelayerManager) ShowTxsChain33ToEthTxBurn(param interface{}, result *interface{}) error {

	return nil
}
