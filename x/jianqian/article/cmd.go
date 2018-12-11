package article

import (
	"fmt"
	"github.com/QOSGroup/qstars/config"
	"github.com/QOSGroup/qstars/wire"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

const (
	flag_authoraddress="authoraddress"
	flag_originalAuthor="originalAuthor"
	flag_articleHash="articleHash"
	flag_shareAuthor="shareAuthor"
	flag_shareOriginalAuthor="shareOriginalAuthor"
	flag_shareCommunity="shareCommunity"
	flag_shareInvestor="shareInvestor"
	flag_endInvestDate="endInvestDate"
	flag_endBuyDate="endBuyDate"


)

// SendTxCmd will create a send tx and sign it with the given key.
func NewArticleCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "NewArticle",
		Short: "add Article and send tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()
			authorAddress := viper.GetString(flag_authoraddress)
			originalAuthor := viper.GetString(flag_originalAuthor)
			articleHash := viper.GetString(flag_articleHash)
			shareAuthor := viper.GetString(flag_shareAuthor)
			shareOriginalAuthor := viper.GetString(flag_shareOriginalAuthor)
			shareCommunity := viper.GetString(flag_shareCommunity)
			shareInvestor := viper.GetString(flag_shareInvestor)
			endInvestDate := viper.GetString(flag_endInvestDate)
			endBuyDate := viper.GetString(flag_endBuyDate)


			result := NewArticle(cdc,config.GetCLIContext().Config,authorAddress,originalAuthor,articleHash,shareAuthor,shareOriginalAuthor,shareCommunity,
				shareInvestor,endInvestDate,endBuyDate)
			fmt.Println(result)
			return nil
		},
	}

	cmd.Flags().String(flag_authoraddress, "", "NewArticle author address")
	cmd.Flags().String(flag_originalAuthor, "", "NewArticle original address")
	cmd.Flags().String(flag_articleHash, "", "NewArticle article hash")
	cmd.Flags().String(flag_shareAuthor, "", "NewArticle share author ")
	cmd.Flags().String(flag_shareOriginalAuthor, "", "NewArticle  share original author")
	cmd.Flags().String(flag_shareCommunity, "", "NewArticle share community")
	cmd.Flags().String(flag_shareInvestor, "", "NewArticle share investor")
	cmd.Flags().String(flag_endInvestDate, "", "NewArticle end invest date")
	cmd.Flags().String(flag_endBuyDate, "", "NewArticle end buy date")

	return cmd
}

func QueryArticleCmd(cdc *wire.Codec) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "QueryArticle",
		Short: "query  Article and send tx",
		RunE: func(cmd *cobra.Command, args []string) error {
			defer func() {
				if r := recover(); r != nil {
					fmt.Println(r)
				}
			}()

			articleHash := viper.GetString(flag_articleHash)

			result:= GetArticle(cdc,articleHash)

			fmt.Println(result)

			return nil
		},
	}


	cmd.Flags().String(flag_articleHash, "", "query article hash")

	return cmd
}