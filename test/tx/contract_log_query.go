package main

import (
	"fmt"

	tc "github.com/fractalplatform/fractal/test/common"
	"github.com/fractalplatform/fractal/rpc"
	jww "github.com/spf13/jwalterweatherman"
)

func init() {
	jww.SetLogThreshold(jww.LevelTrace)
	jww.SetStdoutThreshold(jww.LevelInfo)
}

func main() {
	for i := 0; i < 1000; i++ {
		result, err := tc.GetBlockAndResult(rpc.BlockNumber(i))
		if err != nil {
			jww.ERROR.Println("get block and result failed", err)
		}
		detailtxs := result.DetailTxs
		for i := 0; i < len(detailtxs); i++ {
			details := detailtxs[i].InternalTxs
			for j := 0; j < len(details); j++ {
				logs := details[j].InterlnalLogs
				for m := 0; m < len(logs); m++ {
					log := logs[m]
					fmt.Println(log.Action.AssetID, log.Action.From)
				}
			}
		}
		receipts := result.Receipts
		for i := 0; i < len(result.Receipts); i++ {
			receipt := receipts[i]
			actionresults := receipt.ActionResults
			for j := 0; j < len(actionresults); j++ {
				actions := actionresults[j].GasAllot
				for k := 0; k < len(actionresults); k++ {
					fmt.Println(actions[k].Account.String(), actions[k].Reason, actions[k].Gas)
				}
			}
		}

	}
}
