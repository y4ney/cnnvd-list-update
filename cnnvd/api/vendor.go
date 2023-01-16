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

// ReqVendor 供应商选择列表请求参数
type ReqVendor struct {
	VendorKeyword string `json:"vendorKeyword"` // 供应商关键词
}

// ResVendor 供应商选择列表响应参数
type ResVendor struct {
	ResCode          // 响应码
	Data    []Vendor `json:"data,omitempty"` // 供应商选择列表
}

// Vendor 供应商选择列表
type Vendor struct {
	Label string `json:"label,omitempty" gorm:"uniqueIndex size:256"`
	Value string `json:"value,omitempty"`
}

type TableVendor struct {
	gorm.Model
	Vendor
}

func (v *TableVendor) TableName() string {
	return VendorTable
}
func (v *TableVendor) createTable(db *gorm.DB) error {
	if db.Migrator().HasTable(v.TableName()) {
		return nil
	}
	if err := db.Migrator().CreateTable(&TableVendor{}); err != nil {
		return err
	}
	return nil
}

func (v *ReqVendor) Name() string {
	return VendorName
}

func (v *ReqVendor) Fetch() (*[]TableVendor, error) {
	// 获取供应商信息
	resVendor, err := Post[*ResVendor](v, utils.FormatURL(Domain, APIVendor))
	if err != nil {
		return nil, xerrors.Errorf("【%s】fail to fetch:%w\n", v.Name(), err)
	}
	var vendors []TableVendor
	for _, data := range resVendor.Data {
		vendors = append(vendors, TableVendor{Vendor: data})
	}
	log.Printf("【%s】fetch successfully!", v.Name())
	return &vendors, nil
}

func (v *ReqVendor) Save(data *[]Vendor, dir string) error {
	path := filepath.Join(dir, VendorFile)
	err := utils.Write(path, data)
	if err != nil {
		return xerrors.Errorf("【%s】fail to save :%w\n", v.Name(), err)
	}
	log.Printf("【%s】save %s successfully", v.Name(), path)
	return nil
}

func (v *ReqVendor) StoreByFile(db *gorm.DB, dir string) error {
	var mysql TableVendor
	if err := mysql.createTable(db); err != nil {
		return xerrors.Errorf("【%s】fail to create table :%w\n", v.Name(), err)
	}
	file := filepath.Join(dir, VendorFile)
	vendors, err := v.read(file)
	if err != nil {
		return xerrors.Errorf("【%s】fail to read %s:%w\n", v.Name(), VendorFile, err)
	}
	for _, vendor := range vendors {
		db.Create(&vendor)
		log.Printf("【%s】store %s successfully", v.Name(), vendor.Label)
	}
	return nil
}

func (v *ReqVendor) StoreByRequest(db *gorm.DB) error {
	var mysql TableVendor
	if err := mysql.createTable(db); err != nil {
		return xerrors.Errorf("【%s】fail to create table :%w\n", v.Name(), err)
	}
	vendors, err := v.Fetch()
	if err != nil {
		return xerrors.Errorf("【%s】fail to fetch :%w\n", v.Name(), err)
	}
	for _, vendor := range *vendors {
		db.Create(&vendor)
		log.Printf("【%s】store %s successfully", v.Name(), vendor.Label)
	}
	return nil
}
func (v *ReqVendor) read(file string) ([]TableVendor, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, xerrors.Errorf("fail to read :%w", err)
	}
	var vendor []TableVendor
	err = json.Unmarshal(data, &vendor)
	if err != nil {
		return nil, xerrors.Errorf("fail to unmarshal:%w", err)
	}
	return vendor, nil
}
