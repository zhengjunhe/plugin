pragma solidity ^0.5.0;

import "./DplatformBank.sol";
import "./EthereumBank.sol";
import "../Oracle.sol";
import "../DplatformBridge.sol";

/**
 * @title BridgeBank
 * @dev Bank contract which coordinates asset-related functionality.
 *      DplatformBank manages the minting and burning of tokens which
 *      represent Dplatform based assets, while EthereumBank manages
 *      the locking and unlocking of Ethereum and ERC20 token assets
 *      based on Ethereum.
 **/

contract BridgeBank is DplatformBank, EthereumBank {

    using SafeMath for uint256;
    
    address public operator;
    Oracle public oracle;
    DplatformBridge public dplatformBridge;

    /*
    * @dev: Constructor, sets operator
    */
    constructor (
        address _operatorAddress,
        address _oracleAddress,
        address _dplatformBridgeAddress
    )
        public
    {
        operator = _operatorAddress;
        oracle = Oracle(_oracleAddress);
        dplatformBridge = DplatformBridge(_dplatformBridgeAddress);
    }

    /*
    * @dev: Modifier to restrict access to operator
    */
    modifier onlyOperator() {
        require(
            msg.sender == operator,
            'Must be BridgeBank operator.'
        );
        _;
    }

    /*
    * @dev: Modifier to restrict access to the oracle
    */
    modifier onlyOracle()
    {
        require(
            msg.sender == address(oracle),
            "Access restricted to the oracle"
        );
        _;
    }

    /*
    * @dev: Modifier to restrict access to the dplatform bridge
    */
    modifier onlyDplatformBridge()
    {
        require(
            msg.sender == address(dplatformBridge),
            "Access restricted to the dplatform bridge"
        );
        _;
    }

   /*
    * @dev: Fallback function allows operator to send funds to the bank directly
    *       This feature is used for testing and is available at the operator's own risk.
    */
    function() external payable onlyOperator {}

    /*
    * @dev: Creates a new BridgeToken
    *
    * @param _symbol: The new BridgeToken's symbol
    * @return: The new BridgeToken contract's address
    */
    function createNewBridgeToken(
        string memory _symbol
    )
        public
        onlyOperator
        returns(address)
    {
        return deployNewBridgeToken(_symbol);
    }

    /*
     * @dev: Mints new BankTokens
     *
     * @param _dplatformSender: The sender's Dplatform address in bytes.
     * @param _ethereumRecipient: The intended recipient's Ethereum address.
     * @param _dplatformTokenAddress: The currency type
     * @param _symbol: dplatform token symbol
     * @param _amount: number of dplatform tokens to be minted
     */
     function mintBridgeTokens(
        bytes memory _dplatformSender,
        address payable _intendedRecipient,
        address _bridgeTokenAddress,
        string memory _symbol,
        uint256 _amount
    )
        public
        onlyDplatformBridge
    {
        return mintNewBridgeTokens(
            _dplatformSender,
            _intendedRecipient,
            _bridgeTokenAddress,
            _symbol,
            _amount
        );
    }

    /*
     * @dev: Burns bank tokens
     *
     * @param _dplatformReceiver: The _dplatform receiver address in bytes.
     * @param _dplatformTokenAddress: The currency type
     * @param _amount: number of dplatform tokens to be burned
     */
    function burnBridgeTokens(
        bytes memory _dplatformReceiver,
        address _dplatformTokenAddress,
        uint256 _amount
    )
        public
    {
        return burnDplatformTokens(
            msg.sender,
            _dplatformReceiver,
            _dplatformTokenAddress,
             _amount
        );
    }

    /*
    * @dev: Locks received Ethereum funds.
    *
    * @param _recipient: bytes representation of destination address.
    * @param _token: token address in origin chain (0x0 if ethereum)
    * @param _amount: value of deposit
    */
    function lock(
        bytes memory _recipient,
        address _token,
        uint256 _amount
    )
        public
        availableNonce()
        payable
    {
        string memory symbol;

        // Ethereum deposit
        if (msg.value > 0) {
          require(
              _token == address(0),
              "Ethereum deposits require the 'token' address to be the null address"
            );
          require(
              msg.value == _amount,
              "The transactions value must be equal the specified amount (in wei)"
            );

          // Set the the symbol to ETH
          symbol = "ETH";
          // ERC20 deposit
        } else {
          require(
              BridgeToken(_token).transferFrom(msg.sender, address(this), _amount),
              "Contract token allowances insufficient to complete this lock request"
          );
          // Set symbol to the ERC20 token's symbol
          symbol = BridgeToken(_token).symbol();
        }

        lockFunds(
            msg.sender,
            _recipient,
            _token,
            symbol,
            _amount
        );
    }

   /*
    * @dev: Unlocks Ethereum and ERC20 tokens held on the contract.
    *
    * @param _recipient: recipient's Ethereum address
    * @param _token: token contract address
    * @param _symbol: token symbol
    * @param _amount: wei amount or ERC20 token count
\   */
     function unlock(
        address payable _recipient,
        address _token,
        string memory _symbol,
        uint256 _amount
    )
        public
        onlyDplatformBridge
        hasLockedFunds(
            _token,
            _amount
        )
        canDeliver(
            _token,
            _amount
        )
    {
        unlockFunds(
            _recipient,
            _token,
            _symbol,
            _amount
        );
    }

    /*
    * @dev: Exposes an item's current status.
    *
    * @param _id: The item in question.
    * @return: Boolean indicating the lock status.
    */
    function getDplatformDepositStatus(
        bytes32 _id
    )
        public
        view
        returns(bool)
    {
        return isLockedDplatformDeposit(_id);
    }

    /*
    * @dev: Allows access to a Dplatform deposit's information via its unique identifier.
    *
    * @param _id: The deposit to be viewed.
    * @return: Original sender's Ethereum address.
    * @return: Intended Dplatform recipient's address in bytes.
    * @return: The lock deposit's currency, denoted by a token address.
    * @return: The amount locked in the deposit.
    * @return: The deposit's unique nonce.
    */
    function viewDplatformDeposit(
        bytes32 _id
    )
        public
        view
        returns(bytes memory, address payable, address, uint256)
    {
        return getDplatformDeposit(_id);
    }

}