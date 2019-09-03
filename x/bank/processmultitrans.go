package bank

import (
	"encoding/hex"
	"fmt"
	"github.com/QOSGroup/qbase/account"
	"github.com/QOSGroup/qbase/txs"
	qbasetypes "github.com/QOSGroup/qbase/types"
	"github.com/QOSGroup/qstars/client/utils"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/types"
	"github.com/QOSGroup/qstars/utility"
	"github.com/QOSGroup/qstars/wire"
	"github.com/QOSGroup/qstars/x/bank/tx"
	"github.com/tendermint/go-amino"
	"github.com/tendermint/tendermint/crypto/ed25519"

	"github.com/pkg/errors"
	"strconv"
	"strings"
	"time"
)

const ACCOUNT_NOT_EXIST = "Account is not exsit."

func MultiSendDirect(cdc *wire.Codec, fromstrs []string, tos []qbasetypes.Address, scoins []types.Coins, rcoins []types.Coins) (*SendResult, error) {
	if (len(scoins) == 0) || (len(rcoins) == 0) {
		return nil, errors.New("coins不能为空")
	}
	privs := []ed25519.PrivKeyEd25519{}
	froms := []qbasetypes.Address{}
	nns := []int64{}
	for _, fromstr := range fromstrs {
		_, addrben32, priv := utility.PubAddrRetrievalFromAmino(fromstr, cdc)
		from, err := types.AccAddressFromBech32(addrben32)
		fmt.Println("from addr:", addrben32)
		fromaddr := account.AddressStoreKey(from)
		if err != nil {
			return nil, err
		}
		privs = append(privs, priv)
		froms = append(froms, from)
		fmt.Println("from:   ", fromaddr)
		acc, err := config.GetCLIContext().QOSCliContext.GetAccount(fromaddr, cdc)

		if err != nil {
			if err.Error() == ACCOUNT_NOT_EXIST {
				nn := int64(1)
				nns = append(nns, nn)
			} else {
				return nil, err
			}
		} else {
			nn := int64(acc.Nonce)
			nn++
			nns = append(nns, nn)
		}

	}

	var sendcoins [][]qbasetypes.BaseCoin
	for _, echosendercoin := range scoins {
		var sendcoinstmp []qbasetypes.BaseCoin
		for _, coin := range echosendercoin {
			sendcoinstmp = append(sendcoinstmp, qbasetypes.BaseCoin{
				Name:   coin.Denom,
				Amount: qbasetypes.NewInt(coin.Amount.Int64()),
			})
		}
		sendcoins = append(sendcoins, sendcoinstmp)
	}

	var receivecoins [][]qbasetypes.BaseCoin
	for _, echoreceivecoin := range rcoins {
		var receivecoinstmp []qbasetypes.BaseCoin
		for _, coin := range echoreceivecoin {
			receivecoinstmp = append(receivecoinstmp, qbasetypes.BaseCoin{
				Name:   coin.Denom,
				Amount: qbasetypes.NewInt(coin.Amount.Int64()),
			})
		}
		receivecoins = append(receivecoins, receivecoinstmp)
	}

	t := tx.NewTransferMultiple(froms, tos, sendcoins, receivecoins)
	var msg *txs.TxStd
	tochainid := config.GetCLIContext().Config.QOSChainID
	//fromchainid := config.GetCLIContext().Config.QSCChainID
	msg = genStdSendMultiTx(cdc, t, privs, tochainid, tochainid, nns)
	//--------------------------------------------------------------------
	rrr, _ := cdc.MarshalJSON(msg)
	fmt.Println()
	fmt.Println((string)(rrr))

	//--------------------------------------------------------------------
	ccc, _ := cdc.MarshalBinaryBare(msg)
	fmt.Println()
	fmt.Println(" \n", ccc)

	for _, fromstr1 := range fromstrs {
		fmt.Println("Private key:         ", fromstr1)
	}
	fmt.Println("nonce:              ", nns)
	//-----------------------------Signature-----------------------------------
	i := 0
	for _, nn := range nns {
		//sigdata := append(msg.BuildSignatureBytes(nn,tochainid), Int2Byte(nn)...)
		//sigdata,_ := cdc.MarshalBinaryBare(msg.Signature)
		//-------------------------------------------------------------------------

		privateStr := hex.EncodeToString(froms[i].Bytes())
		fmt.Println("from:            ", (privateStr))

		sigdata := msg.BuildSignatureBytes(nn, tochainid)
		encodedStr := hex.EncodeToString(sigdata)
		fmt.Println("Need to signdata hex:  ", (encodedStr))
		fmt.Println("Need to signdata byte: ", sigdata)

		signed1, _ := privs[i].Sign(sigdata)
		encodedStr = hex.EncodeToString(signed1)
		fmt.Println("signature hex:     ", (encodedStr))
		fmt.Println("signature byte:    ", (signed1))
		i++
	}

	for _, to := range tos {
		encodedStr := hex.EncodeToString(to.Bytes())
		fmt.Println("to: ", (encodedStr))
	}
	cliCtx := *config.GetCLIContext().QOSCliContext
	response, commitresult, err := utils.SendTx(cliCtx, cdc, msg)

	result := &SendResult{}
	if err != nil {
		result.Hash = ""
		result.Error = err.Error()
		result.Code = "1"
		return result, nil
	}
	result.Hash = response
	height := strconv.FormatInt(commitresult.Height, 10)
	result.Heigth = height

	return result, nil
}

