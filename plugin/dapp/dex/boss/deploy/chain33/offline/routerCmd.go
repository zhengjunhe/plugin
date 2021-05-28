package offline

import (
	"math/rand"
	"time"

	"github.com/33cn/plugin/plugin/dapp/dex/contracts/pancake-swap-periphery/src/pancakeFactory"
	"github.com/33cn/plugin/plugin/dapp/dex/contracts/pancake-swap-periphery/src/pancakeRouter"

	"github.com/33cn/chain33/system/crypto/secp256k1"

	"github.com/33cn/chain33/common"
	erc20 "github.com/33cn/plugin/plugin/dapp/cross2eth/contracts/erc20/generated"

	"github.com/33cn/chain33/types"
	"github.com/golang/protobuf/proto"
	"github.com/spf13/cobra"
)

// 创建ERC20合约
func createERC20ContractCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "erc20",
		Short: "create ERC20 contract",
		Run:   createERC20Contract,
	}
	addCreateERC20ContractFlags(cmd)
	return cmd
}

func addCreateERC20ContractFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("caller", "c", "", "the caller's private key")
	cmd.MarkFlagRequired("caller")
	cmd.Flags().StringP("name", "a", "", "REC20 name")
	cmd.MarkFlagRequired("name")
	cmd.Flags().StringP("symbol", "s", "", "REC20 symbol")
	cmd.MarkFlagRequired("symbol")
	cmd.Flags().StringP("supply", "m", "", "REC20 supply")
	cmd.MarkFlagRequired("supply")

	cmd.Flags().StringP("note", "n", "", "transaction note info (optional)")
	cmd.Flags().Float64P("fee", "f", 0, "contract gas fee (optional)")
}

func createERC20Contract(cmd *cobra.Command, args []string) {
	caller, _ := cmd.Flags().GetString("caller")
	name, _ := cmd.Flags().GetString("name")
	symbol, _ := cmd.Flags().GetString("symbol")
	supply, _ := cmd.Flags().GetString("supply")

	privateKey, _ := cmd.Flags().GetString("caller")
	expire, _ := cmd.Flags().GetString("expire")
	note, _ := cmd.Flags().GetString("note")
	fee, _ := cmd.Flags().GetFloat64("fee")
	paraName, _ := cmd.Flags().GetString("paraName")
	chainID, _ := cmd.Flags().GetInt32("chainID")
	feeInt64 := int64(fee*1e4) * 1e4
	info := &TxCreateInfo{
		privateKey: privateKey,
		expire:     expire,
		note:       note,
		fee:        feeInt64,
		paraName:   paraName,
		chainID:    chainID,
	}
	createPara := name + "," + symbol + "," + supply + "," + caller
	content, err := createContractAndSign(info, erc20.ERC20Bin, erc20.ERC20ABI, createPara, "erc20")
	if nil != err {
		writeContractFile("./erc20", content)
	}
}

func createAndSignEvmTx(chainID int32, action proto.Message, execer, privateKeyStr, contract2call, expire string, fee int64) (string, error) {
	tx := &types.Transaction{Execer: []byte(execer), Payload: types.Encode(action), Fee: 0, To: contract2call}

	expireInt64, err := types.ParseExpire(expire)
	if nil != err {
		return "", err
	}

	if expireInt64 > types.ExpireBound {
		if expireInt64 < int64(time.Second*120) {
			expireInt64 = int64(time.Second * 120)
		}
		//用秒数来表示的时间
		tx.Expire = types.Now().Unix() + expireInt64/int64(time.Second)
	} else {
		tx.Expire = expireInt64
	}

	tx.Fee = int64(1e7)
	if tx.Fee < fee {
		tx.Fee += fee
	}

	random := rand.New(rand.NewSource(time.Now().UnixNano()))
	tx.Nonce = random.Int63()
	tx.ChainID = chainID

	var driver secp256k1.Driver
	privateKeySli, err := common.FromHex(privateKeyStr)
	if nil != err {
		return "", err
	}
	privateKey, err := driver.PrivKeyFromBytes(privateKeySli)
	if nil != err {
		return "", err
	}

	tx.Sign(types.SECP256K1, privateKey)
	txData := types.Encode(tx)
	dataStr := common.ToHex(txData)

	return dataStr, nil
}

func createRouterCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "router",
		Short: "create and sign offline router contract",
		Run:   createRouterContract,
	}
	addCreateRouterFlags(cmd)
	return cmd
}

func addCreateRouterFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("caller", "c", "", "the caller address")
	cmd.MarkFlagRequired("caller")

	cmd.Flags().StringP("expire", "", "120s", "transaction expire time (optional)")
	cmd.Flags().StringP("note", "n", "", "transaction note info (optional)")
	cmd.Flags().Float64P("fee", "f", 0, "contract gas fee (optional)")

	cmd.Flags().StringP("factory", "t", "", "address of factory")
	cmd.MarkFlagRequired("factory")
	cmd.Flags().StringP("weth9", "w", "", "address of weth9")
	cmd.MarkFlagRequired("weth9")

}

