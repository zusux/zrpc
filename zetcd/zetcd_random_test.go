package zetcd

import (
	"math/rand"
	"reflect"
	"testing"
	"time"
)

func TestNewRandom(t *testing.T) {
	type args struct {
		seed int64
	}
	seed := time.Now().Unix()
	tests := []struct {
		name string
		args args
		want *random
	}{
		{
			name: "NewRandom",
			args: args{seed: seed},
			want: NewRandom(seed),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewRandom(tt.args.seed); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewRandom() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_random_GetPoint(t *testing.T) {
	type fields struct {
		r *rand.Rand
	}
	seed := time.Now().Unix()
	type args struct {
		count int
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int
		wantErr bool
	}{
		{
			name: "getPoint",
			fields: fields{r:rand.New(rand.NewSource(seed))},
			args: args{count: 10},
			want: 20,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := &random{
				r: tt.fields.r,
			}
			got, err := r.GetPoint(tt.args.count)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPoint() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got == tt.want {
				t.Errorf("GetPoint() got = %v, want %v", got, tt.want)
			}
			t.Logf("GetPoint() got = %v, want %v", got, tt.want)
		})
	}
}
