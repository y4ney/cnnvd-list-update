package api

import (
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"testing"
)

func TestReqVulType_StoreByFile(t *testing.T) {
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
		args    args
		wantErr bool
	}{
		{
			name:    "happy test",
			args:    args{db, "testdata"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqVulType{}
			if err := r.StoreByFile(tt.args.db, tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("StoreByFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReqVulType_StoreByRequest(t *testing.T) {
	type args struct {
		db *gorm.DB
	}
	db, err := gorm.Open(mysql.Open(DatabaseSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "happy test for all vul type",
			args:    args{db},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqVulType{}
			if err := r.StoreByRequest(tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("StoreByRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
