package offline

import (
	"context"
	"encoding/json"
	"fmt"
	tml "github.com/BurntSushi/toml"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"time"
)

//查询deploy 私钥的nonce信息，并输出到文件中
type queryCmd struct {
}

func (q *queryCmd) queryCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "query address", //first step
		Short: " query gasPrice,nonce from the spec address",
		Run:   q.query, //对要部署的factory合约进行签名
	}
	q.addQueryFlags(cmd)
	return cmd
}

func (q *queryCmd) addQueryFlags(cmd *cobra.Command) {
	cmd.Flags().StringP("address", "a", "", "deploy address")
	cmd.MarkFlagRequired("address")
}

func (q *queryCmd) query(cmd *cobra.Command, args []string) {
	url, _ := cmd.Flags().GetString("rpc_laddr")
	addr, _ := cmd.Flags().GetString("address")

	client, err := ethclient.Dial(url)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	price, err := client.SuggestGasPrice(ctx)
	if err != nil {
		panic(err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), common.HexToAddress(addr))
	if nil != err {
		fmt.Println("err:", err)
	}
	var info SignCmd
	info.From = addr
	info.GasPrice = price.Uint64()
	info.Nonce = nonce
	info.Timestamp = time.Now().Format(time.RFC3339)
	writeToFile("accountinfo.txt", &info)
	return

}

func paraseFile(file string, result interface{}) error {
	_, err := os.Stat(file)
	if err != nil {
		fmt.Println(err.Error())
		return err
	}
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	b, err := ioutil.ReadAll(f)
	if err != nil {
		panic(err)
	}
	return json.Unmarshal(b, result)

}

func writeToFile(fileName string, content interface{}) {
	jbytes, err := json.MarshalIndent(content, "", "\t")
	if err != nil {
		panic(err)
	}

	err = ioutil.WriteFile(fileName, jbytes, 0666)
	if err != nil {
		fmt.Println("Failed to write to file:", fileName)
	}
	fmt.Println("tx is written to file: ", fileName, "writeContent:", string(jbytes))
}
func InitCfg(filepath string, cfg interface{}) {
	if _, err := tml.DecodeFile(filepath, cfg); err != nil {
		fmt.Println(err)
		os.Exit(-1)
	}
	return
}
