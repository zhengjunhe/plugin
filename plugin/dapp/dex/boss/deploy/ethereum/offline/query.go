package offline

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"time"
)
type queryCmd struct{

}

func (q*queryCmd )queryAddrInfoCmd()*cobra.Command{
	cmd := &cobra.Command{
		Use:   "query address", //first step
		Short: " query gasPrice,nonce by address",
		Run:   q.query,//对要部署的factory合约进行签名
	}
	q.addQueryFlags(cmd)
	return cmd
}

func  (q*queryCmd )addQueryFlags(cmd *cobra.Command){
	cmd.Flags().StringP("address","a","","account address")
	cmd.MarkFlagRequired("address")
}


func (q*queryCmd )query(cmd*cobra.Command,args []string){
	url,_:=cmd.Flags().GetString("rpc_laddr")
	addr,_:=cmd.Flags().GetString("address")
	client, err := ethclient.Dial(url)
	ctx:=context.Background()
	price,err:=client.SuggestGasPrice(ctx)
	if err!=nil{
		panic(err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), 	common.HexToAddress(addr))
	if nil != err {
		fmt.Println("err:",err)
	}
	var info SignFactoryCmd
	info.From=addr
	info.GasPrice=price.Uint64()
	info.Nonce=nonce
	info.Timestamp=time.Now().String()

	writeToFile("accountinfo.txt",&info)
	return

}

//deploay Factory contractor

type deploayFactory struct{
	FactoryAddr string
	TxHash string
	Nonce uint64
	SignedRawTx string

}
func (d* deploayFactory) deployFactoryCmd()*cobra.Command{
	cmd := &cobra.Command{
		Use:   "send tx", //first step
		Short: " send signed raw tx",
		Run:   d.send,//对要部署的factory合约进行签名
	}
	d.addSendFlags(cmd)
	return cmd
}

func (d* deploayFactory)  addSendFlags(cmd *cobra.Command){
	cmd.Flags().StringP("file","f","accountinfo.txt","multi params")
	cmd.MarkFlagRequired("file")
}

func (d* deploayFactory) send(cmd*cobra.Command,args []string) {
	filePath, _ := cmd.Flags().GetString("file")
	url,_:=cmd.Flags().GetString("rpc_laddr")
	//解析文件数据
	 err := paraseFile(filePath,d)
	if err != nil {
		return
	}

		//tx:=signedTx.(string)
		tx := new(types.Transaction)
		err = tx.UnmarshalJSON(common.FromHex(d.SignedRawTx))
		if err != nil {
			panic(err)
		}
		client, err := ethclient.Dial(url)
		err = client.SendTransaction(context.Background(), tx)
		if err != nil {
			fmt.Println("err:", err)
		}

		txhash := tx.Hash().String()
		var writedata=make(map[string]interface{})
		writedata["hash"]=txhash
		writedata["factoryAddr"]=d.FactoryAddr
		writedata["timestamp"]=time.Now().String()
		writeToFile("factory.txt",writedata)
		return

}



func paraseFile(file string,result interface{})error{
	_, err := os.Stat(file)
	if err!=nil{
		fmt.Println(err.Error())
		return  err
	}
	f,err:=os.Open(file)
	if err!=nil{
		panic(err)
	}
	b,err:=ioutil.ReadAll(f)
	if err!=nil{
		panic(err)
	}
	return  json.Unmarshal(b,result)

}


func writeToFile(fileName string, content interface{}) {
	jbytes,err:=	json.MarshalIndent(content,"","\t")
	if err!=nil{
		panic(err)
	}

	err = ioutil.WriteFile(fileName,  jbytes, 0666)
	if err != nil {
		fmt.Println("Failed to write to file:", fileName)
	}
	fmt.Println("tx is written to file: ", fileName,"writeContent:",string(jbytes))
}

