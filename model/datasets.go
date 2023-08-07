package model

import (
	"math"
)

type MDatasets struct {
	Id          int    `json:"id"`
	Name        string `json:"name" gorm:"not null"`
	CreatBy     string `json:"created_by" gorm:"not null"`
	CreatTime   string `json:"creat_time" gorm:"default:null"`
	UpdateTime  string `json:"updateTime" gorm:"default:null"`
	Description string `json:"description" gorm:"default:null"`
	WordCount   int    `json:"word_count" gorm:"default:0"`
	UserId      int    `json:"userId" gorm:"default:0"`
	Remarks     string `json:"remarks" gorm:"default:null"`
}

func GetAllDatasets(id int, page int, pageSize int) ([]*MDatasets, int64, error) {
	var datasets []*MDatasets
	var count int64
	var err error

	// 获取满足条件的记录总数
	err = DB.Model(&MDatasets{}).Where("user_id = ?", id).Count(&count).Error
	if err != nil {
		return nil, 0, err
	}
	println(count)

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

	// 查询数据集数据
	err = DB.Where("user_id = ?", id).Offset(startIdx).Limit(pageSize).Find(&datasets).Error
	if err != nil {
		return nil, 0, err
	}

	return datasets, count, nil
}

func (datasets *MDatasets) DatasetsInsert() error {
	var err error
	err = DB.Create(datasets).Error
	return err
}

func (datasets *MDatasets) Delete() error {
	var err error
	err = DB.Delete(datasets).Error
	return err
}

func DeleteDatasetsById(id string) (err error) {
	datasets := MDatasets{CreatBy: id}
	err = DB.Where(datasets).First(&datasets).Error
	if err != nil {
		return err
	}
	return datasets.Delete()
}
