package common

import (
	"encoding/json"
	"fmt"
	"github.com/QOSGroup/qstars/x/jianqian/comm"
	"github.com/QOSGroup/qstars/x/jianqian/tx"

	"github.com/QOSGroup/qstars/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"log"
)

const (
	flag_authoraddress="authoraddress"
	flag_authorotheraddress="authorOtherAddr"
	flag_articletype="articleType"
	flag_articleHash="articleHash"
	flag_shareAuthor="shareAuthor"

	flag_shareOriginalAuthor="shareOriginalAuthor"
	flag_shareCommunity="shareCommunity"
	flag_shareInvestor="shareInvestor"
	flag_endInvestDate="endInvestDate"
	flag_endBuyDate="endBuyDate"
	flag_cointype="coinType"

)

// SendTxCmd will create a send tx and sign it with the given key.
func RechargeCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "quzhuan",
		Short: "quzhuan  common",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()

			authorAddress := viper.GetString(flag_authoraddress)
			articleType := viper.GetString(flag_articletype)
			articleHash := viper.GetString(flag_articleHash)
			shareAuthor := viper.GetString(flag_shareAuthor)
			shareOriginalAuthor := viper.GetString(flag_shareOriginalAuthor)
			shareCommunity := viper.GetString(flag_shareCommunity)
			shareInvestor := viper.GetString(flag_shareInvestor)
			endInvestDate := viper.GetString(flag_endInvestDate)
			endBuyDate := viper.GetString(flag_endBuyDate)
			coinType := viper.GetString(flag_cointype)


			privkey := tx.GetConfig().Dappowner
			argss:=[]string{authorAddress,articleType,articleHash,shareAuthor,shareOriginalAuthor,shareCommunity,shareInvestor,endInvestDate,endBuyDate,coinType}



			argstr,_:=json.Marshal(argss)

			fmt.Println("args",string(argstr))


			result:=comm.CommHandler(cdc,comm.ArticleTxFlag,privkey,string(argstr))



			log.Printf(result)

			return nil
		},
	}

	cmd.Flags().String(flag_authoraddress, "", "NewArticle author address")
	cmd.Flags().String(flag_authorotheraddress, "", "NewArticle Other address")
	cmd.Flags().String(flag_articletype, "", "NewArticle article type")
	cmd.Flags().String(flag_articleHash, "", "NewArticle article hash")
	cmd.Flags().String(flag_shareAuthor, "", "NewArticle share author ")
	cmd.Flags().String(flag_shareOriginalAuthor, "", "NewArticle  share original author")
	cmd.Flags().String(flag_shareCommunity, "", "NewArticle share community")
	cmd.Flags().String(flag_shareInvestor, "", "NewArticle share investor")
	cmd.Flags().String(flag_endInvestDate, "", "NewArticle end invest date")
	cmd.Flags().String(flag_endBuyDate, "", "NewArticle end buy date")
	cmd.Flags().String(flag_cointype, "", "NewArticle coin type")

	return cmd
}