func createRouterContract(cmd *cobra.Command, args []string) {
	factory, _ := cmd.Flags().GetString("factory")
	weth9, _ := cmd.Flags().GetString("weth9")

	privateKey, _ := cmd.Flags().GetString("caller")
	expire, _ := cmd.Flags().GetString("expire")
	note, _ := cmd.Flags().GetString("note")
	fee, _ := cmd.Flags().GetFloat64("fee")
	paraName, _ := cmd.Flags().GetString("paraName")
	chainID, _ := cmd.Flags().GetInt32("chainID")
	feeInt64 := int64(fee*1e4) * 1e4
	info := &TxCreateInfo{
		privateKey: privateKey,
		expire:     expire,
		note:       note,
		fee:        feeInt64,
		paraName:   paraName,
		chainID:    chainID,
	}
	//constructor(address _factory, address _WETH)
	createPara := factory + "," + weth9
	content, err := createContractAndSign(info, pancakeRouter.PancakeRouterBin, pancakeRouter.PancakeRouterABI, createPara, "pancakeRouter")
	if nil != err {
		writeContractFile("./pancakeRouter", content)
	}
}

func createWeth9Cmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "weth9",
		Short: "create and sign offline weth9 contract",
		Run:   createWeth9,
	}
	addCreateWeth9Flags(cmd)
	return cmd
}

func createWeth9(cmd *cobra.Command, args []string) {
	privateKey, _ := cmd.Flags().GetString("caller")
	expire, _ := cmd.Flags().GetString("expire")
	note, _ := cmd.Flags().GetString("note")
	fee, _ := cmd.Flags().GetFloat64("fee")
	paraName, _ := cmd.Flags().GetString("paraName")
	chainID, _ := cmd.Flags().GetInt32("chainID")
	feeInt64 := int64(fee*1e4) * 1e4
	info := &TxCreateInfo{
		privateKey: privateKey,
		expire:     expire,
		note:       note,
		fee:        feeInt64,
		paraName:   paraName,
		chainID:    chainID,
	}
	//constructor(address _feeToSetter) public
	createPara := ""
	content, err := createContractAndSign(info, pancakeRouter.WETH9Bin, pancakeRouter.WETH9ABI, createPara, "WETH9")
	if nil != err {
		writeContractFile("./weth9", content)
	}
}

func addCreateWeth9Flags(cmd *cobra.Command) {
	cmd.Flags().StringP("caller", "c", "", "the caller address")
	cmd.MarkFlagRequired("caller")

	cmd.Flags().StringP("expire", "", "120s", "transaction expire time (optional)")
	cmd.Flags().StringP("note", "n", "", "transaction note info (optional)")
	cmd.Flags().Float64P("fee", "f", 0, "contract gas fee (optional)")
}

func createFactoryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "factory",
		Short: "create and sign offline factory contract",
		Run:   createFactoryContract,
	}
	addCreateFactoryContractFlags(cmd)
	return cmd
}

func createFactoryContract(cmd *cobra.Command, args []string) {
	feeToSetter, _ := cmd.Flags().GetString("feeToSetter")

	privateKey, _ := cmd.Flags().GetString("caller")
	expire, _ := cmd.Flags().GetString("expire")
	note, _ := cmd.Flags().GetString("note")
	fee, _ := cmd.Flags().GetFloat64("fee")
	paraName, _ := cmd.Flags().GetString("paraName")
	chainID, _ := cmd.Flags().GetInt32("chainID")
	feeInt64 := int64(fee*1e4) * 1e4
	info := &TxCreateInfo{
		privateKey: privateKey,
		expire:     expire,
		note:       note,
		fee:        feeInt64,
		paraName:   paraName,
		chainID:    chainID,
	}
	//constructor(address _feeToSetter) public
	createPara := feeToSetter
	content, err := createContractAndSign(info, pancakeFactory.PancakeFactoryBin, pancakeFactory.PancakeFactoryABI, createPara, "PancakeFactory")
	if nil != err {
		writeContractFile("./factory", content)
	}
}

func addCreateFactoryContractFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("caller", "c", "", "the caller address")
	cmd.MarkFlagRequired("caller")

	cmd.Flags().StringP("expire", "", "120s", "transaction expire time (optional)")
	cmd.Flags().StringP("note", "n", "", "transaction note info (optional)")
	cmd.Flags().Float64P("fee", "f", 0, "contract gas fee (optional)")

	cmd.Flags().StringP("feeToSetter", "a", "", "address for fee to Setter")
	cmd.MarkFlagRequired("feeToSetter")
}
