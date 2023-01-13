package api

import (
	"github.com/0yaney0/cnnvd-list-update/utils"
	"golang.org/x/xerrors"
	"log"
	"path/filepath"
)

// ReqProduct 产品选择列表请求参数
type ReqProduct struct {
	ProductKeyword string `json:"productKeyword"` // 产品关键词
}

// ResProduct 产品选择列表响应参数
type ResProduct struct {
	ResCode           // 响应码
	Data    []Product `json:"data,omitempty"` // 产品列表
}

// Product 产品列表
type Product struct {
	Label string `json:"label,omitempty"`
	Value string `json:"value,omitempty"`
}

func (p *ReqProduct) Fetch() error {
	resProduct, err := Post[*ResProduct](p, utils.FormatURL(Domain, APIProduct))
	if err != nil {
		return xerrors.Errorf("fail to request Product:%w\n", err)
	}
	err = utils.Write(filepath.Join(utils.DefaultCacheDir(), "product.json"), resProduct.Data)
	if err != nil {
		return xerrors.Errorf("fail to save Product:%w\n", err)
	}
	log.Printf("save %s successfully", filepath.Join(utils.CNNVDListDir(), "product.json"))
	return nil
}
