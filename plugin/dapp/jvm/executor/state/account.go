package state

import (
	"fmt"

	"github.com/33cn/chain33/common"
	"github.com/33cn/chain33/common/db"
	"github.com/33cn/chain33/common/log/log15"
	chain33Types "github.com/33cn/chain33/types"
	"github.com/golang/protobuf/proto"
	"github.com/33cn/plugin/plugin/dapp/jvm/types"
)

var (
	// ContractDataPrefix 在StateDB中合约账户保存的键值有以下几种
	// 合约数据，前缀+合约地址，第一次生成合约时设置，后面不会发生变化
	ContractDataPrefix = "mavl-jvm-data: "

	// ContractStatePrefix 合约状态，前缀+合约地址，保存合约nonce以及其它数据，可变
	ContractStatePrefix = "mavl-jvm-state: "

	// ContractStateItemKey 合约中存储的具体状态数据，包含两个参数：合约地址、状态KEY
	ContractStateItemKey = "mavl-jvm-state:%v:%v"
	// 注意，合约账户本身也可能有余额信息，这部分在CoinsAccount处理
)

// ContractAccount 合约账户对象
type ContractAccount struct {
	mdb *MemoryStateDB

	// 合约代码地址
	Addr string

	// 合约固定数据
	Data types.JVMContractData

	// 合约状态数据
	State types.JVMContractState

	// 当前的状态数据缓存
	stateCache map[string][]byte
}

// NewContractAccount 创建一个新的合约对象
// 注意，此时合约对象有可能已经存在也有可能不存在
// 需要通过LoadContract进行判断
func NewContractAccount(addr string, db *MemoryStateDB) *ContractAccount {
	if len(addr) == 0 || db == nil {
		log15.Error("NewContractAccount error, something is missing", "contract addr", addr, "db", db)
		return nil
	}
	ca := &ContractAccount{Addr: addr, mdb: db}
	ca.State.Storage = make(map[string][]byte)
	ca.stateCache = make(map[string][]byte)
	return ca
}

// GetState 获取状态数据；
// 获取数据分为两层，一层是从当前的缓存中获取，如果获取不到，再从localdb中获取
func (c *ContractAccount) GetState(key string) []byte {
	if val, ok := c.stateCache[key]; ok {
		return val
	}
	keyStr := c.GetStateItemKey(c.Addr, key)
	// 如果缓存中取不到数据，则只能到本地数据库中查询
	val, err := c.mdb.StateDB.Get([]byte(keyStr))
	if err != nil {
		log15.Debug("GetState error!", "key", key, "error", err)
		return nil
	}
	c.stateCache[key] = val
	return val
}

// SetValue2Local 设置本地数据，用于帮助辅助查找
func (c *ContractAccount) SetValue2Local(key string, value []byte) error {
	keyStr := c.GetLocalDataKey(c.Addr, key)
	c.mdb.addChange(localStorageChange{
		baseChange: baseChange{},
		account:    c.Addr,
		key:        []byte(keyStr),
		data:       value,
	})
	return c.mdb.LocalDB.Set([]byte(keyStr), value)
}

// SetState 设置状态数据
func (c *ContractAccount) SetState(key string, value []byte) error {
	c.mdb.addChange(storageChange{
		baseChange: baseChange{},
		account:    c.Addr,
		key:        []byte(key),
		prevalue:   c.GetState(key),
	})
	c.stateCache[key] = value
	keyStr := c.GetStateItemKey(c.Addr, key)
	return c.mdb.StateDB.Set([]byte(keyStr), value)
}

// TransferState 从原有的存储在一个对象，将状态数据分散存储到多个KEY，保证合约可以支撑大量状态数据
func (c *ContractAccount) TransferState() {
	if len(c.State.Storage) > 0 {
		storage := c.State.Storage
		// 为了保证不会造成新、旧数据并存的情况，需要将旧的状态数据清空
		c.State.Storage = make(map[string][]byte)

		// 从旧的区块迁移状态数据到新的区块，模拟状态数据变更的操作
		for key, value := range storage {
			c.SetState(key, value)
		}
		// 更新本合约的状态数据（删除旧的map存储信息）
		c.mdb.UpdateState(c.Addr)
		return
	}
}

// 从外部恢复合约数据
func (c *ContractAccount) restoreData(data []byte) {
	var content types.JVMContractData
	err := proto.Unmarshal(data, &content)
	if err != nil {
		log15.Error("read contract data error", c.Addr)
		return
	}

	c.Data = content
}

// 从外部恢复合约状态
func (c *ContractAccount) resotreState(data []byte) {
	var content types.JVMContractState
	err := proto.Unmarshal(data, &content)
	if err != nil {
		log15.Error("read contract state error", c.Addr)
		return
	}
	c.State = content
	if c.State.Storage == nil {
		c.State.Storage = make(map[string][]byte)
	}
}

// LoadContract 从数据库中加载合约信息（在只有合约地址的情况下）
func (c *ContractAccount) LoadContract(db db.KV) {
	// 加载代码数据
	data, err := db.Get(c.GetDataKey())
	if err != nil {
		log15.Error("StateDBGetState LoadContract:GetDataKey failed")
		return
	}
	c.restoreData(data)

	// 加载状态数据
	data, err = db.Get(c.GetStateKey())
	if err != nil {
		log15.Error("StateDBGetState LoadContract:GetStateKey failed")
		return
	}
	c.resotreState(data)
}

