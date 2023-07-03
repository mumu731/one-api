package model

// import ()

type Order struct {
	Id        int    `json:"id"`
	OrderNo   string `json:"orderNo" gorm:"unique;index"`
	CreatTime string `json:"creatTime" gorm:"not null;"`
}

func GetAllOrders(startIdx int, num int) (orders []*Order, err error) {
	err = DB.Order("id desc").Limit(num).Offset(startIdx).Find(&orders).Error
	return orders, err
}
