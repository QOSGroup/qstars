package article

import (
	"fmt"
	"github.com/QOSGroup/qbase/txs"
	"github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/client/utils"
	"github.com/QOSGroup/qstars/config"
	qstartypes "github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/x/jianqian/tx"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"
	"strconv"
	"strings"
)


//上传新作品
//
//authoraddress          作者地址(必填)
//originalAuthor          原创作者地址(为空表示原创)
//articleHash            作品唯一标识hash
//shareAuthor            作者收入比例(必填)
//shareOriginalAuthor    原创收入比例(转载作品必填)
//shareCommunity         社区收入比例(必填)
//shareInvestor          投资者收入比例(必填)
//endInvestDate          投资结束时间(必填)
//endBuyDate             广告位购买结果时间(必填)
func NewArticle ( cdc *amino.Codec, authorAddress, originalAuthor,  articleHash,  shareAuthor, shareOriginalAuthor,
	shareCommunity,  shareInvestor,  endInvestDate,  endBuyDate string)string {
	privkey:=tx.GetConfig().Dappowner
	authorAddr,err:=qstartypes.AccAddressFromBech32(authorAddress)
	if err!=nil{
		return err.Error()
	}
	var originaladdr types.Address
	if strings.TrimSpace(originalAuthor)!="" {
		originaladdr, _ = qstartypes.AccAddressFromBech32(originalAuthor)
		if err != nil {
			return err.Error()
		}
	}

	authshare,_:=strconv.Atoi(shareAuthor)
	origshare,_:=strconv.Atoi(shareOriginalAuthor)
	commushare,_:=strconv.Atoi(shareCommunity)
	invesshare,_:=strconv.Atoi(shareInvestor)

	tx:=NewArticlesTx(authorAddr,originaladdr,articleHash,authshare,origshare,commushare,invesshare,endInvestDate,endBuyDate,types.ZeroInt())

	fmt.Println("tx",tx)
	_, addrben32, priv := utility.PubAddrRetrievalFromAmino(privkey, cdc)
	from, err := qstartypes.AccAddressFromBech32(addrben32)
	fmt.Println("from=",from)
	//key := account.AddressStoreKey(from)
	//acc, err := config.GetCLIContext().QOSCliContext.GetAccount(key, cdc)
	if err != nil {
		return "{Code:\"1\",Reason:\""+err.Error()+"\"}"
	}
	var nn int64=0
	//nn = int64(acc.Nonce)
	nn++
	chainid := config.GetCLIContext().Config.QOSChainID

	fmt.Println("chainid",chainid)
	txsd:=genStdSendTx(cdc,tx,priv,chainid,nn)
	cliCtx := *config.GetCLIContext().QSCCliContext

	_,_,err1:=utils.SendTx(cliCtx,cdc,txsd)

	if err1!=nil{
		return "{Code:\"1\",Reason:\""+err1.Error()+"\"}"
	}
	return "{Code:\"0\",Reason:\"\"}"
}


//封装公链交易信息
func genStdSendTx(cdc *amino.Codec, sendTx txs.ITx, priKey ed25519.PrivKeyEd25519, chainid string, nonce int64) *txs.TxStd {
	gas := types.NewInt(int64(0))
	stx := txs.NewTxStd(sendTx, chainid, gas)
	signature, _ := stx.SignTx(priKey, nonce)
	stx.Signature = []txs.Signature{txs.Signature{
		Pubkey:    priKey.PubKey(),
		Signature: signature,
		Nonce:     nonce,
	}}
	return stx
}




// GetArticle process of get Article
func GetArticle(cdc *amino.Codec, key string) (string, error) {

	//cliCtx := *config.GetCLIContext().QSCCliContext
	//value, err := cliCtx.QueryArticle([]byte(key))
	//if err != nil {
	//	return "", err
	//}
	//var articles jianqian.Articles
	//cdc.UnmarshalJSON(value,&articles)
	//
	//fmt.Println(articles)
	//result:=string(value)
	return "", nil
}