package api

import (
	"github.com/0yaney0/cnnvd-list-update/utils"
	"golang.org/x/xerrors"
	"log"
	"math"
)

// ReqVulList cnnvd漏洞列表请求参数
type ReqVulList struct {
	PageIndex   int    `json:"pageIndex"`   // 分页索引
	PageSize    int    `json:"pageSize"`    // 分页大小
	Keyword     string `json:"keyword"`     // 关键字
	HazardLevel string `json:"hazardLevel"` // 漏洞等级
	VulType     string `json:"vulType"`     // 漏洞类型
	Vendor      string `json:"vendor"`      // 供应商
	Product     string `json:"product"`     // 产品
	DateType    string `json:"dateType"`    // 数据类型
}

// ResVulList 供应商选择列表响应参数
type ResVulList struct {
	ResCode         // 响应码
	Data    VulList `json:"data,omitempty"` // cnnvd漏洞列表
}

// VulList cnnvd漏洞列表
type VulList struct {
	Total     int      `json:"total,omitempty"`
	Records   []Record `json:"records,omitempty"`
	PageIndex int      `json:"pageIndex,omitempty"`
	PageSize  int      `json:"pageSize,omitempty"`
}

// Record 漏洞列表记录
type Record struct {
	Id          string      `json:"id,omitempty"`
	VulName     string      `json:"vulName,omitempty"`
	CnnvdCode   string      `json:"cnnvdCode,omitempty"`
	CveCode     string      `json:"cveCode,omitempty"`
	HazardLevel interface{} `json:"hazardLevel,omitempty"`
	CreateTime  string      `json:"createTime,omitempty"`
	PublishTime string      `json:"publishTime,omitempty"`
	UpdateTime  string      `json:"updateTime,omitempty"`
	TypeName    string      `json:"typeName,omitempty"`
	VulType     string      `json:"vulType,omitempty"`
}

func NewReqVulList(keyword string) *ReqVulList {
	return &ReqVulList{
		Keyword:  keyword,
		PageSize: PageSize,
	}
}

func (r *ReqVulList) Fetch() error {
	num, err := r.getPageNum()
	if err != nil {
		return xerrors.Errorf("fail to get page num:%w\n", err)
	}
	for i := 1; i <= num; i++ {
		r.PageIndex = i
		resList, err := Post[*ResVulList](r, utils.FormatURL(Domain, APIVulList))
		if err != nil {
			return xerrors.Errorf("fail to request vulList:%w\n", err)
		}
		log.Printf("第%v/%v页", i, num)
		for _, record := range resList.Data.Records {
			// 并保存到本地
			vulDetail := ReqVulDetail{
				Id:        record.Id,
				VulType:   record.VulType,
				CnnvdCode: record.CnnvdCode,
			}
			if err := vulDetail.Fetch(); err != nil {
				return err
			}
		}
	}

	return nil
}

func (r *ReqVulList) getPageNum() (num int, err error) {
	result, err := Post[*ResVulList](r, utils.FormatURL(Domain, APIVulList))
	if err != nil {
		return 0, xerrors.Errorf("fail to request vulList:%w\n", err)
	}
	return int(math.Ceil(float64(result.Data.Total) / float64(r.PageSize))), nil
}

// GetLatestCNNVD 获取最新的漏洞编号
func (r *ReqVulList) GetLatestCNNVD() (string, error) {
	//请求漏洞列表
	resList, err := Post[*ResVulList](r, utils.FormatURL(Domain, APIVulList))
	if err != nil {
		return "", xerrors.Errorf("fail to request vulList:%w\n", err)
	}

	// 设置第一个漏洞为最新的漏洞
	var latestCNNVD string
	latestCNNVD = resList.Data.Records[0].CnnvdCode
	for j := 1; j <= r.PageSize-1; j++ {
		latestCNNVD, err = LatestCNNVD(resList.Data.Records[j].CnnvdCode, latestCNNVD)
		if err != nil {
			return "", xerrors.Errorf("fail to get latest cnnvd:%w\n", err)
		}
	}
	return latestCNNVD, nil
}
