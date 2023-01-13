package api

import (
	"github.com/0yaney0/cnnvd-list-update/utils"
	"golang.org/x/xerrors"
	"log"
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
	Label string `json:"label,omitempty"`
	Value string `json:"value,omitempty"`
}

func (v *ReqVendor) Fetch() error {
	// 获取供应商信息
	resVendor, err := Post[*ResVendor](v, utils.FormatURL(Domain, APIVendor))
	if err != nil {
		return xerrors.Errorf("fail to request Vendor:%w\n", err)
	}
	err = utils.Write(filepath.Join(utils.DefaultCacheDir(), "vendor.json"), resVendor.Data)
	if err != nil {
		return xerrors.Errorf("fail to save Vendor:%w\n", err)
	}
	log.Printf("save %s successfully", filepath.Join(utils.CNNVDListDir(), "vendor.json"))
	return nil
}
