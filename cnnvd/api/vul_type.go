package api

import (
	"github.com/0yaney0/cnnvd-list-update/utils"
	"golang.org/x/xerrors"
	"log"
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

func (v *ReqVulType) Fetch() error {
	// 获取漏洞类型
	resVulType, err := Post[*ResVulType](ReqVulType{}, utils.FormatURL(Domain, APIVulType))
	if err != nil {
		return xerrors.Errorf("fail to request VulType:%w\n", err)
	}
	err = utils.Write(filepath.Join(utils.DefaultCacheDir(), "vul_type.json"), resVulType.Data)
	if err != nil {
		return xerrors.Errorf("fail to save VulType:%w\n", err)
	}
	log.Printf("save %s successfully", filepath.Join(utils.CNNVDListDir(), "vul_type.json"))
	return nil
}
