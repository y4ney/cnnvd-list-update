package api

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"reflect"
	"testing"
)

func TestReqProduct_Fetch(t *testing.T) {
	type fields struct {
		ProductKeyword string
	}
	tests := []struct {
		name    string
		fields  fields
		want    *[]Product
		wantErr bool
	}{
		{
			name:   "happy test one product",
			fields: fields{ProductKeyword: "PowerDNS Recursor"},
			want: &[]Product{
				{
					Label: "1000",
					Value: "PowerDNS Recursor",
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqProduct{
				ProductKeyword: tt.fields.ProductKeyword,
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

func TestReqProduct_Save(t *testing.T) {
	type fields struct {
		ProductKeyword string
	}
	type args struct {
		data *[]Product
		dir  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "happy test one product",
			args: args{
				data: &[]Product{
					{
						Label: "1000",
						Value: "PowerDNS Recursor",
					},
				},
				dir: "testdata",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqProduct{
				ProductKeyword: tt.fields.ProductKeyword,
			}
			if err := r.Save(tt.args.data, tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReqProduct_StoreByFile(t *testing.T) {
	type fields struct {
		ProductKeyword string
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
			name: "happy test a product",
			args: args{
				db:  db,
				dir: "testdata",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqProduct{
				ProductKeyword: tt.fields.ProductKeyword,
			}
			if err := r.StoreByFile(tt.args.db, tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("StoreByFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReqProduct_StoreByRequest(t *testing.T) {
	type fields struct {
		ProductKeyword string
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
			name:    "happy test for all product",
			fields:  fields{ProductKeyword: ""},
			args:    args{db},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqProduct{
				ProductKeyword: tt.fields.ProductKeyword,
			}
			if err := r.StoreByRequest(tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("StoreByRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
