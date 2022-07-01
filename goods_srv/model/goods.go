package model

import (
	"database/sql/driver"
	"encoding/json"
	"gorm.io/gorm"
)

type gormList struct {

}

func (g gormList) Value() (driver.Value,error) {
	return json.Marshal(g)
}

func (g *gormList) Scan(value interface{}) error {
	return json.Unmarshal(value.([]byte),&g)
}


type Category struct {
	gorm.Model
	ParentCategoryID int32
	ParementCategory *Category `gorm:"-"`
	Name string `gorm:"type:varchar(20);not null"`
	Level int32 `gorm:"type:int;not null;default:1"`
	IsTab bool `gorm:"default:false;nut null"`
}

type Brands struct {
	gorm.Model
	Name string `gorm:"type:varchar(20);not null"`
	Logo string `gorm:"type:varchar(200);default:'';not null;"`
}

type GoodsCategoryBrand struct {
	gorm.Model
	CategoryID int32 `gorm:"type:int;uniqueIndex"`
	Category Category `gorm:"-"`

	BrandsID int32 `gorm:"type:int;uniqueIndex"`
	Brands Brands `gorm:"-"`
}

func (GoodsCategoryBrand) TableName() string {
	return "goodcategorybrand"
}

type Banner struct {
	gorm.Model
	Image string `gorm:"type:varchar(200);not null"`
	Url string `gorm:"type:varchar(200);not null"`
	Index int32 `gorm:"type:int;not null;default:1"`
}

type Goods struct {
	gorm.Model
	CategoryId int32 `gorm:"type:int;not null"`
	Category Category `gorm:"-"`
	BrandSID int32 `gorm:"type:int;not null"`
	Brands Brands `gorm:"-"`
	OnSale bool `gorm:"default:false;not null"`
	ShipFree bool `gorm:"default:false;not null"`
	IsNew bool `gorm:"default:false;not null"`
	IsHot bool `gorm:"default:false;not null"`
	Name string `gorm:"type:varchar(50);not null"`
	GoodsSn string `gorm:"type:varchar(50);not null"`
	ClickNum int32 `gorm:"type:int;default:0;not null"`
	SoldNum int32 `gorm:"type:int;default:0;not null"`
	FavNum int32 `gorm:"type:int;default:0;not null"`
	MarketPrice float32 `gorm:"not null"`
	ShopPrice float32 `gorm:"not null"`
	GoodBrif string `gorm:"type:varchar(100);not null"`
	Images gormList `gorm:"type:varchar(1000);not null"`
	DescImages gormList `gorm:"type:varchar(1000);not null"`
	GoodsFrontImages string `gorm:"type:varchar(200);not null"`
}

type GoodsImages struct {
	GoodsID int
	Images string
}