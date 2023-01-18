package api

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func TestReqVulList_Fetch(t *testing.T) {
	type fields struct {
		PageIndex   int
		PageSize    int
		Keyword     string
		HazardLevel string
		VulType     string
		Vendor      string
		Product     string
		DateType    string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *[]Record
		wantErr bool
	}{
		{
			name:   "happy test for a few vul",
			fields: fields{Keyword: "CNNVD-198801-"},
			want: &[]Record{
				{
					Id:          "d19706219d1648da9a8d008eb3a7aeec",
					VulName:     "ftpd CWD ~root命令漏洞",
					CnnvdCode:   "CNNVD-198801-002",
					CveCode:     "CVE-1999-0082",
					HazardLevel: 1,
					CreateTime:  "2022-09-21",
					PublishTime: "1988-11-11",
					UpdateTime:  "2010-12-03",
					VulType:     "0",
				},
				{
					Id:          "f93433d1575641779b48f03fbd35c4ef",
					VulName:     "Berkeley Sendmail 5.x DEBUG远程执行任意命令漏洞",
					CnnvdCode:   "CNNVD-198801-001",
					CveCode:     "CVE-1999-0095",
					HazardLevel: 1,
					CreateTime:  "2022-09-21",
					PublishTime: "1988-10-01",
					UpdateTime:  "2019-06-12",
					TypeName:    "其他",
					VulType:     "0",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqVulList{
				PageIndex:   tt.fields.PageIndex,
				PageSize:    tt.fields.PageSize,
				Keyword:     tt.fields.Keyword,
				HazardLevel: tt.fields.HazardLevel,
				VulType:     tt.fields.VulType,
				Vendor:      tt.fields.Vendor,
				Product:     tt.fields.Product,
				DateType:    tt.fields.DateType,
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

func TestReqVulList_Save(t *testing.T) {
	type fields struct {
		PageIndex   int
		PageSize    int
		Keyword     string
		HazardLevel string
		VulType     string
		Vendor      string
		Product     string
		DateType    string
	}
	type args struct {
		data *[]Record
		dir  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "happy test for a few vuls",
			fields: fields{},
			args: args{
				data: &[]Record{
					{
						Id:          "d19706219d1648da9a8d008eb3a7aeec",
						VulName:     "ftpd CWD ~root命令漏洞",
						CnnvdCode:   "CNNVD-198801-002",
						CveCode:     "CVE-1999-0082",
						HazardLevel: 1,
						CreateTime:  "2022-09-21",
						PublishTime: "1988-11-11",
						UpdateTime:  "2010-12-03",
						VulType:     "0",
					},
					{
						Id:          "f93433d1575641779b48f03fbd35c4ef",
						VulName:     "Berkeley Sendmail 5.x DEBUG远程执行任意命令漏洞",
						CnnvdCode:   "CNNVD-198801-001",
						CveCode:     "CVE-1999-0095",
						HazardLevel: 1,
						CreateTime:  "2022-09-21",
						PublishTime: "1988-10-01",
						UpdateTime:  "2019-06-12",
						TypeName:    "其他",
						VulType:     "0",
					},
				},
				dir: "testdata",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqVulList{
				PageIndex:   tt.fields.PageIndex,
				PageSize:    tt.fields.PageSize,
				Keyword:     tt.fields.Keyword,
				HazardLevel: tt.fields.HazardLevel,
				VulType:     tt.fields.VulType,
				Vendor:      tt.fields.Vendor,
				Product:     tt.fields.Product,
				DateType:    tt.fields.DateType,
			}
			if err := r.Save(tt.args.data, tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReqVulList_StoreByFile(t *testing.T) {
	type fields struct {
		PageIndex   int
		PageSize    int
		Keyword     string
		HazardLevel string
		VulType     string
		Vendor      string
		Product     string
		DateType    string
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
			name: "happy test for a few vul",
			args: args{
				db:  db,
				dir: "testdata",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqVulList{
				PageIndex:   tt.fields.PageIndex,
				PageSize:    tt.fields.PageSize,
				Keyword:     tt.fields.Keyword,
				HazardLevel: tt.fields.HazardLevel,
				VulType:     tt.fields.VulType,
				Vendor:      tt.fields.Vendor,
				Product:     tt.fields.Product,
				DateType:    tt.fields.DateType,
			}
			if err := r.StoreByFile(tt.args.db, tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("StoreByFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReqVulList_StoreByRequest(t *testing.T) {
	type fields struct {
		PageIndex   int
		PageSize    int
		Keyword     string
		HazardLevel string
		VulType     string
		Vendor      string
		Product     string
		DateType    string
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
			name: "happy test for many vuls",
			fields: fields{
				PageSize: 100,
				Keyword:  "CNNVD-202301-",
			},
			args:    args{db},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqVulList{
				PageIndex:   tt.fields.PageIndex,
				PageSize:    tt.fields.PageSize,
				Keyword:     tt.fields.Keyword,
				HazardLevel: tt.fields.HazardLevel,
				VulType:     tt.fields.VulType,
				Vendor:      tt.fields.Vendor,
				Product:     tt.fields.Product,
				DateType:    tt.fields.DateType,
			}
			if err := r.StoreByRequest(tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("StoreByRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
