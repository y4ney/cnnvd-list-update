package api

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
	ReceviceVulDetail interface{}                    `json:"receviceVulDetail,omitempty"` // 接收到的漏洞详情
}

// CNNVDDetail CNNVD详情
type CNNVDDetail struct {
	Id                 interface{} `json:"id,omitempty"`
	VulName            string      `json:"vulName,omitempty"`
	CnnvdCode          string      `json:"cnnvdCode,omitempty"`
	CveCode            string      `json:"cveCode,omitempty"`
	PublishTime        string      `json:"publishTime,omitempty"`
	IsOfficial         int         `json:"isOfficial,omitempty"`
	Vendor             string      `json:"vendor,omitempty"`
	HazardLevel        interface{} `json:"hazardLevel,omitempty"`
	VulType            string      `json:"vulType,omitempty"`
	VulTypeName        string      `json:"vulTypeName,omitempty"`
	VulDesc            string      `json:"vulDesc"`
	AffectedProduct    interface{} `json:"affectedProduct,omitempty"`
	AffectedVendor     string      `json:"affectedVendor,omitempty"`
	ProductDesc        interface{} `json:"productDesc,omitempty"`
	AffectedSystem     interface{} `json:"affectedSystem,omitempty"`
	ReferUrl           string      `json:"referUrl,omitempty"`
	PatchId            interface{} `json:"patchId,omitempty"`
	Patch              string      `json:"patch,omitempty"`
	Deleted            interface{} `json:"deleted,omitempty"`
	Version            interface{} `json:"version,omitempty"`
	CreateUid          interface{} `json:"createUid,omitempty"`
	CreateUname        interface{} `json:"createUname,omitempty"`
	CreateTime         interface{} `json:"createTime,omitempty"`
	UpdateUid          interface{} `json:"updateUid,omitempty"`
	UpdateUname        interface{} `json:"updateUname,omitempty"`
	UpdateTime         string      `json:"updateTime,omitempty"`
	CnnvdFiledShow     string      `json:"cnnvdFiledShow,omitempty"`
	CveVulVO           interface{} `json:"cveVulVO,omitempty"`
	CveFiledShow       interface{} `json:"cveFiledShow,omitempty"`
	IbmVulVO           interface{} `json:"ibmVulVO,omitempty"`
	IbmFiledShow       interface{} `json:"ibmFiledShow,omitempty"`
	IcsCertVulVO       interface{} `json:"icsCertVulVO,omitempty"`
	IcsCertFiledShow   interface{} `json:"icsCertFiledShow,omitempty"`
	MicrosoftVulVO     interface{} `json:"microsoftVulVO,omitempty"`
	MicrosoftFiledShow interface{} `json:"microsoftFiledShow,omitempty"`
	HuaweiVulVO        interface{} `json:"huaweiVulVO,omitempty"`
	HuaweiFiledShow    interface{} `json:"huaweiFiledShow,omitempty"`
	NvdVulVO           interface{} `json:"nvdVulVO,omitempty"`
	NvdFiledShow       interface{} `json:"nvdFiledShow,omitempty"`
	Varchar1           string      `json:"varchar1,omitempty"`
}

func (v *ReqVulDetail) Fetch() error {
	resDetail, err := Post[*ResVulDetail](v, utils.FormatURL(Domain, APIVulDetail))
	if err != nil {
		return xerrors.Errorf("fail to request %s's vulDetail:%w\n", resDetail.Data.CnnvdCode, err)
	}
	err = SaveCNNVDPerYear(utils.CNNVDListDir(), resDetail.Data.CnnvdCode, resDetail.Data)
	if err != nil {
		return xerrors.Errorf("fail to save %s:%w\n", resDetail.Data.CnnvdCode, err)
	}
	log.Printf("save %s successfully", resDetail.Data.CnnvdCode)
	return nil
}
