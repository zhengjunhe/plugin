package executor

import (
	"encoding/json"
	"fmt"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/executor/common"
	"github.com/33cn/plugin/plugin/dapp/x2ethereum/types"
	"github.com/gogo/protobuf/codec"
	"strconv"
)

// QueryEthProphecyParams defines the params for the following queries:
// - 'custom/ethbridge/prophecies/'
type QueryEthProphecyParams struct {
	EthereumChainID       int               `json:"ethereum_chain_id"`
	BridgeContractAddress common.EthAddress `json:"bridge_contract_address"`
	Nonce                 int               `json:"nonce"`
	Symbol                string            `json:"symbol"`
	TokenContractAddress  common.EthAddress `json:"token_contract_address"`
	EthereumSender        common.EthAddress `json:"ethereum_sender"`
}

// QueryEthProphecyParams creates a new QueryEthProphecyParams
func NewQueryEthProphecyParams(ethereumChainID int, bridgeContractAddress string, nonce int, symbol string, tokenContractAddress string, ethereumSender string) QueryEthProphecyParams {
	return QueryEthProphecyParams{
		EthereumChainID:       ethereumChainID,
		BridgeContractAddress: NewEthAddressByProto(bridgeContractAddress),
		Nonce:                 nonce,
		Symbol:                symbol,
		TokenContractAddress:  NewEthAddressByProto(tokenContractAddress),
		EthereumSender:        NewEthAddressByProto(ethereumSender),
	}
}

// Query Result Payload for an eth prophecy query
type QueryEthProphecyResponse struct {
	ID     string                 `json:"id"`
	Status types.ProphecyStatus   `json:"status"`
	Claims []types.EthBridgeClaim `json:"claims"`
}

func NewQueryEthProphecyResponse(id string, status types.ProphecyStatus, claims []types.EthBridgeClaim) QueryEthProphecyResponse {
	return QueryEthProphecyResponse{
		ID:     id,
		Status: status,
		Claims: claims,
	}
}

func (response QueryEthProphecyResponse) String() string {
	prophecyJSON, err := json.Marshal(response)
	if err != nil {
		return fmt.Sprintf("Error marshalling json: %v", err)
	}

	return string(prophecyJSON)
}

// query endpoints supported by the oracle Querier
const (
	QueryEthProphecy = "prophecies"
)

// NewQuerier is the module level router for state queries
func NewQuerier(keeper types.OracleKeeper, cdc *codec.Codec, codespace sdk.CodespaceType) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) (res []byte, err error) {
		switch path[0] {
		case QueryEthProphecy:
			return queryEthProphecy(ctx, cdc, req, keeper, codespace)
		default:
			return nil, sdk.ErrUnknownRequest("unknown ethbridge query endpoint")
		}
	}
}

func queryEthProphecy(ctx sdk.Context, cdc *codec.Codec, req abci.RequestQuery, keeper types.OracleKeeper, codespace sdk.CodespaceType) (res []byte, errSdk sdk.Error) {
	var params types.QueryEthProphecyParams

	if err := cdc.UnmarshalJSON(req.Data, &params); err != nil {
		return []byte{}, sdk.ErrInternal(sdk.AppendMsgToErr("failed to parse params: %s", err.Error()))
	}
	id := strconv.Itoa(params.EthereumChainID) + strconv.Itoa(params.Nonce) + params.EthereumSender.String()
	prophecy, errSdk := keeper.GetProphecy(ctx, id)
	if errSdk != nil {
		return []byte{}, oracletypes.ErrProphecyNotFound(codespace)
	}

	bridgeClaims, errSdk := types.MapOracleClaimsToEthBridgeClaims(params.EthereumChainID, params.BridgeContractAddress, params.Nonce, params.Symbol, params.TokenContractAddress, params.EthereumSender, prophecy.ValidatorClaims, types.CreateEthClaimFromOracleString)
	if errSdk != nil {
		return []byte{}, errSdk
	}

	response := types.NewQueryEthProphecyResponse(prophecy.ID, prophecy.Status, bridgeClaims)

	bz, err := cdc.MarshalJSONIndent(response, "", "  ")
	if err != nil {
		panic(err)
	}

	return bz, nil
}
