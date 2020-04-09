package setup

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
	var genesiskey, _ = crypto.GenerateKey()
	genesisAddr := crypto.PubkeyToAddress(genesiskey.PublicKey)

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
		DeployPrivateKey: genesiskey,
		Deployer:         genesisAddr,
		Operator:         genesisAddr,
		InitValidators:   InitValidators,
		ValidatorPriKey:  ValidatorPriKey,
		InitPowers:       InitPowers,
	}

	alloc := make(core.GenesisAlloc)
	genesisAccount := core.GenesisAccount{
		Balance:    big.NewInt(10000000000 * 10000),
		PrivateKey: crypto.FromECDSA(genesiskey),
	}
	v0 := core.GenesisAccount{
		Balance:    big.NewInt(10000000000),
		PrivateKey: crypto.FromECDSA(ValidatorPriKey[0]),
	}

	v1 := core.GenesisAccount{
		Balance:    big.NewInt(10000000000),
		PrivateKey: crypto.FromECDSA(ValidatorPriKey[1]),
	}

	alloc[genesisAddr] = genesisAccount
	alloc[InitValidators[0]] = v0
	alloc[InitValidators[1]] = v1
	gasLimit := uint64(0)
	sim := backends.NewSimulatedBackend(alloc, gasLimit)

	return sim, para
}
