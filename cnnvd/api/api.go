package api

import (
	"encoding/json"
	"fmt"
	"github.com/0yaney0/cnnvd-list-update/utils"
	"golang.org/x/xerrors"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

const (
	Domain         = "www.cnnvd.org.cn"
	APIVulDetail   = "web/cnnvdVul/getCnnnvdDetailOnDatasource"
	APIVulList     = "web/homePage/cnnvdVulList"
	APIVendor      = "web/homePage/getVendorSelectList"
	APIVulType     = "web/homePage/vulTypeList"
	APIProduct     = "web/homePage/getProductSelectList"
	PageSize       = 100
	Retry          = 5
	FirstYear      = 1988
	DatabaseSource = "root:1600850141yangli@tcp(localhost:3306)/cnnvd?charset=utf8mb4&parseTime=true&loc=Asia%2FShanghai"
	ProductName    = "cnnvd/api/product.go"
	VendorName     = "cnnvd/api/vendor.go"
	VulDetailName  = "cnnvd/api/vul_detail.go"
	VulListName    = "cnnvd/api/vul_list.go"
	ProductTable   = "product"
	VendorTable    = "vendor"
	VulListTable   = "vul_list"
	VulDetailTable = "vul_detail"
	VendorFile     = "vendor.json"
	ProductFile    = "product.json"
	VulListFile    = "vul_list"
	VulDetailFile  = "vul_detail"
)

// ResCode 响应码
type ResCode struct {
	Code    int    `json:"code,omitempty"`    // 代码
	Success bool   `json:"success,omitempty"` // 是否成功
	Message string `json:"message,omitempty"` // 消息
	Time    string `json:"time,omitempty"`    // 时间
}
type CNNVD struct {
	Year  string
	Month string
	ID    int
}

func NewCNNVD(str string) (*CNNVD, error) {
	s := strings.Split(str, "-")
	if len(s) != 3 {
		return nil, xerrors.Errorf("invalid CNNVD-ID format: %s\n", str)
	}

	id, err := strconv.Atoi(s[2])
	if err != nil {
		return nil, xerrors.Errorf("fail to convert %s's id:%w\n", str, err)
	}
	return &CNNVD{
		Year:  s[1][:4],
		Month: s[1][4:],
		ID:    id,
	}, nil
}

func (c *CNNVD) GetDate() (*time.Time, error) {
	date, err := time.Parse("2006-01", fmt.Sprintf("%v-%v", c.Year, c.Month))
	if err != nil {
		return nil, xerrors.Errorf("fail to get date:%w\n", err)
	}
	return &date, nil
}

func (c *CNNVD) Before(item *CNNVD) (bool, error) {
	cDate, err := c.GetDate()
	if err != nil {
		return false, err
	}
	itemDate, err := item.GetDate()
	if err != nil {
		return false, err
	}
	if cDate.After(*itemDate) {
		return false, nil
	}
	if cDate.Before(*itemDate) {
		return true, nil
	}
	if c.ID < item.ID {
		return true, nil
	}
	return false, nil

}

func (c *CNNVD) After(item *CNNVD) (bool, error) {
	cDate, err := c.GetDate()
	if err != nil {
		return false, err
	}
	itemDate, err := item.GetDate()
	if err != nil {
		return false, err
	}
	if cDate.Before(*itemDate) {
		return false, nil
	}
	if cDate.After(*itemDate) {
		return true, nil
	}
	if c.ID > item.ID {
		return true, nil
	}
	return false, nil
}

func (c *CNNVD) Equal(item *CNNVD) bool {
	return c.Year == item.Year && c.Month == item.Month && c.ID == item.ID
}

func LatestCNNVD(str1, str2 string) (string, error) {
	cnnvd1, err := NewCNNVD(str1)
	if err != nil {
		return "", xerrors.Errorf("fail to new %s:%w\n", str1, err)
	}
	cnnvd2, err := NewCNNVD(str2)
	if err != nil {
		return "", xerrors.Errorf("fail to new %s:%w\n", str2, err)
	}
	flag, err := cnnvd1.After(cnnvd2)
	if err != nil {
		return "", err
	}
	if flag {
		return str1, nil
	}
	return str2, nil
}

func Post[T *ResVendor | *ResVulType | *ResVulDetail | *ResVulList | *ResProduct](req any, url string) (res T, err error) {
	data, err := utils.FormatBody(req)
	if err != nil {
		return nil, xerrors.Errorf("fail to format request body:%w\n", err)
	}
	resBody, err := utils.Post(url, data, Retry)
	if err != nil {
		return nil, xerrors.Errorf("fail to post data:%w\n", err)
	}
	err = json.Unmarshal(resBody, &res)
	if err != nil {
		return nil, xerrors.Errorf("fail to unmarshal response body:%w\n", err)
	}
	return res, err
}

// SaveCNNVDPerYear 存储每年的漏洞
func SaveCNNVDPerYear(dirPath string, cnnvdID string, data interface{}) error {
	cnnvd, err := NewCNNVD(cnnvdID)
	if err != nil {
		return xerrors.Errorf("fail to new %s:%w\n", cnnvdID, err)
	}

	yearDir := filepath.Join(dirPath, cnnvd.Year)
	monthDir := filepath.Join(yearDir, cnnvd.Month)
	if err = os.MkdirAll(monthDir, os.ModePerm); err != nil {
		return err
	}

	filePath := filepath.Join(monthDir, fmt.Sprintf("%s.json", cnnvdID))
	if err = utils.Write(filePath, data); err != nil {
		return xerrors.Errorf("failed to write file: %w\n", err)
	}
	return nil
}

func FormatKeyword(year int, month string) string {
	// TODO 将year和month先解析为时间，然后在format成keyword
	return fmt.Sprintf("CNNVD-%v%s-", year, month)
}
