package api

import (
	"encoding/json"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"os"
	"reflect"
	"testing"
)

func GetTestData1() (*[]HazardLevel, error) {
	var level []HazardLevel
	data, err := os.ReadFile("./testdata/hazard_level_example.json")
	if err != nil {
		return nil, err
	}
	err = json.Unmarshal(data, &level)
	if err != nil {
		return nil, err
	}
	return &level, nil
}
func TestReqHazardLevel_Fetch(t *testing.T) {
	level, err := GetTestData1()
	if err != nil {
		panic(err)
	}
	tests := []struct {
		name    string
		want    *[]HazardLevel
		wantErr bool
	}{
		{
			name:    "happy test",
			want:    level,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqHazardLevel{}
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

func TestReqHazardLevel_Save(t *testing.T) {
	level, err := GetTestData1()
	if err != nil {
		panic(err)
	}
	type args struct {
		data *[]HazardLevel
		dir  string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy",
			args: args{
				data: level,
				dir:  "testdata",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqHazardLevel{}
			if err := r.Save(tt.args.data, tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("Save() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReqHazardLevel_StoreByFile(t *testing.T) {
	db, err := gorm.Open(mysql.Open(DatabaseSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	type args struct {
		db  *gorm.DB
		dir string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "happy test",
			args: args{
				db:  db,
				dir: "testdata",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqHazardLevel{}
			if err := r.StoreByFile(tt.args.db, tt.args.dir); (err != nil) != tt.wantErr {
				t.Errorf("StoreByFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestReqHazardLevel_StoreByRequest(t *testing.T) {
	db, err := gorm.Open(mysql.Open(DatabaseSource), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	type args struct {
		db *gorm.DB
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "happy test",
			args:    args{db},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &ReqHazardLevel{}
			if err := r.StoreByRequest(tt.args.db); (err != nil) != tt.wantErr {
				t.Errorf("StoreByRequest() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
