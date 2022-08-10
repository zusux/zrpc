package zetcd

import (
	"reflect"
	"testing"
)

func TestNewRoundRobin(t *testing.T) {
	tests := []struct {
		name string
		want *roundRobin
	}{
		{
			name: "NewRoundRobin",
			want: &roundRobin{
				c: 0,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRoundRobin(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRoundRobin() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_roundRobin_GetPoint(t *testing.T) {
	type fields struct {
		c uint64
	}
	f := fields{c: 0}
	type args struct {
		count int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []int
		wantErr bool
	}{
		{
			name: "GetPoint",
			fields: f,
			args: args{count: 1000},
			want: []int{0,1,2,3,4,5,6},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &roundRobin{
				c: tt.fields.c,
			}
			for _,v := range tt.want{
				got, err := r.GetPoint(tt.args.count)
				if (err != nil) != tt.wantErr {
					t.Errorf("GetPoint() error = %v, wantErr %v", err, tt.wantErr)
					return
				}
				t.Logf("GetPoint() got = %v, want %v", got, v)
				if got != v {
					t.Errorf("GetPoint() got = %v, want %v", got, v)
				}
			}

		})
	}
}