//add the string input chainid
func genStdSendMultiTx(cdc *amino.Codec, sendTx txs.ITx, priKeys []ed25519.PrivKeyEd25519, fromchainid string, tochainid string, nonces []int64) *txs.TxStd {
	gas := qbasetypes.NewInt(int64(config.MaxGas))
	stx := txs.NewTxStd(sendTx, tochainid, gas)
	i := 0
	for _, priKey := range priKeys {
		//fmt.Println("signed data:%d",i,stx.BuildSignatureBytes(nonces[i], tochainid));
		signature, _ := stx.SignTx(priKey, nonces[i], fromchainid, tochainid)
		stx.Signature = append(stx.Signature, txs.Signature{
			Pubkey:    priKey.PubKey(),
			Signature: signature,
			Nonce:     nonces[i],
		})
		i++
	}
	return stx
}

func MultiSendViaQStars(cdc *wire.Codec, fromstrs []string, tos []qbasetypes.Address, scoins []types.Coins, rcoins []types.Coins) (*SendResult, error) {
	if (len(scoins) == 0) || (len(rcoins) == 0) {
		return nil, errors.New("coins不能为空")
	}
	privs := []ed25519.PrivKeyEd25519{}
	froms := []qbasetypes.Address{}
	nns := []int64{}
	nnsqos := []int64{}
	for _, fromstr := range fromstrs {
		_, addrben32, priv := utility.PubAddrRetrievalFromAmino(fromstr, cdc)
		from, err := types.AccAddressFromBech32(addrben32)
		fmt.Println("from addr:", addrben32)
		fromaddr := account.AddressStoreKey(from)
		if err != nil {
			return nil, err
		}
		privs = append(privs, priv)
		froms = append(froms, from)
		fmt.Println("QOS from:   ", fromaddr)

		acc, err := config.GetCLIContext().QOSCliContext.GetAccount(fromaddr, cdc)

		if err != nil {
			if err.Error() == ACCOUNT_NOT_EXIST {
				nn := int64(1)
				nns = append(nns, nn)
			} else {
				return nil, err
			}
		} else {
			nn := int64(acc.Nonce)
			nn++
			nns = append(nns, nn)
		}
		fmt.Println("QSC from:   ", fromaddr)
		accqsc, err := config.GetCLIContext().QSCCliContext.GetAccount(fromaddr, cdc)

		if err != nil {
			if err.Error() == ACCOUNT_NOT_EXIST {
				nn := int64(1)
				nnsqos = append(nnsqos, nn)
			} else {
				return nil, err
			}
		} else {
			nn := int64(accqsc.Nonce)
			nn++
			nnsqos = append(nnsqos, nn)
		}
	}

	var sendcoins [][]qbasetypes.BaseCoin
	for _, echosendercoin := range scoins {
		var sendcoinstmp []qbasetypes.BaseCoin
		for _, coin := range echosendercoin {
			sendcoinstmp = append(sendcoinstmp, qbasetypes.BaseCoin{
				Name:   coin.Denom,
				Amount: qbasetypes.NewInt(coin.Amount.Int64()),
			})
		}
		sendcoins = append(sendcoins, sendcoinstmp)
	}

	var receivecoins [][]qbasetypes.BaseCoin
	for _, echoreceivecoin := range rcoins {
		var receivecoinstmp []qbasetypes.BaseCoin
		for _, coin := range echoreceivecoin {
			receivecoinstmp = append(receivecoinstmp, qbasetypes.BaseCoin{
				Name:   coin.Denom,
				Amount: qbasetypes.NewInt(coin.Amount.Int64()),
			})
		}
		receivecoins = append(receivecoins, receivecoinstmp)
	}

	t := tx.NewTransferMultiple(froms, tos, sendcoins, receivecoins)
	var msg *txs.TxStd
	tochainid := config.GetCLIContext().Config.QOSChainID
	fromchainid := config.GetCLIContext().Config.QSCChainID
	msg = genStdSendMultiWrapTx(cdc, t, privs, fromchainid, tochainid, nns, nnsqos)

	//--------------------------------------------------------------------
	rrr, _ := cdc.MarshalJSON(msg)
	fmt.Println()
	fmt.Println((string)(rrr))

	//--------------------------------------------------------------------
	ccc, _ := cdc.MarshalBinaryBare(msg)

	fmt.Println(" \n", ccc)
	fmt.Println()
	for _, fromstr1 := range fromstrs {
		fmt.Println("Private key:         ", fromstr1)
	}
	fmt.Println("qos nonce:              ", nns)
	fmt.Println("qsc nonce:              ", nnsqos)
	//-----------------------------Signature-----------------------------------
	i := 0
	for _, nn := range nns {
		//sigdata := append(msg.BuildSignatureBytes(nn,tochainid), Int2Byte(nn)...)
		//sigdata,_ := cdc.MarshalBinaryBare(msg.Signature)
		//-------------------------------------------------------------------------

		privateStr := hex.EncodeToString(froms[i].Bytes())
		fmt.Println("from:            ", i, (privateStr))

		sigdata := msg.BuildSignatureBytes(nn, tochainid)
		encodedStr := hex.EncodeToString(sigdata)
		fmt.Println("Need to signdata hex:  ", i, (encodedStr))
		fmt.Println("Need to signdata byte: ", i, sigdata)

		signed1, _ := privs[i].Sign(sigdata)
		encodedStr = hex.EncodeToString(signed1)
		fmt.Println("signature hex:     ", i, (encodedStr))
		fmt.Println("signature byte:    ", i, (signed1))
		i++
	}

	//msg = genMultiStdWrapTx(cdc, t, priv, config.GetCLIContext().Config.QOSChainID, config.GetCLIContext().Config.QSCChainID,nn,qscnonce)

	for _, to := range tos {
		encodedStr := hex.EncodeToString(to.Bytes())
		fmt.Println("to: ", (encodedStr))
	}
	cliCtx := *config.GetCLIContext().QSCCliContext
	response, commitresult, err := utils.SendTx(cliCtx, cdc, msg)

	result := &SendResult{}
	if err != nil {
		result.Hash = ""
		result.Error = err.Error()
		result.Code = "1"
		return result, nil
	}
	result.Hash = response
	height := strconv.FormatInt(commitresult.Height, 10)
	result.Heigth = height
	waittime, err := strconv.Atoi(config.GetCLIContext().Config.WaitingForQosResult)
	if err != nil {
		panic("WaitingForQosResult should be able to convert to integer." + err.Error())
	}

	counter := 0
	for {
		if counter >= waittime {
			fmt.Println("time out")
			result.Error = "time out"
			break
		}
		resultstr, err := fetchResult(cdc, height, commitresult.Hash.String())
		if err != nil {
			fmt.Println("get result error:" + err.Error())
			result.Error = err.Error()
		}
		if resultstr == "BankStub" {
			result.Error = ""
			result.Result = resultstr
			result.Code = "-1"
		} else if resultstr != "" && resultstr != "-1" {
			fmt.Printf("get result:[%+v]\n", resultstr)
			rs := []rune(resultstr)
			index1 := strings.Index(resultstr, " ")

			result.Error = ""
			result.Result = string(rs[index1+1:])
			result.Code = string(rs[:index1])
			break
		}
		time.Sleep(500 * time.Millisecond)
		counter++
	}

	return result, nil
}

