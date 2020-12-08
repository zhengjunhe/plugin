package ethtxs

import (
	"context"
	"crypto/ecdsa"
	"fmt"
	"math/big"
	"strings"
	"testing"

	dplatformTypes "github.com/33cn/dplatform/types"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/ebrelayer/ethcontract/generated"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/ebrelayer/ethinterface"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/ebrelayer/events"
	ebrelayerTypes "github.com/33cn/plugin/plugin/dapp/x2ethereum/ebrelayer/types"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
	x2ethTypes "github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind/backends"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/core"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

var (
	dplatformAddr  = "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
	ethAddr      = "0x92C8b16aFD6d423652559C6E266cBE1c29Bfd84f"
	ethTokenAddr = "0x0000000000000000000000000000000000000000"
)

type suiteContracts struct {
	suite.Suite
	para            *DeployPara
	sim             *ethinterface.SimExtend
	x2EthContracts  *X2EthContracts
	x2EthDeployInfo *X2EthDeployInfo
}

func TestRunSuiteX2Ethereum(t *testing.T) {
	log := new(suiteContracts)
	suite.Run(t, log)
}

func (c *suiteContracts) SetupSuite() {
	var err error
	c.para, c.sim, c.x2EthContracts, c.x2EthDeployInfo, err = DeployContracts()
	require.Nil(c.T(), err)
}

func (c *suiteContracts) Test_GetOperator() {
	operator, err := GetOperator(c.sim, c.para.InitValidators[0], c.x2EthDeployInfo.BridgeBank.Address)
	require.Nil(c.T(), err)
	assert.Equal(c.T(), operator.String(), c.para.Operator.String())
}

func (c *suiteContracts) Test_IsActiveValidator() {
	bret, err := IsActiveValidator(c.para.InitValidators[0], c.x2EthContracts.Valset)
	require.Nil(c.T(), err)
	assert.Equal(c.T(), bret, true)

	key, _ := crypto.GenerateKey()
	addr := crypto.PubkeyToAddress(key.PublicKey)
	bret, err = IsActiveValidator(addr, c.x2EthContracts.Valset)
	require.Nil(c.T(), err) // ???
	assert.Equal(c.T(), bret, false)
}

func (c *suiteContracts) Test_IsProphecyPending() {
	claimID := crypto.Keccak256Hash(big.NewInt(50).Bytes())
	bret, err := IsProphecyPending(claimID, c.para.InitValidators[0], c.x2EthContracts.DplatformBridge)
	require.Nil(c.T(), err)
	assert.Equal(c.T(), bret, false)
}

func (c *suiteContracts) Test_LogLockToEthBridgeClaim() {
	to := common.FromHex(dplatformAddr)
	event := &events.LockEvent{
		From:   c.para.InitValidators[0],
		To:     to,
		Token:  common.HexToAddress(ethTokenAddr),
		Symbol: "eth",
		Value:  big.NewInt(10000 * 10000 * 10000),
		Nonce:  big.NewInt(1),
	}
	witnessClaim, err := LogLockToEthBridgeClaim(event, 1, c.x2EthDeployInfo.BridgeBank.Address.String(), 18)
	require.Nil(c.T(), err)
	assert.NotEmpty(c.T(), witnessClaim)
	assert.Equal(c.T(), witnessClaim.EthereumChainID, int64(1))
	assert.Equal(c.T(), witnessClaim.BridgeBrankAddr, c.x2EthDeployInfo.BridgeBank.Address.String())
	assert.Equal(c.T(), witnessClaim.TokenAddr, ethTokenAddr)
	assert.Equal(c.T(), witnessClaim.Symbol, event.Symbol)
	assert.Equal(c.T(), witnessClaim.EthereumSender, event.From.String())
	assert.Equal(c.T(), witnessClaim.DplatformReceiver, string(event.To))
	assert.Equal(c.T(), witnessClaim.Amount, "100")
	assert.Equal(c.T(), witnessClaim.Nonce, event.Nonce.Int64())
	assert.Equal(c.T(), witnessClaim.Decimal, int64(18))

	event.Token = common.HexToAddress("0x0000000000000000000000000000000000000001")
	_, err = LogLockToEthBridgeClaim(event, 1, c.x2EthDeployInfo.BridgeBank.Address.String(), 18)
	require.NotNil(c.T(), err)
	assert.Equal(c.T(), err, ebrelayerTypes.ErrAddress4Eth)
}

func (c *suiteContracts) Test_LogBurnToEthBridgeClaim() {
	to := common.FromHex(dplatformAddr)
	event := &events.BurnEvent{
		OwnerFrom:       c.para.InitValidators[0],
		DplatformReceiver: to,
		Token:           common.HexToAddress(ethTokenAddr),
		Symbol:          "bty",
		Amount:          big.NewInt(100),
		Nonce:           big.NewInt(2),
	}
	witnessClaim, err := LogBurnToEthBridgeClaim(event, 1, c.x2EthDeployInfo.BridgeBank.Address.String(), 8)
	require.Nil(c.T(), err)
	assert.NotEmpty(c.T(), witnessClaim)
	assert.Equal(c.T(), witnessClaim.EthereumChainID, int64(1))
	assert.Equal(c.T(), witnessClaim.BridgeBrankAddr, c.x2EthDeployInfo.BridgeBank.Address.String())
	assert.Equal(c.T(), witnessClaim.TokenAddr, ethTokenAddr)
	assert.Equal(c.T(), witnessClaim.Symbol, event.Symbol)
	assert.Equal(c.T(), witnessClaim.EthereumSender, event.OwnerFrom.String())
	assert.Equal(c.T(), witnessClaim.DplatformReceiver, string(event.DplatformReceiver))
	assert.Equal(c.T(), witnessClaim.Amount, "100")
	assert.Equal(c.T(), witnessClaim.Nonce, event.Nonce.Int64())
	assert.Equal(c.T(), witnessClaim.Decimal, int64(8))
}

func (c *suiteContracts) Test_ParseBurnLockTxReceipt_DplatformMsgToProphecyClaim() {
	claimType := events.MsgBurn
	dplatformToEth := types.ReceiptDplatformToEth{
		DplatformSender:    dplatformAddr,
		EthereumReceiver: ethAddr,
		TokenContract:    ethTokenAddr,
		IssuerDotSymbol:  "bty",
		Amount:           "100",
		Decimals:         8,
	}

	log := &dplatformTypes.ReceiptLog{
		Ty:  types.TyWithdrawDplatformLog,
		Log: dplatformTypes.Encode(&dplatformToEth),
	}

	var logs []*dplatformTypes.ReceiptLog
	logs = append(logs, log)

	receipt := &dplatformTypes.ReceiptData{
		Ty:   types.TyWithdrawDplatformLog,
		Logs: logs,
	}

	dplatformMsg := ParseBurnLockTxReceipt(claimType, receipt)
	require.NotNil(c.T(), dplatformMsg)
	assert.Equal(c.T(), dplatformMsg.ClaimType, claimType)
	assert.Equal(c.T(), dplatformMsg.DplatformSender, []byte(dplatformToEth.DplatformSender))
	assert.Equal(c.T(), dplatformMsg.EthereumReceiver, common.HexToAddress(dplatformToEth.EthereumReceiver))
	assert.Equal(c.T(), dplatformMsg.TokenContractAddress, common.HexToAddress(dplatformToEth.TokenContract))
	assert.Equal(c.T(), dplatformMsg.Symbol, dplatformToEth.IssuerDotSymbol)
	assert.Equal(c.T(), dplatformMsg.Amount.String(), "100")

	prophecyClaim := DplatformMsgToProphecyClaim(*dplatformMsg)
	assert.Equal(c.T(), dplatformMsg.ClaimType, prophecyClaim.ClaimType)
	assert.Equal(c.T(), dplatformMsg.DplatformSender, prophecyClaim.DplatformSender)
	assert.Equal(c.T(), dplatformMsg.EthereumReceiver, prophecyClaim.EthereumReceiver)
	assert.Equal(c.T(), dplatformMsg.TokenContractAddress, prophecyClaim.TokenContractAddress)
	assert.Equal(c.T(), strings.ToLower(dplatformMsg.Symbol), prophecyClaim.Symbol)
	assert.Equal(c.T(), dplatformMsg.Amount, prophecyClaim.Amount)
}

func (c *suiteContracts) Test_RecoverContractHandler() {
	_, _, err := RecoverContractHandler(c.sim, c.x2EthDeployInfo.BridgeRegistry.Address, c.x2EthDeployInfo.BridgeRegistry.Address)
	require.Nil(c.T(), err)
}

func (c *suiteContracts) Test_RecoverOracleInstance() {
	oracleInstance, err := RecoverOracleInstance(c.sim, c.x2EthDeployInfo.BridgeRegistry.Address, c.x2EthDeployInfo.BridgeRegistry.Address)
	require.Nil(c.T(), err)
	require.NotNil(c.T(), oracleInstance)
}

func (c *suiteContracts) Test_GetDeployHeight() {
	height, err := GetDeployHeight(c.sim, c.x2EthDeployInfo.BridgeRegistry.Address, c.x2EthDeployInfo.BridgeRegistry.Address)
	require.Nil(c.T(), err)
	assert.True(c.T(), height > 0)
}

func (c *suiteContracts) Test_CreateBridgeToken() {
	operatorInfo := &OperatorInfo{
		PrivateKey: c.para.DeployPrivateKey,
		Address:    crypto.PubkeyToAddress(c.para.DeployPrivateKey.PublicKey),
	}
	tokenAddr, err := CreateBridgeToken("bty", c.sim, operatorInfo, c.x2EthDeployInfo, c.x2EthContracts)
	require.Nil(c.T(), err)
	c.sim.Commit()

	addr, err := GetToken2address(c.x2EthContracts.BridgeBank, "bty")
	require.Nil(c.T(), err)
	assert.Equal(c.T(), addr, tokenAddr)

	dplatformSender := []byte("14KEKbYtKKQm4wMthSK9J4La4nAiidGozt")
	amount := int64(100)
	ethReceiver := c.para.InitValidators[2]
	claimID := crypto.Keccak256Hash(dplatformSender, ethReceiver.Bytes(), big.NewInt(amount).Bytes())
	authOracle, err := PrepareAuth(c.sim, c.para.ValidatorPriKey[0], c.para.InitValidators[0])
	require.Nil(c.T(), err)
	signature, err := SignClaim4Eth(claimID, c.para.ValidatorPriKey[0])
	require.Nil(c.T(), err)

	_, err = c.x2EthContracts.Oracle.NewOracleClaim(
		authOracle,
		events.ClaimTypeLock,
		dplatformSender,
		ethReceiver,
		common.HexToAddress(tokenAddr),
		"bty",
		big.NewInt(amount),
		claimID,
		signature)
	require.Nil(c.T(), err)
	c.sim.Commit()

	balanceNew, err := GetBalance(c.sim, tokenAddr, ethReceiver.String())
	require.Nil(c.T(), err)
	require.Equal(c.T(), balanceNew, "100")

	dplatformReceiver := "1GTxrmuWiXavhcvsaH5w9whgVxUrWsUMdV"
	{
		amount := "10"
		bn := big.NewInt(1)
		bn, _ = bn.SetString(x2ethTypes.TrimZeroAndDot(amount), 10)
		txhash, err := Burn(hexutil.Encode(crypto.FromECDSA(c.para.ValidatorPriKey[2])), tokenAddr, dplatformReceiver, c.x2EthDeployInfo.BridgeBank.Address, bn, c.x2EthContracts.BridgeBank, c.sim)
		require.NoError(c.T(), err)
		c.sim.Commit()

		balanceNew, err = GetBalance(c.sim, tokenAddr, ethReceiver.String())
		require.Nil(c.T(), err)
		require.Equal(c.T(), balanceNew, "90")

		status := GetEthTxStatus(c.sim, common.HexToHash(txhash))
		fmt.Println()
		fmt.Println(status)
	}

	{
		amount := "10"
		bn := big.NewInt(1)
		bn, _ = bn.SetString(x2ethTypes.TrimZeroAndDot(amount), 10)
		_, err := ApproveAllowance(hexutil.Encode(crypto.FromECDSA(c.para.ValidatorPriKey[2])), tokenAddr, c.x2EthDeployInfo.BridgeBank.Address, bn, c.sim)
		require.Nil(c.T(), err)
		c.sim.Commit()

		_, err = BurnAsync(hexutil.Encode(crypto.FromECDSA(c.para.ValidatorPriKey[2])), tokenAddr, dplatformReceiver, bn, c.x2EthContracts.BridgeBank, c.sim)
		require.Nil(c.T(), err)
		c.sim.Commit()

		balanceNew, err = GetBalance(c.sim, tokenAddr, ethReceiver.String())
		require.Nil(c.T(), err)
		require.Equal(c.T(), balanceNew, "80")
	}
}

func (c *suiteContracts) Test_CreateERC20Token() {
	operatorInfo := &OperatorInfo{
		PrivateKey: c.para.DeployPrivateKey,
		Address:    crypto.PubkeyToAddress(c.para.DeployPrivateKey.PublicKey),
	}
	tokenAddr, err := CreateERC20Token("testc", c.sim, operatorInfo)
	require.Nil(c.T(), err)
	c.sim.Commit()

	amount := "10000000000000"
	bn := big.NewInt(1)
	bn, _ = bn.SetString(x2ethTypes.TrimZeroAndDot(amount), 10)

	_, err = MintERC20Token(tokenAddr, c.para.Deployer.String(), bn, c.sim, operatorInfo)
	require.Nil(c.T(), err)
	c.sim.Commit()

	balance, err := GetDepositFunds(c.sim, tokenAddr)
	require.Nil(c.T(), err)
	assert.Equal(c.T(), balance, amount)

	amount = "100"
	bn = big.NewInt(1)
	bn, _ = bn.SetString(x2ethTypes.TrimZeroAndDot(amount), 10)
	txhash, err := TransferToken(tokenAddr, hexutil.Encode(crypto.FromECDSA(c.para.DeployPrivateKey)), c.para.InitValidators[0].String(), bn, c.sim)
	require.Nil(c.T(), err)
	c.sim.Commit()

	_, err = c.sim.TransactionReceipt(context.Background(), common.HexToHash(txhash))
	require.Nil(c.T(), err)
	balance, err = GetBalance(c.sim, tokenAddr, c.para.InitValidators[0].String())
	require.Nil(c.T(), err)
	assert.Equal(c.T(), balance, amount)

	{
		amount = "100"
		bn := big.NewInt(1)
		bn, _ = bn.SetString(x2ethTypes.TrimZeroAndDot(amount), 10)
		dplatformReceiver := "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
		_, err = LockEthErc20Asset(hexutil.Encode(crypto.FromECDSA(c.para.DeployPrivateKey)), tokenAddr, dplatformReceiver, bn, c.sim, c.x2EthContracts.BridgeBank, c.x2EthDeployInfo.BridgeBank.Address)
		require.Nil(c.T(), err)
		c.sim.Commit()

		balance, err = GetBalance(c.sim, tokenAddr, c.para.Deployer.String())
		require.Nil(c.T(), err)
		fmt.Println(balance)
		assert.Equal(c.T(), balance, "9999999999800")
	}

	{
		amount := "800"
		bn := big.NewInt(1)
		bn, _ = bn.SetString(x2ethTypes.TrimZeroAndDot(amount), 10)
		_, err = ApproveAllowance(hexutil.Encode(crypto.FromECDSA(c.para.DeployPrivateKey)), tokenAddr, c.x2EthDeployInfo.BridgeBank.Address, bn, c.sim)
		require.Nil(c.T(), err)
		c.sim.Commit()

		dplatformReceiver := "14KEKbYtKKQm4wMthSK9J4La4nAiidGozt"
		_, err = LockEthErc20AssetAsync(hexutil.Encode(crypto.FromECDSA(c.para.DeployPrivateKey)), tokenAddr, dplatformReceiver, bn, c.sim, c.x2EthContracts.BridgeBank)
		require.Nil(c.T(), err)
		c.sim.Commit()

		balance, err = GetBalance(c.sim, tokenAddr, c.para.Deployer.String())
		require.Nil(c.T(), err)
		fmt.Println(balance)
		assert.Equal(c.T(), balance, "9999999999000")
	}
}

func (c *suiteContracts) Test_GetLockedFunds() {
	balance, err := GetLockedFunds(c.x2EthContracts.BridgeBank, "")
	require.Nil(c.T(), err)
	assert.Equal(c.T(), balance, "0")
}

func PrepareTestEnv() (*ethinterface.SimExtend, *DeployPara) {
	genesiskey, _ := crypto.GenerateKey()
	alloc := make(core.GenesisAlloc)
	genesisAddr := crypto.PubkeyToAddress(genesiskey.PublicKey)
	genesisAccount := core.GenesisAccount{
		Balance:    big.NewInt(10000000000 * 10000),
		PrivateKey: crypto.FromECDSA(genesiskey),
	}
	alloc[genesisAddr] = genesisAccount

	var InitValidators []common.Address
	var ValidatorPriKey []*ecdsa.PrivateKey
	for i := 0; i < 4; i++ {
		key, _ := crypto.GenerateKey()
		addr := crypto.PubkeyToAddress(key.PublicKey)
		InitValidators = append(InitValidators, addr)
		ValidatorPriKey = append(ValidatorPriKey, key)

		account := core.GenesisAccount{
			Balance:    big.NewInt(100000000 * 100),
			PrivateKey: crypto.FromECDSA(key),
		}
		alloc[addr] = account
	}
	gasLimit := uint64(100000000)
	sim := new(ethinterface.SimExtend)
	sim.SimulatedBackend = backends.NewSimulatedBackend(alloc, gasLimit)

	InitPowers := []*big.Int{big.NewInt(80), big.NewInt(10), big.NewInt(10), big.NewInt(10)}
	para := &DeployPara{
		DeployPrivateKey: genesiskey,
		Deployer:         genesisAddr,
		Operator:         genesisAddr,
		InitValidators:   InitValidators,
		ValidatorPriKey:  ValidatorPriKey,
		InitPowers:       InitPowers,
	}

	return sim, para
}

func DeployContracts() (*DeployPara, *ethinterface.SimExtend, *X2EthContracts, *X2EthDeployInfo, error) {
	ctx := context.Background()
	sim, para := PrepareTestEnv()

	callMsg := ethereum.CallMsg{
		From: para.Deployer,
		Data: common.FromHex(generated.BridgeBankBin),
	}

	_, err := sim.EstimateGas(ctx, callMsg)
	if nil != err {
		panic("failed to estimate gas due to:" + err.Error())
	}
	x2EthContracts, x2EthDeployInfo, err := DeployAndInit(sim, para)
	if nil != err {
		return nil, nil, nil, nil, err
	}
	sim.Commit()

	return para, sim, x2EthContracts, x2EthDeployInfo, nil
}
