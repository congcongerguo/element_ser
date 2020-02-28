package transaction

import (
	"elementser/src/common"

	"github.com/gin-gonic/gin"
)

type TransactionList struct{}

type Item struct {
	Order_no  string `json:"order_no"`
	Timestamp string `json:"timestamp"`
	UserNmae  string `json:"username"`
	Price     string `json:"price"`
	Status    string `json:"status"`
}

func (TransactionList) List(c *gin.Context) {
	var total uint64

	list := [1]Item{}
	total = 1

	list[0].Order_no = "aaa"
	list[0].Timestamp = "bbb"
	list[0].UserNmae = "ccc"
	list[0].Price = "ddd"
	list[0].Status = "success"

	common.ResSuccessPage(c, total, &list)
}
