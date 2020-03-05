package executor

import "github.com/33cn/plugin/plugin/dapp/x2ethereum/types"

// Ethbridge module event types
var (
	EventTypeCreateClaim    = "create_claim"
	EventTypeProphecyStatus = "prophecy_status"
	EventTypeBurn           = "burn"
	EventTypeLock           = "lock"

	AttributeKeyEthereumSender  = "ethereum_sender"
	AttributeKeyChain33Receiver = "chain33_receiver"
	AttributeKeyAmount          = "amount"
	AttributeKeyStatus          = "status"
	AttributeKeyClaimType       = "claim_type"

	AttributeKeyEthereumChainID  = "ethereum_chain_id"
	AttributeKeyTokenContract    = "token_contract_address"
	AttributeKeyChain33Sender    = "chain33_sender"
	AttributeKeyEthereumReceiver = "ethereum_receiver"

	AttributeValueCategory = types.ModuleName
)
