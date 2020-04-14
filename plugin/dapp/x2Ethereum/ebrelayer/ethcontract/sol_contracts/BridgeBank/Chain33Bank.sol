pragma solidity ^0.5.0;

import "../../node_modules/openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./BridgeToken.sol";

/**
 * @title Chain33Bank
 * @dev Manages the deployment and minting of ERC20 compatible BridgeTokens
 *      which represent assets based on the Chain33 blockchain.
 **/

contract Chain33Bank {

    using SafeMath for uint256;

    uint256 public bridgeTokenCount;
    mapping(address => bool) public bridgeTokenWhitelist;
    mapping(bytes32 => Chain33Deposit) chain33Deposits;

    struct Chain33Deposit {
        bytes chain33Sender;
        address payable ethereumRecipient;
        address bridgeTokenAddress;
        uint256 amount;
        bool locked;
    }

    /*
    * @dev: Event declarations
    */
    event LogNewBridgeToken(
        address _token,
        string _symbol
    );

    event LogBridgeTokenMint(
        address _token,
        string _symbol,
        uint256 _amount,
        address _beneficiary
    );

    /*
    * @dev: Constructor, sets bridgeTokenCount
    */
    constructor () public {
        bridgeTokenCount = 0;
    }

    /*
    * @dev: Creates a new Chain33Deposit with a unique ID
    *
    * @param _chain33Sender: The sender's Chain33 address in bytes.
    * @param _ethereumRecipient: The intended recipient's Ethereum address.
    * @param _token: The currency type
    * @param _amount: The amount in the deposit.
    * @return: The newly created Chain33Deposit's unique id.
    */
    function newChain33Deposit(
        bytes memory _chain33Sender,
        address payable _ethereumRecipient,
        address _token,
        uint256 _amount
    )
        internal
        returns(bytes32)
    {
        bytes32 depositID = keccak256(
            abi.encodePacked(
                _chain33Sender,
                _ethereumRecipient,
                _token,
                _amount
            )
        );

        chain33Deposits[depositID] = Chain33Deposit(
            _chain33Sender,
            _ethereumRecipient,
            _token,
            _amount,
            true
        );

        return depositID;
    }

    /*
     * @dev: Deploys a new BridgeToken contract
     *
     * @param _symbol: The BridgeToken's symbol
     */
    function deployNewBridgeToken(
        string memory _symbol
    )
        internal
        returns(address)
    {
        bridgeTokenCount = bridgeTokenCount.add(1);

        // Deploy new bridge token contract
        BridgeToken newBridgeToken = (new BridgeToken)(_symbol);

        // Set address in tokens mapping
        address newBridgeTokenAddress = address(newBridgeToken);
        bridgeTokenWhitelist[newBridgeTokenAddress] = true;

        emit LogNewBridgeToken(
            newBridgeTokenAddress,
            _symbol
        );

        return newBridgeTokenAddress;
    }

    /*
     * @dev: Mints new chain33 tokens
     *
     * @param _chain33Sender: The sender's Chain33 address in bytes.
     * @param _ethereumRecipient: The intended recipient's Ethereum address.
     * @param _chain33TokenAddress: The currency type
     * @param _symbol: chain33 token symbol
     * @param _amount: number of chain33 tokens to be minted
\    */
     function mintNewBridgeTokens(
        bytes memory _chain33Sender,
        address payable _intendedRecipient,
        address _bridgeTokenAddress,
        string memory _symbol,
        uint256 _amount
    )
        internal
    {
        // Must be whitelisted bridge token
        require(
            bridgeTokenWhitelist[_bridgeTokenAddress],
            "Token must be a whitelisted bridge token"
        );

        // Mint bridge tokens
        require(
            BridgeToken(_bridgeTokenAddress).mint(
                _intendedRecipient,
                _amount
            ),
            "Attempted mint of bridge tokens failed"
        );

        newChain33Deposit(
            _chain33Sender,
            _intendedRecipient,
            _bridgeTokenAddress,
            _amount
        );

        emit LogBridgeTokenMint(
            _bridgeTokenAddress,
            _symbol,
            _amount,
            _intendedRecipient
        );
    }

    /*
    * @dev: Checks if an individual Chain33Deposit exists.
    *
    * @param _id: The unique Chain33Deposit's id.
    * @return: Boolean indicating if the Chain33Deposit exists in memory.
    */
    function isLockedChain33Deposit(
        bytes32 _id
    )
        internal
        view
        returns(bool)
    {
        return(chain33Deposits[_id].locked);
    }

  /*
    * @dev: Gets an item's information
    *
    * @param _Id: The item containing the desired information.
    * @return: Sender's address.
    * @return: Recipient's address in bytes.
    * @return: Token address.
    * @return: Amount of ethereum/erc20 in the item.
    * @return: Unique nonce of the item.
    */
    function getChain33Deposit(
        bytes32 _id
    )
        internal
        view
        returns(bytes memory, address payable, address, uint256)
    {
        Chain33Deposit memory deposit = chain33Deposits[_id];

        return(
            deposit.chain33Sender,
            deposit.ethereumRecipient,
            deposit.bridgeTokenAddress,
            deposit.amount
        );
    }
}