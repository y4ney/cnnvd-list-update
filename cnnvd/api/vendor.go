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

func (r *ReqVendor) Name() string {
	return VendorName
}

func (r *ReqVendor) Fetch() (*[]Vendor, error) {
	// 获取供应商信息
	resVendor, err := Post[*ResVendor](r, utils.FormatURL(Domain, APIVendor))
	if err != nil {
		return nil, xerrors.Errorf("【%s】fail to fetch:%w\n", r.Name(), err)
	}
	var vendors []Vendor
	for _, data := range resVendor.Data {
		vendors = append(vendors, data)
	}
	log.Printf("【%s】fetch successfully!", r.Name())
	return &vendors, nil
}

func (r *ReqVendor) Save(data *[]Vendor, dir string) error {
	path := filepath.Join(dir, VendorFile)
	err := utils.Write(path, data)
	if err != nil {
		return xerrors.Errorf("【%s】fail to save :%w\n", r.Name(), err)
	}
	log.Printf("【%s】save %s successfully", r.Name(), path)
	return nil
}

func (r *ReqVendor) StoreByFile(db *gorm.DB, dir string) error {
	if err := CreateTable(db, VendorTable); err != nil {
		return xerrors.Errorf("【%s】fail to create table :%w\n", r.Name(), err)
	}
	file := filepath.Join(dir, VendorFile)
	vendors, err := r.read(file)
	if err != nil {
		return xerrors.Errorf("【%s】fail to read %s:%w\n", r.Name(), VendorFile, err)
	}
	r.store(db, vendors)
	return nil
}

func (r *ReqVendor) StoreByRequest(db *gorm.DB) error {
	if err := CreateTable(db, VendorTable); err != nil {
		return xerrors.Errorf("【%s】fail to create table :%w\n", r.Name(), err)
	}
	vendors, err := r.Fetch()
	if err != nil {
		return xerrors.Errorf("【%s】fail to fetch :%w\n", r.Name(), err)
	}
	r.store(db, vendors)
	return nil
}

func (r *ReqVendor) read(file string) (*[]Vendor, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, xerrors.Errorf("fail to read :%w", err)
	}
	var vendors []Vendor
	err = json.Unmarshal(data, &vendors)
	if err != nil {
		return nil, xerrors.Errorf("fail to unmarshal:%w", err)
	}
	return &vendors, nil
}

func (r *ReqVendor) store(db *gorm.DB, data *[]Vendor) {
	var vendors []TableVendor
	for _, vendor := range *data {
		vendors = append(vendors, TableVendor{Vendor: vendor})
	}
	db.CreateInBatches(&vendors, 100)
}
