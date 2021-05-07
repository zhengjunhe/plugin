package main

import (
	"context"
	"fmt"
	"math/big"
	"strings"
	"time"

	ebTypes "github.com/33cn/plugin/plugin/dapp/cross2eth/ebrelayer/types"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/crypto"

	"github.com/33cn/chain33/rpc/jsonclient"
	rpctypes "github.com/33cn/chain33/rpc/types"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/abi"
	"github.com/33cn/plugin/plugin/dapp/evm/executor/vm/common"
	"github.com/spf13/cobra"
)

//TokenAddressCmd...
func MultiSignCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "multisign",
		Short: "deploy,setup and trasfer multisign",
		Args:  cobra.MinimumNArgs(1),
	}
	cmd.AddCommand(
		DeployMultiSignCmd(),
		SetupCmd(),
		TransferCmd(),
	)
	return cmd
}

func DeployMultiSignCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "deploy",
		Short: "deploy mulsign to chain33",
		Run:   DeployMultiSign,
	}
	return cmd
}

func DeployMultiSign(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	var res rpctypes.Reply
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Manager.DeployMulsign2Chain33", nil, &res)
	ctx.Run()
}

func SetupCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "Setup",
		Short: "Setup owners to contract",
		Run:   SetupOwner,
	}
	SetupOwnerFlags(cmd)
	return cmd
}

func SetupOwnerFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("owner", "o", "", "owners's address, separated by ','")
	_ = cmd.MarkFlagRequired("owner")
	cmd.Flags().StringP("operator", "k", "", "operator address")
	_ = cmd.MarkFlagRequired("operator")

}

func SetupOwner(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	ownersStr, _ := cmd.Flags().GetString("owner")
	operator, _ := cmd.Flags().GetString("operator")
	owners := strings.Split(ownersStr, ",")

	para := ebTypes.SetupMulSign{
		Operator: operator,
		Owners:   owners,
	}
	var res rpctypes.Reply
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Manager.SetupOwner4Chain33", para, &res)
	ctx.Run()
}

func TransferCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "transfer",
		Short: "transfer via safe",
		Run:   Transfer,
	}
	TransferFlags(cmd)
	return cmd
}

func TransferFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("receiver", "r", "", "receive address")
	_ = cmd.MarkFlagRequired("receiver")

	cmd.Flags().Float64P("amount", "a", 0, "amount to transfer")
	_ = cmd.MarkFlagRequired("amount")

	cmd.Flags().StringP("token", "t", "", "erc20 address")
}

func Transfer(cmd *cobra.Command, args []string) {
	rpcLaddr, _ := cmd.Flags().GetString("rpc_laddr")
	receiver, _ := cmd.Flags().GetString("receiver")
	tokenAddr, _ := cmd.Flags().GetString("token")
	amount, _ := cmd.Flags().GetFloat64("amount")

	para := ebTypes.SafeTransfer{
		To:     receiver,
		Token:  tokenAddr,
		Amount: amount,
	}
	var res rpctypes.Reply
	ctx := jsonclient.NewRPCCtx(rpcLaddr, "Manager.SafeTransfer4Chain33", para, &res)
	ctx.Run()
}

//threshold int, to, paymentToken, paymentReceiver string, payment int64
func SetupOwnerProc(safe string) error {
	owners := []string{"0x0f2e821517D4f64a012a04b668a6b1aa3B262e08", "0xee760B2E502244016ADeD3491948220B3b1dd789", "0x21B5f4C2F6Ff418fa0067629D9D76AE03fB4a2d2"}
	_ = recoverBinancePrivateKey()
	auth, err := PrepareAuth(privateKey, deployerAddr, GasLimitTxExec)
	if nil != err {
		return err
	}

	gnosisSafeAddr := common.HexToAddress(safe)
	gnosisSafeInt, err := gnosisSafe.NewGnosisSafe(gnosisSafeAddr, ethClient)
	if nil != err {
		return err
	}

	//_owners []common.Address, _threshold *big.Int, to common.Address, data []byte,
	// fallbackHandler common.Address, paymentToken common.Address,
	// payment *big.Int, paymentReceiver common.Address
	var _owners []common.Address
	for _, onwer := range owners {
		_owners = append(_owners, common.HexToAddress(onwer))
	}
	AddressZero := common.HexToAddress("0x0000000000000000000000000000000000000000")

	//safe.setup([user1.address, user2.address], 1, AddressZero, "0x", handler.address, AddressZero, 0, AddressZero)
	setupTx, err := gnosisSafeInt.Setup(auth, _owners, big.NewInt(int64(len(_owners))), AddressZero, []byte{'0', 'x'}, AddressZero, AddressZero, big.NewInt(int64(0)), AddressZero)
	if nil != err {
		panic(fmt.Sprintf("Failed to setupTx with err:%s", err.Error()))
		return err
	}

	{
		fmt.Println("\nsetupTx tx hash:", setupTx.Hash().String())
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("setupTx timeout")
			case <-oneSecondtimeout.C:
				_, err := ethClient.TransactionReceipt(context.Background(), setupTx.Hash())
				if err == ethereum.NotFound {
					fmt.Println("\n No receipt received yet for setupTx  and continue to wait")
					continue
				} else if err != nil {
					panic("SetupOwner failed due to" + err.Error())
				}
				fmt.Println("\n Succeed to setup Tx")
				return nil
			}
		}
	}

	return nil
}

