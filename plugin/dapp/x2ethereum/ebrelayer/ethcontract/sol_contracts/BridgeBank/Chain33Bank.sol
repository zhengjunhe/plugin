pragma solidity ^0.5.0;

import "../../openzeppelin-solidity/contracts/math/SafeMath.sol";
import "./BridgeToken.sol";

/**
 * @title DplatformBank
 * @dev Manages the deployment and minting of ERC20 compatible BridgeTokens
 *      which represent assets based on the Dplatform blockchain.
 **/

contract DplatformBank {

    using SafeMath for uint256;

    uint256 public bridgeTokenCount;
    mapping(address => bool) public bridgeTokenWhitelist;
    mapping(bytes32 => bool) public bridgeTokenCreated;
    mapping(bytes32 => DplatformDeposit) dplatformDeposits;
    mapping(bytes32 => DplatformBurn) dplatformBurns;
    mapping(address => DepositBurnCount) depositBurnCounts;
    mapping(bytes32 => address) public token2address;

    struct DplatformDeposit {
        bytes dplatformSender;
        address payable ethereumRecipient;
        address bridgeTokenAddress;
        uint256 amount;
        bool exist;
        uint256 nonce;
    }

    struct DepositBurnCount {
        uint256 depositCount;
        uint256 burnCount;
    }

    struct DplatformBurn {
        bytes dplatformSender;
        address payable ethereumOwner;
        address bridgeTokenAddress;
        uint256 amount;
        uint256 nonce;
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

    event LogDplatformTokenBurn(
        address _token,
        string _symbol,
        uint256 _amount,
        address _ownerFrom,
        bytes _dplatformReceiver,
        uint256 _nonce
    );

    /*
     * @dev: Modifier to make sure this symbol not created now
     */
     modifier notCreated(string memory _symbol)
     {
         require(
             !hasBridgeTokenCreated(_symbol),
             "The symbol has been created already"
         );
         _;
     }

     /*
     * @dev: Modifier to make sure this symbol not created now
     */
     modifier created(string memory _symbol)
     {
         require(
             hasBridgeTokenCreated(_symbol),
             "The symbol has not been created yet"
         );
         _;
     }

    /*
    * @dev: Constructor, sets bridgeTokenCount
    */
    constructor () public {
        bridgeTokenCount = 0;
    }

    /*
    * @dev: check whether this symbol has been created yet or not
    *
    * @param _symbol: token symbol
    * @return: true or false
    */
    function hasBridgeTokenCreated(string memory _symbol) public view returns(bool) {
        bytes32 symHash = keccak256(abi.encodePacked(_symbol));
        return bridgeTokenCreated[symHash];
    }

    /*
    * @dev: Creates a new DplatformDeposit with a unique ID
    *
    * @param _dplatformSender: The sender's Dplatform address in bytes.
    * @param _ethereumRecipient: The intended recipient's Ethereum address.
    * @param _token: The currency type
    * @param _amount: The amount in the deposit.
    * @return: The newly created DplatformDeposit's unique id.
    */
    function newDplatformDeposit(
        bytes memory _dplatformSender,
        address payable _ethereumRecipient,
        address _token,
        uint256 _amount
    )
        internal
        returns(bytes32)
    {
        DepositBurnCount memory depositBurnCount = depositBurnCounts[_token];
        depositBurnCount.depositCount = depositBurnCount.depositCount.add(1);
        depositBurnCounts[_token] = depositBurnCount;

        bytes32 depositID = keccak256(
            abi.encodePacked(
                _dplatformSender,
                _ethereumRecipient,
                _token,
                _amount,
                depositBurnCount.depositCount
            )
        );

        dplatformDeposits[depositID] = DplatformDeposit(
            _dplatformSender,
            _ethereumRecipient,
            _token,
            _amount,
            true,
            depositBurnCount.depositCount
        );

        return depositID;
    }

    /*
    * @dev: Creates a new DplatformBurn with a unique ID
        *
        * @param _dplatformSender: The sender's Dplatform address in bytes.
        * @param _ethereumOwner: The owner's Ethereum address.
        * @param _token: The token Address
        * @param _amount: The amount to be burned.
        * @param _nonce: The nonce indicates the burn count for this token
        * @return: The newly created DplatformBurn's unique id.
        */
        function newDplatformBurn(
            bytes memory _dplatformSender,
            address payable _ethereumOwner,
            address _token,
            uint256 _amount,
            uint256 nonce
        )
            internal
            returns(bytes32)
        {
            bytes32 burnID = keccak256(
                abi.encodePacked(
                    _dplatformSender,
                    _ethereumOwner,
                    _token,
                    _amount,
                    nonce
                )
            );

            dplatformBurns[burnID] = DplatformBurn(
                _dplatformSender,
                _ethereumOwner,
                _token,
                _amount,
                nonce
            );

            return burnID;
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
        notCreated(_symbol)
        returns(address)
    {
        bridgeTokenCount = bridgeTokenCount.add(1);

        // Deploy new bridge token contract
        BridgeToken newBridgeToken = (new BridgeToken)(_symbol);

        // Set address in tokens mapping
        address newBridgeTokenAddress = address(newBridgeToken);
        bridgeTokenWhitelist[newBridgeTokenAddress] = true;
        bytes32 symHash = keccak256(abi.encodePacked(_symbol));
        bridgeTokenCreated[symHash] = true;
        depositBurnCounts[newBridgeTokenAddress] = DepositBurnCount(
            uint256(0),
            uint256(0));
        token2address[symHash] = newBridgeTokenAddress;

        emit LogNewBridgeToken(
            newBridgeTokenAddress,
            _symbol
        );

        return newBridgeTokenAddress;
    }

    /*
     * @dev: Mints new dplatform tokens
     *
     * @param _dplatformSender: The sender's Dplatform address in bytes.
     * @param _ethereumRecipient: The intended recipient's Ethereum address.
     * @param _dplatformTokenAddress: The currency type
     * @param _symbol: dplatform token symbol
     * @param _amount: number of dplatform tokens to be minted
     */
     function mintNewBridgeTokens(
        bytes memory _dplatformSender,
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

        newDplatformDeposit(
            _dplatformSender,
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
     * @dev: Burn dplatform tokens
     *
     * @param _from: The address to be burned from
     * @param _dplatformReceiver: The receiver's Dplatform address in bytes.
     * @param _dplatformTokenAddress: The token address of dplatform asset issued on ethereum
     * @param _amount: number of dplatform tokens to be minted
     */
    function burnDplatformTokens(
        address payable _from,
        bytes memory _dplatformReceiver,
        address _dplatformTokenAddress,
        uint256 _amount
    )
        internal
    {
        // Must be whitelisted bridge token
        require(
            bridgeTokenWhitelist[_dplatformTokenAddress],
            "Token must be a whitelisted bridge token"
        );

        // burn bridge tokens
        BridgeToken bridgeTokenInstance = BridgeToken(_dplatformTokenAddress);
        bridgeTokenInstance.burnFrom(_from, _amount);

        DepositBurnCount memory depositBurnCount = depositBurnCounts[_dplatformTokenAddress];
        require(
            depositBurnCount.burnCount + 1 > depositBurnCount.burnCount,
            "burn nonce is not available"
        );
        depositBurnCount.burnCount = depositBurnCount.burnCount.add(1);
        depositBurnCounts[_dplatformTokenAddress] = depositBurnCount;

        newDplatformBurn(
            _dplatformReceiver,
            _from,
            _dplatformTokenAddress,
            _amount,
            depositBurnCount.burnCount
        );

        emit LogDplatformTokenBurn(
            _dplatformTokenAddress,
            bridgeTokenInstance.symbol(),
            _amount,
            _from,
            _dplatformReceiver,
            depositBurnCount.burnCount
        );
    }

    /*
    * @dev: Checks if an individual DplatformDeposit exists.
    *
    * @param _id: The unique DplatformDeposit's id.
    * @return: Boolean indicating if the DplatformDeposit exists in memory.
    */
    function isLockedDplatformDeposit(
        bytes32 _id
    )
        internal
        view
        returns(bool)
    {
        return(dplatformDeposits[_id].exist);
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
    function getDplatformDeposit(
        bytes32 _id
    )
        internal
        view
        returns(bytes memory, address payable, address, uint256)
    {
        DplatformDeposit memory deposit = dplatformDeposits[_id];

        return(
            deposit.dplatformSender,
            deposit.ethereumRecipient,
            deposit.bridgeTokenAddress,
            deposit.amount
        );
    }

    function getToken2address(string memory _symbol)
        created(_symbol)
        public view returns(address)
    {
        bytes32 symHash = keccak256(abi.encodePacked(_symbol));
        return token2address[symHash];
    }
}