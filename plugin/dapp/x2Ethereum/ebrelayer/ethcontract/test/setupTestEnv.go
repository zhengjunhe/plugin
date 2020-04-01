package test

import (
	"crypto/ecdsa"
	"github.com/33cn/plugin/plugin/dapp/x2Ethereum/ebrelayer/ethtxs"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"math/big"
)

func PrepareTestEnv() (bind.ContractBackend, *ethtxs.DeployPara) {
	genesiskey, _ := crypto.GenerateKey()
	alloc := make(core.GenesisAlloc)
	genesisAddr := crypto.PubkeyToAddress(genesiskey.PublicKey)
	genesisAccount := core.GenesisAccount{
		Balance:    big.NewInt(10000000000 * 10000),
		PrivateKey: crypto.FromECDSA(genesiskey),
	}
	alloc[genesisAddr] = genesisAccount
	gasLimit := uint64(2999280)
	sim := backends.NewSimulatedBackend(alloc, gasLimit)

	var InitValidators []common.Address
	var ValidatorPriKey []*ecdsa.PrivateKey
	for i := 0; i < 3; i++ {
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)
		InitValidators = append(InitValidators, addr)
		ValidatorPriKey = append(ValidatorPriKey, key)
	}

	InitPowers := []*big.Int{big.NewInt(80), big.NewInt(10), big.NewInt(10)}

	para := &ethtxs.DeployPara{
		PrivateKey:     genesiskey,
		Deployer:       genesisAddr,
		Operator:       genesisAddr,
		InitValidators: InitValidators,
		ValidatorPriKey:ValidatorPriKey,
		InitPowers:     InitPowers,
	}

	return sim, para
}
