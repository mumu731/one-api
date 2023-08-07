package model

import (
	"math"
)

type Midjourney struct {
	Id        int    `json:"id"`
	ImageUrl  string `json:"image_url" gorm:"default:null"`
	CreatedAt string `json:"created_at" gorm:"default null"`
	UpdateAt  string `json:"update_at" gorm:"default null"`
	Prompt    string `json:"prompt" gorm:"default:null"`
	Status    string `json:"status" gorm:"default:null"`
	MessageId string `json:"message_id" gorm:"default:null"`
	UserId    int    `json:"user_id" gorm:"default:0"`
	Remarks   string `json:"remarks" gorm:"default:null"`
}

func GetAllPicture(userId int, page int, pageSize int) ([]*Midjourney, int64, error) {
	var pictures []*Midjourney
	var count int64
	var err error

	// 获取满足条件的记录总数
	err = DB.Model(&Midjourney{}).Where("user_id = ?", userId).Count(&count).Error
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
	err = DB.Where("user_id = ?", userId).Order("id desc").Offset(startIdx).Limit(pageSize).Find(&pictures).Error
	if err != nil {
		return nil, 0, err
	}

	return pictures, count, nil
}

func (picture *Midjourney) MidjourneyInsert() error {
	var err error
	err = DB.Create(picture).Error
	return err
}

func (picture *Midjourney) MidjourneyUpdate() error {
	// This can update zero values
	return DB.Model(picture).Where("message_id = ?", picture.MessageId).Updates(picture).Error
}

func (picture *Midjourney) MidjourneyDelete() error {
	var err error
	err = DB.Delete(picture).Error
	return err
}
