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
		want    *[]TableProduct
		wantErr bool
	}{
		{
			name:   "happy path",
			fields: fields{"PowerDNS Recursor"},
			want: &[]TableProduct{
				{
					Product: Product{
						Label: "1000",
						Value: "PowerDNS Recursor",
					},
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ReqProduct{
				ProductKeyword: tt.fields.ProductKeyword,
			}
			got, err := p.Fetch()
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
			name:   "happy path",
			fields: fields{""},
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
			p := &ReqProduct{
				ProductKeyword: tt.fields.ProductKeyword,
			}
			if err := p.Save(tt.args.data, tt.args.dir); (err != nil) != tt.wantErr {
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
			name:   "happy path",
			fields: fields{""},
			args: args{
				db:  db,
				dir: "testdata",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ReqProduct{
				ProductKeyword: tt.fields.ProductKeyword,
			}
			if err := p.StoreByFile(tt.args.db, tt.args.dir); (err != nil) != tt.wantErr {
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
			name:    "happy path",
			fields:  fields{"Croogo"},
			args:    args{db},
			wantErr: false,
		},
		{
			name:    "sad path",
			fields:  fields{"PowerDNS Recursor"},
			args:    args{db},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &ReqProduct{
				ProductKeyword: tt.fields.ProductKeyword,
			}
			if err := p.StoreByRequest(tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("StoreByRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
