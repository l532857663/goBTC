package ord

import (
	"fmt"
	"goBTC/elastic"
	"goBTC/models"
)

func DealWithBrc20Info(tokenInfo *elastic.OrdToken, brc20 *models.OrdBRC20) error {
	switch brc20.OP {
	case "deploy":
		fmt.Printf("wch----- deploy: %+v\n", brc20)
		ok, err := CheckBrc20ForDeploy(brc20)
		if err != nil {
			return err
		}
		if !ok {
		}
	case "mint":
		fmt.Printf("wch----- mint: %+v\n", brc20)
		// 判断币种mint完没有
	case "transfer":
		fmt.Printf("wch----- transfer: %+v\n", brc20)
	default:
		fmt.Printf("wch----- default: %+v\n", brc20)
	}
	return nil
}

func CheckBrc20ForDeploy(brc20 *models.OrdBRC20) (bool, error) {
	// ID信息
	dId := GetDeployIdStr(brc20.Tick)
	res, err := elastic.GetDataById(elastic.DeployType, dId)
	if err != nil {
		return false, err
	}
	// 判断之前块高的数据有没有
	fmt.Printf("wch------ res: %+v\n", res)
	// // 之前没数据，添加新数据
	// err := elastic.CreateData(elastic.DeployType, inscribeId, ord)
	// if err != nil {
	// 	return err
	// }
	return true, nil
}
