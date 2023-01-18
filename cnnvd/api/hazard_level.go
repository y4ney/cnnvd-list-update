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

// ReqHazardLevel 威胁等级请求参数
type ReqHazardLevel struct {
}

// ResHazardLevel 威胁等级响应参数
type ResHazardLevel struct {
	ResCode               // 响应码
	Data    []HazardLevel `json:"data,omitempty"` // 威胁等级列表
}

// HazardLevel 产品列表
type HazardLevel struct {
	DictLabel string `json:"dictLabel,omitempty" gorm:"uniqueIndex size:256"`
	DictValue string `json:"dictValue,omitempty"`
}
type TableHazardLevel struct {
	gorm.Model
	HazardLevel
}

func (t *TableHazardLevel) TableName() string {
	return HazardLevelTable
}

func (r *ReqHazardLevel) Name() string {
	return HazardLevelName
}

func (r *ReqHazardLevel) Fetch() (*[]HazardLevel, error) {
	var (
		resHazardLevel ResHazardLevel
		hazardLevel    []HazardLevel
	)

	resBody, err := utils.Fetch("GET", utils.FormatURL(Domain, APIProduct), r, Retry)
	if err != nil {
		return nil, xerrors.Errorf("【%s】fail to fetch:%w\n", r.Name(), err)
	}
	err = json.Unmarshal(resBody, &resHazardLevel)
	if err != nil {
		return nil, xerrors.Errorf("【%s】fail to unmarshal resBody :%w\n", r.Name(), err)
	}
	for _, data := range resHazardLevel.Data {
		hazardLevel = append(hazardLevel, data)
	}
	log.Printf("【%s】fetch successfully!", r.Name())
	return &hazardLevel, nil
}

func (r *ReqHazardLevel) Save(data *[]HazardLevel, dir string) error {
	path := filepath.Join(dir, HazardLevelFile)
	err := utils.Write(path, data)
	if err != nil {
		return xerrors.Errorf("【%s】fail to save :%w\n", r.Name(), err)
	}
	log.Printf("【%s】save %s successfully", r.Name(), path)
	return nil
}

func (r *ReqHazardLevel) StoreByFile(db *gorm.DB, dir string) error {
	if err := CreateTable(db, HazardLevelTable); err != nil {
		return xerrors.Errorf("【%s】fail to create table :%w\n", r.Name(), err)
	}
	file := filepath.Join(dir, HazardLevelFile)
	hazardLevel, err := r.read(file)
	if err != nil {
		return xerrors.Errorf("【%s】fail to read %s:%w\n", r.Name(), file, err)
	}
	r.store(db, hazardLevel)
	return nil
}

func (r *ReqHazardLevel) StoreByRequest(db *gorm.DB) error {
	if err := CreateTable(db, HazardLevelTable); err != nil {
		return xerrors.Errorf("【%s】fail to create table :%w\n", r.Name(), err)
	}
	hazardLevel, err := r.Fetch()
	if err != nil {
		return xerrors.Errorf("【%s】fail to fetch :%w\n", r.Name(), err)
	}
	r.store(db, hazardLevel)
	return nil
}

func (r *ReqHazardLevel) read(file string) (*[]HazardLevel, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, xerrors.Errorf("fail to read :%w", err)
	}
	var hazardLevel []HazardLevel
	err = json.Unmarshal(data, &hazardLevel)
	if err != nil {
		return nil, xerrors.Errorf("fail to unmarshal:%w", err)
	}
	return &hazardLevel, nil
}
func (r *ReqHazardLevel) store(db *gorm.DB, data *[]HazardLevel) {
	var hazardLevels []TableHazardLevel
	for _, hazardLevel := range *data {
		hazardLevels = append(hazardLevels, TableHazardLevel{HazardLevel: hazardLevel})
	}
	db.CreateInBatches(&hazardLevels, 100)
}
