package api

import (
	"encoding/json"
	"gorm.io/gorm"
	"os"
	"path/filepath"
)

import (
	"github.com/0yaney0/cnnvd-list-update/utils"
	"golang.org/x/xerrors"
	"log"
)

// ReqVulDetail cnnvd漏洞详情请求参数
type ReqVulDetail struct {
	Id        string `json:"id"`        // 漏洞id
	VulType   string `json:"vulType"`   // 漏洞类型
	CnnvdCode string `json:"cnnvdCode"` // cnnvd编号
}

// ResVulDetail cnnvd漏洞详情响应参数
type ResVulDetail struct {
	ResCode           // 响应码
	Data    VulDetail `json:"data,omitempty"` // 漏洞详情数据
}

// VulDetail 漏洞详情数据
type VulDetail struct {
	CNNVDDetail       `json:"cnnvdDetail,omitempty"` // CNNVD详情
	ReceviceVulDetail string                         `json:"receviceVulDetail,omitempty"` // 接收到的漏洞详情
}

// CNNVDDetail CNNVD详情
type CNNVDDetail struct {
	Id                 string `json:"id,omitempty"`
	VulName            string `json:"vulName,omitempty"`
	CnnvdCode          string `json:"cnnvdCode,omitempty"`
	CveCode            string `json:"cveCode,omitempty"`
	PublishTime        string `json:"publishTime,omitempty"`
	IsOfficial         int    `json:"isOfficial,omitempty"`
	Vendor             string `json:"vendor,omitempty"`
	HazardLevel        int    `json:"hazardLevel,omitempty"`
	VulType            string `json:"vulType,omitempty"`
	VulTypeName        string `json:"vulTypeName,omitempty"`
	VulDesc            string `json:"vulDesc"`
	AffectedProduct    string `json:"affectedProduct,omitempty"`
	AffectedVendor     string `json:"affectedVendor,omitempty"`
	ProductDesc        string `json:"productDesc,omitempty"`
	AffectedSystem     string `json:"affectedSystem,omitempty"`
	ReferUrl           string `json:"referUrl,omitempty"`
	PatchId            string `json:"patchId,omitempty"`
	Patch              string `json:"patch,omitempty"`
	Deleted            string `json:"deleted,omitempty"`
	Version            string `json:"version,omitempty"`
	CreateUid          string `json:"createUid,omitempty"`
	CreateUname        string `json:"createUname,omitempty"`
	CreateTime         string `json:"createTime,omitempty"`
	UpdateUid          string `json:"updateUid,omitempty"`
	UpdateUname        string `json:"updateUname,omitempty"`
	UpdateTime         string `json:"updateTime,omitempty"`
	CnnvdFiledShow     string `json:"cnnvdFiledShow,omitempty"`
	CveVulVO           string `json:"cveVulVO,omitempty"`
	CveFiledShow       string `json:"cveFiledShow,omitempty"`
	IbmVulVO           string `json:"ibmVulVO,omitempty"`
	IbmFiledShow       string `json:"ibmFiledShow,omitempty"`
	IcsCertVulVO       string `json:"icsCertVulVO,omitempty"`
	IcsCertFiledShow   string `json:"icsCertFiledShow,omitempty"`
	MicrosoftVulVO     string `json:"microsoftVulVO,omitempty"`
	MicrosoftFiledShow string `json:"microsoftFiledShow,omitempty"`
	HuaweiVulVO        string `json:"huaweiVulVO,omitempty"`
	HuaweiFiledShow    string `json:"huaweiFiledShow,omitempty"`
	NvdVulVO           string `json:"nvdVulVO,omitempty"`
	NvdFiledShow       string `json:"nvdFiledShow,omitempty"`
	Varchar1           string `json:"varchar1,omitempty"`
}

type TableVulDetail struct {
	gorm.Model
	VulDetail
}

// TableName 为Record绑定表名
func (t *TableVulDetail) TableName() string {
	return VulDetailTable
}

func (r *ReqVulDetail) Name() string {
	return VulDetailName
}

func (r *ReqVulDetail) Fetch() (*VulDetail, error) {
	if r.Id == "" || r.VulType == "" || r.CnnvdCode == "" {
		return nil, xerrors.New("please specify id,vul type and cnnvd vode")
	}
	resDetail, err := Post[*ResVulDetail](r, utils.FormatURL(Domain, APIVulDetail))
	if err != nil {
		return nil, xerrors.Errorf("【%s】fail to request %s's vulDetail:%w\n", r.Name(), resDetail.Data.CnnvdCode, err)
	}
	log.Printf("【%s】fetch %s successfully", r.Name(), resDetail.Data.CnnvdCode)
	return &resDetail.Data, nil
}

func (r *ReqVulDetail) Save(data *VulDetail, dir string) error {
	vulDetail := filepath.Join(dir, VulDetailFile)
	exist, err := utils.PathExists(vulDetail)
	if err != nil {
		return xerrors.Errorf("【%s】fail to determine whether %s is dir:%w", r.Name(), vulDetail, err)
	}
	if !exist {
		err = os.MkdirAll(vulDetail, os.ModePerm)
		if err != nil {
			return xerrors.Errorf("【%s】fail to mkdir %s:%w", r.Name(), vulDetail, err)
		}
	}
	err = SaveCNNVDPerYear(vulDetail, data.CnnvdCode, data)
	if err != nil {
		return xerrors.Errorf("【%s】fail to save %s:%w\n", r.Name(), data.CnnvdCode, err)
	}
	log.Printf("【%s】save %s successfully", r.Name(), data.CnnvdCode)
	return nil
}

func (r *ReqVulDetail) StoreByFile(db *gorm.DB, dir string) error {
	if err := CreateTable(db, VulDetailFile); err != nil {
		return xerrors.Errorf("【%s】fail to create table :%w\n", r.Name(), err)
	}
	// 遍历文件夹下的文件
	files, err := utils.GetFile(filepath.Join(dir, VulDetailFile))
	if err != nil {
		return xerrors.Errorf("【%s】fail to get all files from %s:%w", r.Name(), filepath.Join(dir, VulDetailFile), err)
	}
	var vuls []VulDetail
	for _, file := range files {
		vul, err := r.read(file)
		if err != nil {
			return xerrors.Errorf("【%s】fail to read %s:%w\n", r.Name(), file, err)
		}
		vuls = append(vuls, *vul)
	}
	r.store(db, &vuls)
	return nil
}

func (r *ReqVulDetail) StoreByRequest(db *gorm.DB) error {
	if err := CreateTable(db, VulDetailTable); err != nil {
		return xerrors.Errorf("【%s】fail to create table :%w\n", r.Name(), err)
	}
	vulDetail, err := r.Fetch()
	if err != nil {
		return xerrors.Errorf("【%s】fail to fetch :%w\n", r.Name(), err)
	}
	var vuls []VulDetail
	vuls = append(vuls, *vulDetail)
	r.store(db, &vuls)
	return nil
}

func (r *ReqVulDetail) read(file string) (*VulDetail, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, xerrors.Errorf("fail to read :%w", err)
	}
	var vulDetail VulDetail
	err = json.Unmarshal(data, &vulDetail)
	if err != nil {
		return nil, xerrors.Errorf("fail to unmarshal:%w", err)
	}
	return &vulDetail, nil
}

func (r *ReqVulDetail) store(db *gorm.DB, data *[]VulDetail) {
	var vuls []TableVulDetail
	for _, vul := range *data {
		vuls = append(vuls, TableVulDetail{VulDetail: vul})
	}
	db.Create(&vuls)
}