func TransferProc(safe, to, token string, fValue float64) error {
	owners := []string{"0x0f2e821517D4f64a012a04b668a6b1aa3B262e08",
		"0xee760B2E502244016ADeD3491948220B3b1dd789",
		"0x21B5f4C2F6Ff418fa0067629D9D76AE03fB4a2d2"}
	_ = recoverBinancePrivateKey()
	auth, err := PrepareAuth(privateKey, deployerAddr, GasLimitTxExec)
	if nil != err {
		return err
	}

	gnosisSafeAddr := common.HexToAddress(safe)
	gnosisSafeInt, err := gnosisSafe.NewGnosisSafe(gnosisSafeAddr, ethClient)
	if nil != err {
		return err
	}
	AddressZero := common.HexToAddress("0x0000000000000000000000000000000000000000")

	//_owners []common.Address, _threshold *big.Int, to common.Address, data []byte,
	// fallbackHandler common.Address, paymentToken common.Address,
	// payment *big.Int, paymentReceiver common.Address
	var _owners []common.Address
	for _, onwer := range owners {
		_owners = append(_owners, common.HexToAddress(onwer))
	}

	//opts *bind.TransactOpts, to common.Address, value *big.Int, data []byte, operation uint8, safeTxGas *big.Int, baseGas *big.Int, gasPrice *big.Int, gasToken common.Address, refundReceiver common.Address, signatures []byte

	_to := common.HexToAddress(to)
	_data := []byte{'0', 'x'}
	safeTxGas := big.NewInt(10 * 10000)
	baseGas := big.NewInt(0)
	gasPrice := big.NewInt(0)
	var value *big.Int = big.NewInt(int64(fValue * 1e18))
	opts := &bind.CallOpts{
		From:    deployerAddr,
		Context: context.Background(),
	}
	//token transfer
	if token != "" {
		_to = common.HexToAddress(token)

		erc20Abi, err := abi.JSON(strings.NewReader(erc20.ERC20ABI))
		if err != nil {
			return err
		}

		tokenInstance, err := erc20.NewERC20(_to, ethClient)
		if err != nil {
			return err
		}
		decimals, err := tokenInstance.Decimals(opts)
		if err != nil {
			return err
		}
		mul := int64(1)
		for i := 0; i < int(decimals); i++ {
			mul *= 10
		}
		value = big.NewInt(int64(fValue * float64(mul)))

		//const data = token.interface.encodeFunctionData("transfer", [address, 500])
		_data, err = erc20Abi.Pack("transfer", common.HexToAddress(to), value)
		if err != nil {
			return err
		}
		//对于erc20这种方式 最后需要将其设置为0
		value = big.NewInt(0)
	}

	nonce, err := gnosisSafeInt.Nonce(opts)
	if err != nil {
		panic("Failed to get Nonce")
		return err
	}

	//opts *bind.TransactOpts, to common.Address, value *big.Int, data []byte, operation uint8, safeTxGas *big.Int, baseGas *big.Int, gasPrice *big.Int, gasToken common.Address, refundReceiver common.Address, _nonce *big.Int
	signContent, err := gnosisSafeInt.GetTransactionHash(opts, _to, value, _data, 0,
		safeTxGas, baseGas, gasPrice, AddressZero, AddressZero, nonce)
	if err != nil {
		panic("Failed to GetTransactionHash")
		return err
	}
	fmt.Println("safe.Nonce =", nonce.String(), "safe.Nonce(int64) =", nonce.Int64())
	sigs := buildSigs(signContent[:])

	execTx, err := gnosisSafeInt.ExecTransaction(auth, _to, value, _data, 0,
		safeTxGas, baseGas, gasPrice, AddressZero, AddressZero, sigs)
	if nil != err {
		panic(fmt.Sprintf("Failed to ExecTransaction with err:%s", err.Error()))
		return err
	}

	{
		fmt.Println("\nExecTransaction tx hash:", execTx.Hash().String())
		timeout := time.NewTimer(300 * time.Second)
		oneSecondtimeout := time.NewTicker(5 * time.Second)
		for {
			select {
			case <-timeout.C:
				panic("ExecTransaction timeout")
			case <-oneSecondtimeout.C:
				_, err := ethClient.TransactionReceipt(context.Background(), execTx.Hash())
				if err == ethereum.NotFound {
					fmt.Println("\n No receipt received yet for ExecTransaction  and continue to wait")
					continue
				} else if err != nil {
					panic("ExecTransaction failed due to" + err.Error())
				}
				fmt.Println("\n Succeed to ExecTransaction Tx")
				return nil
			}
		}
	}

	return nil
}

func buildSigs(data []byte) (sigs []byte) {
	fmt.Println("\nbuildSigs, data:", common.Bytes2Hex(data))

	for _, privateKeyStr := range privateKeyStrs {
		privateKey, err := crypto.ToECDSA(common.FromHex(privateKeyStr))
		if nil != err {
			panic("Failed to recover private key")
			return nil
		}

		signature, err := crypto.Sign(data, privateKey)
		if err != nil {
			panic("Failed to sign data due to:" + err.Error())
			return nil
		}
		signature[64] += 27
		sigs = append(sigs, signature[:]...)
	}
	return
}
