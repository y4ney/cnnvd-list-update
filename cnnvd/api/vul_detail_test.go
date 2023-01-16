package api

import (
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
		want    *TableVulDetail
		wantErr bool
	}{
		// TODO: Add test cases.
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
		data *TableVulList
		dir  string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		// TODO: Add test cases.
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
