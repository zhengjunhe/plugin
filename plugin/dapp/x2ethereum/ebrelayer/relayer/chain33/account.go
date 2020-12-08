package dplatform

import (
	dplatformCommon "github.com/33cn/dplatform/common"
	"github.com/ethereum/go-ethereum/crypto"

	//dbm "github.com/33cn/dplatform/common/db"
	dplatformTypes "github.com/33cn/dplatform/types"
	wcom "github.com/33cn/dplatform/wallet/common"
	x2ethTypes "github.com/33cn/plugin/plugin/dapp/x2ethereum/ebrelayer/types"
)

var (
	dplatformAccountKey = []byte("DplatformAccount4Relayer")
	start             = int(1)
)

//GetAccount ...
func (dplatformRelayer *Relayer4Dplatform) GetAccount(passphrase string) (privateKey, addr string, err error) {
	accountInfo, err := dplatformRelayer.db.Get(dplatformAccountKey)
	if nil != err {
		return "", "", err
	}
	ethAccount := &x2ethTypes.Account4Relayer{}
	if err := dplatformTypes.Decode(accountInfo, ethAccount); nil != err {
		return "", "", err
	}
	decryptered := wcom.CBCDecrypterPrivkey([]byte(passphrase), ethAccount.Privkey)
	privateKey = dplatformCommon.ToHex(decryptered)
	addr = ethAccount.Addr
	return
}

//GetAccountAddr ...
func (dplatformRelayer *Relayer4Dplatform) GetAccountAddr() (addr string, err error) {
	accountInfo, err := dplatformRelayer.db.Get(dplatformAccountKey)
	if nil != err {
		relayerLog.Info("GetValidatorAddr", "Failed to get account from db due to:", err.Error())
		return "", err
	}
	ethAccount := &x2ethTypes.Account4Relayer{}
	if err := dplatformTypes.Decode(accountInfo, ethAccount); nil != err {
		relayerLog.Info("GetValidatorAddr", "Failed to decode due to:", err.Error())
		return "", err
	}
	addr = ethAccount.Addr
	return
}

//ImportPrivateKey ...
func (dplatformRelayer *Relayer4Dplatform) ImportPrivateKey(passphrase, privateKeyStr string) (addr string, err error) {
	privateKeySlice, err := dplatformCommon.FromHex(privateKeyStr)
	if nil != err {
		return "", err
	}
	privateKey, err := crypto.ToECDSA(privateKeySlice)
	if nil != err {
		return "", err
	}

	ethSender := crypto.PubkeyToAddress(privateKey.PublicKey)
	dplatformRelayer.privateKey4Ethereum = privateKey
	dplatformRelayer.ethSender = ethSender
	dplatformRelayer.unlock <- start

	addr = dplatformCommon.ToHex(ethSender.Bytes())
	encryptered := wcom.CBCEncrypterPrivkey([]byte(passphrase), privateKeySlice)
	ethAccount := &x2ethTypes.Account4Relayer{
		Privkey: encryptered,
		Addr:    addr,
	}
	encodedInfo := dplatformTypes.Encode(ethAccount)
	err = dplatformRelayer.db.SetSync(dplatformAccountKey, encodedInfo)

	return
}

//StoreAccountWithNewPassphase ...
func (dplatformRelayer *Relayer4Dplatform) StoreAccountWithNewPassphase(newPassphrase, oldPassphrase string) error {
	accountInfo, err := dplatformRelayer.db.Get(dplatformAccountKey)
	if nil != err {
		relayerLog.Info("StoreAccountWithNewPassphase", "pls check account is created already, err", err)
		return err
	}
	ethAccount := &x2ethTypes.Account4Relayer{}
	if err := dplatformTypes.Decode(accountInfo, ethAccount); nil != err {
		return err
	}
	decryptered := wcom.CBCDecrypterPrivkey([]byte(oldPassphrase), ethAccount.Privkey)
	encryptered := wcom.CBCEncrypterPrivkey([]byte(newPassphrase), decryptered)
	ethAccount.Privkey = encryptered
	encodedInfo := dplatformTypes.Encode(ethAccount)
	return dplatformRelayer.db.SetSync(dplatformAccountKey, encodedInfo)
}

//RestorePrivateKeys ...
func (dplatformRelayer *Relayer4Dplatform) RestorePrivateKeys(passphrase string) error {
	accountInfo, err := dplatformRelayer.db.Get(dplatformAccountKey)
	if nil != err {
		relayerLog.Info("No private key saved for Relayer4Dplatform")
		return nil
	}
	ethAccount := &x2ethTypes.Account4Relayer{}
	if err := dplatformTypes.Decode(accountInfo, ethAccount); nil != err {
		relayerLog.Info("RestorePrivateKeys", "Failed to decode due to:", err.Error())
		return err
	}
	decryptered := wcom.CBCDecrypterPrivkey([]byte(passphrase), ethAccount.Privkey)
	privateKey, err := crypto.ToECDSA(decryptered)
	if nil != err {
		relayerLog.Info("RestorePrivateKeys", "Failed to ToECDSA:", err.Error())
		return err
	}

	dplatformRelayer.rwLock.Lock()
	dplatformRelayer.privateKey4Ethereum = privateKey
	dplatformRelayer.ethSender = crypto.PubkeyToAddress(privateKey.PublicKey)
	dplatformRelayer.rwLock.Unlock()
	dplatformRelayer.unlock <- start
	return nil
}

//func (dplatformRelayer *Relayer4Dplatform) UpdatePrivateKey(Passphrase, privateKey string) error {
//	return nil
//}
