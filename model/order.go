package model

import (
	"errors"
	"math"
	"one-api/common"
)

// import ()

type Order struct {
	Id         int     `json:"id"`
	OrderNo    string  `json:"orderNo" gorm:"unique;index"`
	CreatTime  string  `json:"creatTime" gorm:"not null"`
	PayTime    string  `json:"payTime" gorm:"default:null"`
	UpdateTime string  `json:"updateTime" gorm:"default:null"`
	State      int     `json:"state" gorm:"default:0"`
	UserId     int     `json:"userId" gorm:"default:0"`
	Price      float64 `json:"price" gorm:"default:0"`
	Remarks    string  `json:"remarks" gorm:"default:null"`
}

func GetAllOrders(userId int, page int, pageSize int) ([]*Order, int64, error) {
	var orders []*Order
	var count int64
	var err error

	// 获取满足条件的记录总数
	err = DB.Model(&Order{}).Where("user_id = ?", userId).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}

	// 计算总页数
	totalPages := int(math.Ceil(float64(count) / float64(pageSize)))

	// 校正页码，确保不超出范围
	if page < 1 {
		page = 1
	} else if page > totalPages {
		page = totalPages
	}

	// 计算起始索引和结束索引
	startIdx := (page - 1) * pageSize

	// 查询订单数据
	err = DB.Where("user_id = ?", userId).Order("id desc").Offset(startIdx).Limit(pageSize).Find(&orders).Error
	if err != nil {
		return nil, 0, err
	}

	return orders, count, nil
}

// 根据订单号查询订单
func SearchOrders(orderNo string) (Order, error) {
	var orders []*Order
	err := DB.Where("order_no = ?", orderNo).Find(&orders).Error
	if err != nil {
		return Order{}, err
	}
	if len(orders) == 0 {
		return Order{}, errors.New("订单不存在")
	}
	return *orders[0], nil
}

func (order *Order) Insert() error {
	var err error
	err = DB.Create(order).Error
	return err
}

func Update(orderNo string, newState int, date string) {
	//var err error
	//err = DB.Model(order).Select("order_no", "state").Updates(order).Error
	err := DB.Model(&Order{}).Where("order_no = ?", orderNo).Updates(map[string]interface{}{
		"state":       newState,
		"update_time": date,
	}).Error

	if err != nil {
		common.SysError("failed to update user used quota and request count: " + err.Error())
	}
}

func (order *Order) SelectUpdate() error {
	// This can update zero values
	return DB.Model(order).Select("order_no", "state").Updates(order).Error
}

func (order *Order) Delete() error {
	var err error
	err = DB.Delete(order).Error
	return err
}
