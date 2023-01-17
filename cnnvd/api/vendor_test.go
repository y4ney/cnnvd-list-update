package api

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func TestReqVendor_Fetch(t *testing.T) {
	type fields struct {
		VendorKeyword string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *[]Vendor
		wantErr bool
	}{
		{
			name:   "happy path",
			fields: fields{"北京智慧远景科技产业有限公司"},
			want: &[]Vendor{
				{
					Label: "1000006",
					Value: "北京智慧远景科技产业有限公司",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ReqVendor{
				VendorKeyword: tt.fields.VendorKeyword,
			}
			got, err := v.Fetch()
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

func TestReqVendor_Save(t *testing.T) {
	type fields struct {
		VendorKeyword string
	}
	type args struct {
		data *[]Vendor
		dir  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "test save function",
			fields: fields{},
			args: args{
				data: &[]Vendor{{
					Label: "1000006",
					Value: "北京智慧远景科技产业有限公司",
				}},
				dir: "testdata",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ReqVendor{
				VendorKeyword: tt.fields.VendorKeyword,
			}
			if err := v.Save(tt.args.data, tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReqVendor_StoreByFile(t *testing.T) {
	type fields struct {
		VendorKeyword string
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
			name: "test store by file function",
			args: args{
				db:  db,
				dir: "testdata",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ReqVendor{
				VendorKeyword: tt.fields.VendorKeyword,
			}
			if err := v.StoreByFile(tt.args.db, tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("StoreByFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReqVendor_StoreByRequest(t *testing.T) {
	type fields struct {
		VendorKeyword string
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
			name:    "test store by request",
			fields:  fields{"Docker"},
			args:    args{db},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &ReqVendor{
				VendorKeyword: tt.fields.VendorKeyword,
			}
			if err := v.StoreByRequest(tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("StoreByRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
