package types

//Version4Relayer ...
const Version4Relayer = "0.1.0"

const (
	Chain33BlockChainName    = "Chain33-mainchain"
	EthereumBlockChainName   = "Ethereum-mainchain"
	BTYAddrChain33           = "1111111111111111111114oLvT2"
	NilAddrChain33           = "1111111111111111111114oLvT2"
	EthNilAddr               = "0x0000000000000000000000000000000000000000"
	SYMBOL_ETH               = "ETH"
	SYMBOL_BTY               = "BTY"
	Tx_Status_Pending        = "pending"
	Tx_Status_Success        = "Successful"
	Tx_Status_Failed         = "Failed"
	Source_Chain_Ethereum    = int32(0)
	Source_Chain_Chain33     = int32(1)
	Invalid_Tx_Index         = int64(-1)
	Invalid_Chain33Tx_Status = int32(-1)
)

var Tx_Status_Map = map[int32]string{
	1: Tx_Status_Pending,
	2: Tx_Status_Success,
	3: Tx_Status_Failed,
}
