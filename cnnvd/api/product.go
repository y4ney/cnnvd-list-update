package api

import (
	"encoding/json"
	"github.com/0yaney0/cnnvd-list-update/utils"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
	"log"
	"os"
	"path/filepath"
)

// ReqProduct 产品选择列表请求参数
type ReqProduct struct {
	ProductKeyword string `json:"productKeyword"` // 产品关键词
}

// ResProduct 产品选择列表响应参数
type ResProduct struct {
	ResCode           // 响应码
	Data    []Product `json:"data,omitempty"` // 产品列表
}

// Product 产品列表
type Product struct {
	Label string `json:"label,omitempty" gorm:"uniqueIndex size:256"`
	Value string `json:"value,omitempty"`
}
type TableProduct struct {
	gorm.Model
	Product
}

func (t *TableProduct) TableName() string {
	return ProductTable
}
func (t *TableProduct) createTable(db *gorm.DB) error {
	if db.Migrator().HasTable(t.TableName()) {
		return nil
	}
	if err := db.Migrator().CreateTable(&TableProduct{}); err != nil {
		return err
	}
	return nil
}

func (r *ReqProduct) Name() string {
	return ProductName
}

func (r *ReqProduct) Fetch() (*[]Product, error) {
	resProduct, err := Post[*ResProduct](r, utils.FormatURL(Domain, APIProduct))
	if err != nil {
		return nil, xerrors.Errorf("【%s】fail to fetch:%w\n", r.Name(), err)
	}
	var products []Product
	for _, data := range resProduct.Data {
		products = append(products, data)
	}
	log.Printf("【%s】fetch successfully!", r.Name())
	return &products, nil
}

func (r *ReqProduct) Save(data *[]Product, dir string) error {
	path := filepath.Join(dir, ProductFile)
	err := utils.Write(path, data)
	if err != nil {
		return xerrors.Errorf("【%s】fail to save :%w\n", r.Name(), err)
	}
	log.Printf("【%s】save %s successfully", r.Name(), path)
	return nil
}

func (r *ReqProduct) StoreByFile(db *gorm.DB, dir string) error {
	var mysql TableProduct
	if err := mysql.createTable(db); err != nil {
		return xerrors.Errorf("【%s】fail to create table :%w\n", r.Name(), err)
	}
	file := filepath.Join(dir, ProductFile)
	products, err := r.read(file)
	if err != nil {
		return xerrors.Errorf("【%s】fail to read %s:%w\n", r.Name(), ProductFile, err)
	}
	for _, product := range products {
		db.Create(&product)
		log.Printf("【%s】store %s successfully", r.Name(), product.Label)
	}
	return nil
}

func (r *ReqProduct) StoreByRequest(db *gorm.DB) error {
	var mysql TableProduct
	if err := mysql.createTable(db); err != nil {
		return xerrors.Errorf("【%s】fail to create table :%w\n", r.Name(), err)
	}
	products, err := r.Fetch()
	if err != nil {
		return xerrors.Errorf("【%s】fail to fetch :%w\n", r.Name(), err)
	}
	for _, product := range *products {
		db.Create(&product)
		log.Printf("【%s】store %s successfully", r.Name(), product.Label)
	}
	return nil
}

func (r *ReqProduct) read(file string) ([]TableProduct, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, xerrors.Errorf("fail to read :%w", err)
	}
	var product []TableProduct
	err = json.Unmarshal(data, &product)
	if err != nil {
		return nil, xerrors.Errorf("fail to unmarshal:%w", err)
	}
	return product, nil
}