//add the string input chainid
func genStdSendMultiWrapTx(cdc *amino.Codec, sendTx txs.ITx, priKey []ed25519.PrivKeyEd25519, fromchainid string, tochainid string, qosnonce []int64, qscnonce []int64) *txs.TxStd {
	stx := genStdSendMultiTx(cdc, sendTx, priKey, fromchainid, tochainid, qosnonce)
	tx2 := txs.NewTxStd(nil, fromchainid, stx.MaxGas)
	tx2.ITxs[0] = NewWrapperSendTx(stx)

	for i := 0; i < len(qscnonce); i++ {
		qscn := qscnonce[i]
		signature, _ := tx2.SignTx(priKey[i], qscn, fromchainid, fromchainid)
		tx2.Signature = append(tx2.Signature, txs.Signature{
			Pubkey:    priKey[i].PubKey(),
			Signature: signature,
			Nonce:     qscn,
		})
	}
	return tx2
}

//add the string input chainid
//func genMultiStdWrapTx(cdc *amino.Codec, sendTx txs.ITx, priKey []ed25519.PrivKeyEd25519, tochainid string, fromchainid string, qosnonce []int64, qscnonce []int64) *txs.TxStd {
//
//	stx := genStdSendMultiTx(cdc, sendTx, priKey, tochainid,fromchainid, qosnonce)
//	tx2 := txs.NewTxStd(nil, fromchainid, stx.MaxGas)
//	tx2.ITx = NewWrapperSendTx(stx)
//	signature, _ := tx2.SignTx(priKey, qscnonce,fromchainid,fromchainid)
//	tx2.Signature = []txs.Signature{txs.Signature{
//		Pubkey:    priKey.PubKey(),
//		Signature: signature,
//		Nonce:     qscnonce,
//	}}
//
//	return tx2
//}
