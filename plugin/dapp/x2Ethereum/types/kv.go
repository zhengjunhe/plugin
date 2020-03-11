package types

/*
 * 用户合约存取kv数据时，key值前缀需要满足一定规范
 * 即key = keyPrefix + userKey
 * 需要字段前缀查询时，使用’-‘作为分割符号
 */

var (
	//KeyPrefixStateDB state db key必须前缀
	KeyPrefixStateDB = "mavl-x2ethereum-"
	//KeyPrefixLocalDB local db的key必须前缀
	KeyPrefixLocalDB = "LODB-x2ethereum-"
)

func CalProphecyPrefix() []byte {
	return []byte(KeyPrefixStateDB + string(ProphecyKey))
}

func CalEthBridgeClaimPrefix() []byte {
	return []byte(KeyPrefixStateDB + string(EthBridgeClaimKey))
}

func CalLockPrefix() []byte {
	return []byte(KeyPrefixStateDB + string(LockKey))
}

func CalBurnPrefix() []byte {
	return []byte(KeyPrefixStateDB + string(BurnKey))
}

func CalValidatorMapsPrefix() []byte {
	return []byte(KeyPrefixStateDB + string(ValidatorMapsKey))
}

func CalLastTotalPowerPrefix() []byte {
	return []byte(KeyPrefixStateDB + string(LastTotalPowerKey))
}

func CalConsensusNeededPrefix() []byte {
	return []byte(KeyPrefixStateDB + string(ConsensusNeededKey))
}