// SetCodeAndAbi 设置合约二进制代码
// 会同步生成代码哈希
func (c *ContractAccount) SetCodeAndAbi(code []byte, abi []byte) {
	prevcode := c.Data.GetCode()
	prevabi := c.Data.GetAbi()
	c.mdb.addChange(codeChange{
		baseChange: baseChange{},
		account:    c.Addr,
		prevhash:   c.Data.GetCodeHash(),
		prevcode:   prevcode,
		prevabi:    prevabi,
	})
	c.Data.Code = code
	c.Data.CodeHash = common.Sha256(code)
	c.Data.Abi = abi
}

// SetCreator 设置创建者
func (c *ContractAccount) SetCreator(creator string) {
	if len(creator) == 0 {
		log15.Error("SetCreator error", "creator", creator)
		return
	}
	c.Data.Creator = creator
}

// SetExecName 设置执行名称
func (c *ContractAccount) SetExecName(execName string) {
	if len(execName) == 0 {
		log15.Error("SetExecName error", "execName", execName)
		return
	}
	c.Data.Name = execName
}

// GetCreator get creator
func (c *ContractAccount) GetCreator() string {
	return c.Data.Creator
}

// GetExecName get exec name
func (c *ContractAccount) GetExecName() string {
	return c.Data.Name
}

// GetDataKV 合约固定数据，包含合约代码，以及代码哈希
func (c *ContractAccount) GetDataKV() (kvSet []*chain33Types.KeyValue) {
	c.Data.Addr = c.Addr
	datas, err := proto.Marshal(&c.Data)
	if err != nil {
		log15.Error("marshal contract data error!", "addr", c.Addr, "error", err)
		return
	}
	kvSet = append(kvSet, &chain33Types.KeyValue{Key: c.GetDataKey(), Value: datas})
	return
}

// GetStateKV 获取合约状态数据，包含nonce、是否自杀、存储哈希、存储数据
func (c *ContractAccount) GetStateKV() (kvSet []*chain33Types.KeyValue) {
	datas, err := proto.Marshal(&c.State)
	if err != nil {
		log15.Error("marshal contract state error!", "addr", c.Addr, "error", err)
		return
	}
	kvSet = append(kvSet, &chain33Types.KeyValue{Key: c.GetStateKey(), Value: datas})
	return
}

// BuildDataLog 构建变更日志
func (c *ContractAccount) BuildDataLog() (log *chain33Types.ReceiptLog) {
	logjvmContractData := types.LogJVMContractData{
		Creator:  c.Data.Creator,
		Name:     c.Data.Name,
		Addr:     c.Data.Addr,
		CodeHash: common.ToHex(c.Data.Code),
		AbiHash:  common.ToHex(c.Data.Abi),
	}

	logdatas, err := proto.Marshal(&logjvmContractData)
	if err != nil {
		log15.Error("marshal contract data error!", "addr", c.Addr, "error", err)
		return
	}
	return &chain33Types.ReceiptLog{Ty: types.TyLogContractDataJvm, Log: logdatas}
}

// BuildStateLog 构建变更日志
func (c *ContractAccount) BuildStateLog() (log *chain33Types.ReceiptLog) {
	datas, err := proto.Marshal(&c.State)
	if err != nil {
		log15.Error("marshal contract state log error!", "addr", c.Addr, "error", err)
		return
	}

	return &chain33Types.ReceiptLog{Ty: types.TyLogContractStateJvm, Log: datas}
}

// GetDataKey get data for key
func (c *ContractAccount) GetDataKey() []byte {
	return []byte("mavl-" + c.mdb.ExecutorName + "-jvmContractInfo: " + c.Addr)
}

// GetStateKey get state for key
func (c *ContractAccount) GetStateKey() []byte {
	return []byte("mavl-" + c.mdb.ExecutorName + "-state: " + c.Addr)
}

// GetStateItemKey get state item for key
func (c *ContractAccount) GetStateItemKey(addr, key string) string {
	return fmt.Sprintf("mavl-"+c.mdb.ExecutorName+"-state:%v:%v", addr, key)
}

// GetLocalDataKey get local data for key
func (c *ContractAccount) GetLocalDataKey(addr, key string) string {
	if IsPara {
		return fmt.Sprintf(string(chain33Types.LocalPrefix)+"-"+Title+c.mdb.ExecutorName+"-data-%v:%v", addr, key)
	}

	return fmt.Sprintf(string(chain33Types.LocalPrefix)+"-"+c.mdb.ExecutorName+"-data-%v:%v", addr, key)
}

// Suicide contract suicide
func (c *ContractAccount) Suicide() bool {
	c.State.Suicided = true
	return true
}

// HasSuicided return suicided
func (c *ContractAccount) HasSuicided() bool {
	return c.State.GetSuicided()
}

// Empty judge empty or not
func (c *ContractAccount) Empty() bool {
	return c.Data.GetCodeHash() == nil || len(c.Data.GetCodeHash()) == 0
}

// SetNonce set nonce
func (c *ContractAccount) SetNonce(nonce uint64) {
	c.mdb.addChange(nonceChange{
		baseChange: baseChange{},
		account:    c.Addr,
		prev:       c.State.GetNonce(),
	})
	c.State.Nonce = nonce
}

// GetNonce get nonce
func (c *ContractAccount) GetNonce() uint64 {
	return c.State.GetNonce()
}
