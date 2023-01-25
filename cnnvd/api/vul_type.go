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

// ReqVulType 漏洞类型列表请求参数
type ReqVulType struct {
}

// ResVulType 漏洞类型列表响应参数
type ResVulType struct {
	ResCode           // 响应码
	Data    []VulType `json:"data,omitempty"` // 漏洞类型列表
}

// VulType 漏洞类型列表
type VulType struct {
	Id      string    `json:"id,omitempty"`
	Pid     string    `json:"pid,omitempty"`
	Label   string    `json:"label,omitempty"`
	Value   string    `json:"value,omitempty"`
	VulType []VulType `json:"children,omitempty"`
}

type TableVulType struct {
	gorm.Model
	Id    string `json:"id,omitempty"`
	Pid   string `json:"pid,omitempty"`
	Label string `json:"label,omitempty"`
}

func (t *TableVulType) TableName() string {
	return VulTypeTable
}

func (r *ReqVulType) Name() string {
	return VulTypeName
}

func (r *ReqVulType) Fetch() (*[]VulType, error) {
	var (
		resVulType ResVulType
		vulTypes   []VulType
	)

	// 获取漏洞类型
	resBody, err := utils.Fetch("POST", utils.FormatURL(Domain, APIVulType), r, Retry)
	if err != nil {
		return nil, xerrors.Errorf("【%s】fail to fetch:%w\n", r.Name(), err)
	}

	err = json.Unmarshal(resBody, &resVulType)
	if err != nil {
		return nil, xerrors.Errorf("【%s】fail to unmarshal resBody :%w\n", r.Name(), err)
	}

	for _, vulType := range resVulType.Data {
		vulTypes = append(vulTypes, vulType)
	}
	log.Printf("【%s】fetch successfully!", r.Name())
	return &vulTypes, nil

}
func (r *ReqVulType) Save(data *[]VulType, dir string) error {
	path := filepath.Join(dir, VulTypeFile)
	err := utils.Write(path, data)
	if err != nil {
		return xerrors.Errorf("【%s】fail to save :%w\n", r.Name(), err)
	}
	log.Printf("【%s】save %s successfully", r.Name(), path)
	return nil
}
func (r *ReqVulType) StoreByFile(db *gorm.DB, dir string) error {
	if err := CreateTable(db, VulTypeTable); err != nil {
		return xerrors.Errorf("【%s】fail to create table :%w\n", r.Name(), err)
	}
	file := filepath.Join(dir, VulTypeFile)
	vulTypes, err := r.read(file)
	if err != nil {
		return xerrors.Errorf("【%s】fail to read %s:%w\n", r.Name(), ProductFile, err)
	}
	r.store(db, vulTypes)
	return nil
}

func (r *ReqVulType) StoreByRequest(db *gorm.DB) error {
	if err := CreateTable(db, VulTypeTable); err != nil {
		return xerrors.Errorf("【%s】fail to create table :%w\n", r.Name(), err)
	}
	vulTypes, err := r.Fetch()
	if err != nil {
		return xerrors.Errorf("【%s】fail to fetch :%w\n", r.Name(), err)
	}
	r.store(db, vulTypes)
	return nil
}

func (r *ReqVulType) read(file string) (*[]VulType, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, xerrors.Errorf("fail to read :%w", err)
	}
	var vulTypes []VulType
	err = json.Unmarshal(data, &vulTypes)
	if err != nil {
		return nil, xerrors.Errorf("fail to unmarshal:%w", err)
	}
	return &vulTypes, nil
}

func (r *ReqVulType) store(db *gorm.DB, data *[]VulType) {
	var vulTypes []TableVulType
	traversing(data, &vulTypes)
	db.CreateInBatches(&vulTypes, 100)
}

// traversing 递归遍历data，并转化为 vuls
func traversing(data *[]VulType, vulTypes *[]TableVulType) {
	for _, d := range *data {
		var children *[]VulType
		*vulTypes = append(*vulTypes, TableVulType{
			Id:    d.Id,
			Pid:   d.Pid,
			Label: d.Label,
		})
		if d.VulType != nil {
			children = &d.VulType
			traversing(children, vulTypes)
		}
	}
}
