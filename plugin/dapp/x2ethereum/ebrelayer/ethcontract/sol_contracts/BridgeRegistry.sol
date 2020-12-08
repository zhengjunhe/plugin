pragma solidity ^0.5.0;

contract BridgeRegistry {

    address public dplatformBridge;
    address public bridgeBank;
    address public oracle;
    address public valset;
    uint256 public deployHeight;

    event LogContractsRegistered(
        address _dplatformBridge,
        address _bridgeBank,
        address _oracle,
        address _valset
    );
    
    constructor(
        address _dplatformBridge,
        address _bridgeBank,
        address _oracle,
        address _valset
    )
        public
    {
        dplatformBridge = _dplatformBridge;
        bridgeBank = _bridgeBank;
        oracle = _oracle;
        valset = _valset;
        deployHeight = block.number;

        emit LogContractsRegistered(
            dplatformBridge,
            bridgeBank,
            oracle,
            valset
        );
    }
}