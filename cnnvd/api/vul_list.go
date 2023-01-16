package api

import (
	"encoding/json"
	"github.com/0yaney0/cnnvd-list-update/utils"
	"golang.org/x/xerrors"
	"gorm.io/gorm"
	"log"
	"math"
	"os"
	"path/filepath"
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
	Id          string `json:"id,omitempty"`
	VulName     string `json:"vulName,omitempty"`
	CnnvdCode   string `json:"cnnvdCode,omitempty"`
	CveCode     string `json:"cveCode,omitempty"`
	HazardLevel int    `json:"hazardintel,omitempty"`
	CreateTime  string `json:"createTime,omitempty"`
	PublishTime string `json:"publishTime,omitempty"`
	UpdateTime  string `json:"updateTime,omitempty"`
	TypeName    string `json:"typeName,omitempty"`
	VulType     string `json:"vulType,omitempty"`
}
type TableVulList struct {
	gorm.Model
	Record
}

func (l *TableVulList) TableName() string {
	return VulListTable
}
func (l *TableVulList) createTable(db *gorm.DB) error {
	if db.Migrator().HasTable(VulListTable) {
		return nil
	}
	if err := db.Migrator().CreateTable(&TableVulList{}); err != nil {
		return xerrors.Errorf("fail to create %s table:%w\n", l.TableName(), err)
	}
	return nil
}

func (r *ReqVulList) Name() string {
	return VulListName
}

func NewReqVulList(keyword string) *ReqVulList {
	return &ReqVulList{
		Keyword:  keyword,
		PageSize: PageSize,
	}
}

func (r *ReqVulList) Fetch() (*[]TableVulList, error) {
	num, err := r.getPageNum()
	if err != nil {
		return nil, xerrors.Errorf("【%s】fail to get page num:%w\n", r.Name(), err)
	}
	var vuls []TableVulList
	for i := 1; i <= num; i++ {
		r.PageIndex = i
		resList, err := Post[*ResVulList](r, utils.FormatURL(Domain, APIVulList))
		if err != nil {
			return nil, xerrors.Errorf("【%s】fail to fetch:%w\n", r.Name(), err)
		}
		log.Printf("【%s】第%v/%v页", r.Name(), i, num)
		for _, record := range resList.Data.Records {
			vuls = append(vuls, TableVulList{Record: record})
			log.Printf("【%s】fetch %s successfully!", r.Name(), record.CnnvdCode)
		}
	}
	return &vuls, nil
}

func (r *ReqVulList) Save(data *[]TableVulList, dir string) error {
	vulList := filepath.Join(dir, VulListFile)
	exist, err := utils.PathExists(vulList)
	if err != nil {
		return xerrors.Errorf("【%s】fail to determine whether %s is dir:%w", r.Name(), vulList, err)
	}
	if !exist {
		err = os.MkdirAll(vulList, os.ModePerm)
		if err != nil {
			return xerrors.Errorf("【%s】fail to mkdir %s:%w", r.Name(), vulList, err)
		}
	}

	for _, vul := range *data {
		err = SaveCNNVDPerYear(vulList, vul.CnnvdCode, vul)
		if err != nil {
			return xerrors.Errorf("【%s】fail to save %s:%w\n", r.Name(), vul.CnnvdCode, err)
		}
		log.Printf("【%s】save %s successfully", r.Name(), vul.CnnvdCode)
	}
	return nil
}

func (r *ReqVulList) StoreByFile(db *gorm.DB, dir string) error {
	var mysql TableVulList
	if err := mysql.createTable(db); err != nil {
		return xerrors.Errorf("【%s】fail to create table :%w\n", r.Name(), err)
	}
	// 遍历文件夹下的文件
	files, err := utils.GetFile(filepath.Join(dir, VulListFile))
	if err != nil {
		return xerrors.Errorf("【%s】fail to get all files from %s:%w", r.Name(), filepath.Join(dir, VulListFile), err)
	}
	for _, file := range files {
		vul, err := r.read(file)
		if err != nil {
			return xerrors.Errorf("【%s】fail to read %s:%w\n", r.Name(), file, err)
		}
		db.Create(vul)
		log.Printf("【%s】store %s successfully", r.Name(), vul.CnnvdCode)
	}
	return nil
}

func (r *ReqVulList) StoreByRequest(db *gorm.DB) error {
	var mysql TableVulList
	if err := mysql.createTable(db); err != nil {
		return xerrors.Errorf("【%s】fail to create table :%w\n", r.Name(), err)
	}
	vulList, err := r.Fetch()
	if err != nil {
		return xerrors.Errorf("【%s】fail to fetch :%w\n", r.Name(), err)
	}
	for _, vul := range *vulList {
		mysql.Id = vul.Id
		mysql.VulName = vul.VulName
		mysql.CnnvdCode = vul.CnnvdCode
		mysql.CveCode = vul.CveCode
		mysql.HazardLevel = vul.HazardLevel
		mysql.CreateTime = vul.CreateTime
		mysql.PublishTime = vul.PublishTime
		mysql.UpdateTime = vul.UpdateTime
		mysql.TypeName = vul.TypeName
		mysql.VulType = vul.VulType
		db.Create(&mysql)
		log.Printf("【%s】store %s successfully", r.Name(), mysql.CnnvdCode)
	}
	return nil
}

// GetLatestCNNVD 获取最新的漏洞编号
func (r *ReqVulList) GetLatestCNNVD() (string, error) {
	//请求漏洞列表
	r.PageSize = 10
	resList, err := Post[*ResVulList](r, utils.FormatURL(Domain, APIVulList))
	if err != nil {
		return "", xerrors.Errorf("【%s】fail to fetch:%w\n", r.Name(), err)
	}

	// 设置第一个漏洞为最新的漏洞
	var latestCNNVD string
	latestCNNVD = resList.Data.Records[0].CnnvdCode
	for i := 1; i <= r.PageSize-1; i++ {
		latestCNNVD, err = LatestCNNVD(resList.Data.Records[i].CnnvdCode, latestCNNVD)
		if err != nil {
			return "", xerrors.Errorf("【%s】fail to get latest cnnvd:%w\n", r.Name(), err)
		}
	}
	return latestCNNVD, nil
}

func (r *ReqVulList) getPageNum() (num int, err error) {
	result, err := Post[*ResVulList](r, utils.FormatURL(Domain, APIVulList))
	if err != nil {
		return 0, xerrors.Errorf("【%s】fail to fetch:%w\n", r.Name(), err)
	}
	return int(math.Ceil(float64(result.Data.Total) / float64(r.PageSize))), nil
}

func (r *ReqVulList) read(file string) (*TableVulList, error) {
	data, err := os.ReadFile(file)
	if err != nil {
		return nil, xerrors.Errorf("fail to read :%w", err)
	}
	var vulList TableVulList
	err = json.Unmarshal(data, &vulList)
	if err != nil {
		return nil, xerrors.Errorf("fail to unmarshal:%w", err)
	}
	return &vulList, nil
}
