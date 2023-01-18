package api

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func TestReqVulDetail_Fetch(t *testing.T) {
	type fields struct {
		Id        string
		VulType   string
		CnnvdCode string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *VulDetail
		wantErr bool
	}{
		{
			name: "happy test for one vul",
			fields: fields{
				Id:        "07098e506fb9488188edea1372908d09",
				VulType:   "0",
				CnnvdCode: "CNNVD-202212-2756",
			},
			want: &VulDetail{
				CNNVDDetail: CNNVDDetail{
					VulName:        "Containous Traefik 日志信息泄露漏洞",
					CnnvdCode:      "CNNVD-202212-2756",
					CveCode:        "CVE-2022-23469",
					PublishTime:    "2022-12-08 00:00:00",
					IsOfficial:     1,
					Vendor:         "1000971",
					HazardLevel:    3,
					VulType:        "日志信息泄露",
					VulTypeName:    "日志信息泄露",
					VulDesc:        "Containous Traefik是美国Containous公司的一款反向代理和负载平衡器。\r\nContainous Traefik 2.9.6之前的版本存在日志信息泄露漏洞，该漏洞源于其调试日志中显示授权标头。",
					AffectedVendor: "Containous",
					ReferUrl:       "来源:MISC\r\n链接:https://github.com/traefik/traefik/pull/9574\r\n\r\n来源:MISC\r\n链接:https://github.com/traefik/traefik/releases/tag/v2.9.6\r\n\r\n来源:MISC\r\n链接:https://github.com/traefik/traefik/security/advisories/GHSA-h2ph-vhm7-g4hp\r\n\r\n来源:cxsecurity.com\r\n链接:https://cxsecurity.com/cveshow/CVE-2022-23469/",
					Patch:          "https://github.com/traefik/traefik/security/advisories/GHSA-h2ph-vhm7-g4hp",
					UpdateTime:     "2022-12-13 00:00:00",
					CnnvdFiledShow: "vul_name,cnnvd_code,_code,publish_time,is_official,vendor,hazard_level,vul_type,vul_desc,affected_product,affected_vendor,product_desc,affected_system,refer_url,patch_id,product,update_time,patch",
					Varchar1:       "日志信息泄露",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqVulDetail{
				Id:        tt.fields.Id,
				VulType:   tt.fields.VulType,
				CnnvdCode: tt.fields.CnnvdCode,
			}
			got, err := r.Fetch()
			if (err != nil) != tt.wantErr {
				t.Errorf("Fetch() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Fetch() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestReqVulDetail_Save(t *testing.T) {
	type fields struct {
		Id        string
		VulType   string
		CnnvdCode string
	}
	type args struct {
		data *VulDetail
		dir  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy test for a vul",
			args: args{
				data: &VulDetail{
					CNNVDDetail: CNNVDDetail{
						VulName:        "Containous Traefik 日志信息泄露漏洞",
						CnnvdCode:      "CNNVD-202212-2756",
						CveCode:        "CVE-2022-23469",
						PublishTime:    "2022-12-08 00:00:00",
						IsOfficial:     1,
						Vendor:         "1000971",
						HazardLevel:    3,
						VulType:        "日志信息泄露",
						VulTypeName:    "日志信息泄露",
						VulDesc:        "Containous Traefik是美国Containous公司的一款反向代理和负载平衡器。\r\nContainous Traefik 2.9.6之前的版本存在日志信息泄露漏洞，该漏洞源于其调试日志中显示授权标头。",
						AffectedVendor: "Containous",
						ReferUrl:       "来源:MISC\r\n链接:https://github.com/traefik/traefik/pull/9574\r\n\r\n来源:MISC\r\n链接:https://github.com/traefik/traefik/releases/tag/v2.9.6\r\n\r\n来源:MISC\r\n链接:https://github.com/traefik/traefik/security/advisories/GHSA-h2ph-vhm7-g4hp\r\n\r\n来源:cxsecurity.com\r\n链接:https://cxsecurity.com/cveshow/CVE-2022-23469/",
						Patch:          "https://github.com/traefik/traefik/security/advisories/GHSA-h2ph-vhm7-g4hp",
						UpdateTime:     "2022-12-13 00:00:00",
						CnnvdFiledShow: "vul_name,cnnvd_code,_code,publish_time,is_official,vendor,hazard_level,vul_type,vul_desc,affected_product,affected_vendor,product_desc,affected_system,refer_url,patch_id,product,update_time,patch",
						Varchar1:       "日志信息泄露",
					},
				},
				dir: "testdata",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqVulDetail{
				Id:        tt.fields.Id,
				VulType:   tt.fields.VulType,
				CnnvdCode: tt.fields.CnnvdCode,
			}
			if err := r.Save(tt.args.data, tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReqVulDetail_StoreByFile(t *testing.T) {
	type fields struct {
		Id        string
		VulType   string
		CnnvdCode string
	}
	type args struct {
		db  *gorm.DB
		dir string
	}
	db, err := gorm.Open(mysql.Open(DatabaseSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy test for a vuln",
			args: args{
				db:  db,
				dir: "testdata",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqVulDetail{
				Id:        tt.fields.Id,
				VulType:   tt.fields.VulType,
				CnnvdCode: tt.fields.CnnvdCode,
			}
			if err := r.StoreByFile(tt.args.db, tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("StoreByFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReqVulDetail_StoreByRequest(t *testing.T) {
	type fields struct {
		Id        string
		VulType   string
		CnnvdCode string
	}
	type args struct {
		db *gorm.DB
	}
	db, err := gorm.Open(mysql.Open(DatabaseSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy test for a vul",
			fields: fields{
				Id:        "d19706219d1648da9a8d008eb3a7aeec",
				VulType:   "0",
				CnnvdCode: "CNNVD-198801-002",
			},
			args:    args{db},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqVulDetail{
				Id:        tt.fields.Id,
				VulType:   tt.fields.VulType,
				CnnvdCode: tt.fields.CnnvdCode,
			}
			if err := r.StoreByRequest(tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("StoreByRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
