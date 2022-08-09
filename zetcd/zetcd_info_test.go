package zetcd

import (
	"reflect"
	"testing"
)

func TestKeyInfo_GetRegisterKey(t *testing.T) {
	type fields struct {
		Cluster string
		Name    string
		Kind    string
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "GetRegisterKey",
			fields: fields{
				Cluster: "prod",
				Name:    "book",
				Kind:    "grpc",
			},
			want: "prod/book:grpc",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &KeyInfo{
				Cluster: tt.fields.Cluster,
				Name:    tt.fields.Name,
				Kind:    tt.fields.Kind,
			}
			if got := i.GetRegisterKey(); got != tt.want {
				t.Errorf("GetRegisterKey() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestKeyInfo_SetRegisterKey(t *testing.T) {
	type fields struct {
		Cluster string
		Name    string
		Kind    string
	}
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "SetRegisterKey",
			fields: fields{
				Cluster: "",
				Name:    "",
				Kind:    "",
			},
			args: args{address: "prod/book:grpc"},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &KeyInfo{
				Cluster: tt.fields.Cluster,
				Name:    tt.fields.Name,
				Kind:    tt.fields.Kind,
			}
			if err := i.SetRegisterKey(tt.args.address); (err != nil) != tt.wantErr {
				t.Errorf("SetRegisterKey() error = %v, wantErr %v", err, tt.wantErr)
			}
			if i.Cluster != "prod" || i.Name!="book" || i.Kind!="grpc"{
				t.Errorf("SetRegisterKey() error, keyinfo = %+v", i)
			}
 		})
	}
}

func TestKeyInfo_setCluster(t *testing.T) {
	type fields struct {
		Cluster string
		Name    string
		Kind    string
	}
	type args struct {
		cluster string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "setCluster",
			args: args{cluster: "prod"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &KeyInfo{
				Cluster: tt.fields.Cluster,
				Name:    tt.fields.Name,
				Kind:    tt.fields.Kind,
			}
			i.setCluster(tt.args.cluster)
			if i.Cluster != "prod"{
				t.Errorf("setCluster error cluster=%s not eq prod",i.Cluster)
			}
		})
	}
}

func TestKeyInfo_setKind(t *testing.T) {
	type fields struct {
		Cluster string
		Name    string
		Kind    string
	}
	type args struct {
		kind string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "setKind",
			args: args{kind: "grpc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &KeyInfo{
				Cluster: tt.fields.Cluster,
				Name:    tt.fields.Name,
				Kind:    tt.fields.Kind,
			}
			i.setKind(tt.args.kind)
			if i.Kind != "grpc" {
				t.Errorf("setKind error kind not eq grpc")
			}
		})
	}
}

func TestKeyInfo_setName(t *testing.T) {
	type fields struct {
		Cluster string
		Name    string
		Kind    string
	}
	type args struct {
		name string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "setName",
			args: args{name: "book"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &KeyInfo{
				Cluster: tt.fields.Cluster,
				Name:    tt.fields.Name,
				Kind:    tt.fields.Kind,
			}
			i.setName(tt.args.name)
			if i.Name != "book"{
				t.Errorf("setName error name not eq book")
			}
		})
	}
}

func TestValueInfo_EncodeValue(t *testing.T) {
	type fields struct {
		Kind        string
		Ip          string
		Port        uint32
		Status      uint32
		RequestFlow uint32
		UpdatedAt   uint32
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name:"EncodeValue",
			fields: fields{
				Kind:        "grpc",
				Ip:          "192.168.0.1",
				Port:        3307,
				Status:      0,
				RequestFlow: 0,
				UpdatedAt:   10,
			},
			want:[]byte{10, 4, 103, 114, 112, 99, 18, 11, 49, 57, 50, 46, 49, 54, 56, 46, 48, 46, 49, 24, 235, 25, 32, 0, 40, 0, 48, 10},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ValueInfo{
				Kind:        tt.fields.Kind,
				Ip:          tt.fields.Ip,
				Port:        tt.fields.Port,
				Status:      tt.fields.Status,
				RequestFlow: tt.fields.RequestFlow,
				UpdatedAt:   tt.fields.UpdatedAt,
			}
			got, err := i.EncodeValue()
			if (err != nil) != tt.wantErr {
				t.Errorf("EncodeValue() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncodeValue() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueInfo_SetRegisterAddress(t *testing.T) {
	type fields struct {
		Kind        string
		Ip          string
		Port        uint32
		Status      uint32
		RequestFlow uint32
		UpdatedAt   uint32
	}
	type args struct {
		address string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "SetRegisterAddress",
			fields: fields{
				Kind:        "",
				Ip:          "",
				Port:        0,
				Status:      0,
				RequestFlow: 0,
				UpdatedAt:   0,
			},
			args: args{address: "192.168.0.1:3307"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ValueInfo{
				Kind:        tt.fields.Kind,
				Ip:          tt.fields.Ip,
				Port:        tt.fields.Port,
				Status:      tt.fields.Status,
				RequestFlow: tt.fields.RequestFlow,
				UpdatedAt:   tt.fields.UpdatedAt,
			}
			if err := i.SetRegisterAddress(tt.args.address); (err != nil) != tt.wantErr {
				t.Errorf("SetRegisterAddress() error = %v, wantErr %v", err, tt.wantErr)
			}
			if i.getRegisterAddress() != "192.168.0.1:3307"{
				t.Errorf("SetRegisterAddress() error  wantAddress 192.168.0.1:3307, got=%v", i.getRegisterAddress() )
			}
		})
	}
}

func TestValueInfo_DecodeValue(t *testing.T) {
	type fields struct {
		Kind        string
		Ip          string
		Port        uint32
		Status      uint32
		RequestFlow uint32
		UpdatedAt   uint32
	}
	type args struct {
		infoByte []byte
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "DecodeValue",
			args: args{infoByte: []byte{10, 4, 103, 114, 112, 99, 18, 11, 49, 57, 50, 46, 49, 54, 56, 46, 48, 46, 49, 24, 235, 25, 32, 0, 40, 0, 48, 10}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ValueInfo{
				Kind:        "",
				Ip:          "",
				Port:        0,
				Status:      0,
				RequestFlow: 0,
				UpdatedAt:   0,
			}
			if err := i.DecodeValue(tt.args.infoByte); (err != nil) != tt.wantErr {
				t.Errorf("UndecodeValue() error = %v, wantErr %v", err, tt.wantErr)
			}
			if i.Kind != "grpc" || i.Ip != "192.168.0.1" || i.Port != 3307 || i.Status != 0 || i.RequestFlow != 0 || i.UpdatedAt != 10{
				t.Errorf("decodeValue() error, i= %+v", i)
			}
			t.Logf("decode: %+v",i)
		})
	}
}

func TestValueInfo_getRegisterValue(t *testing.T) {
	type fields struct {
		Kind        string
		Ip          string
		Port        uint32
		Status      uint32
		RequestFlow uint32
		UpdatedAt   uint32
	}
	tests := []struct {
		name   string
		fields fields
		want   string
	}{
		{
			name: "getRegisterAddress",
			fields: fields{
				Kind:        "grpc",
				Ip:          "192.168.0.1",
				Port:        3307,
				Status:      0,
				RequestFlow: 0,
				UpdatedAt:   10,
			},
			want: "192.168.0.1:3307",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ValueInfo{
				Kind:        tt.fields.Kind,
				Ip:          tt.fields.Ip,
				Port:        tt.fields.Port,
				Status:      tt.fields.Status,
				RequestFlow: tt.fields.RequestFlow,
				UpdatedAt:   tt.fields.UpdatedAt,
			}
			if got := i.getRegisterAddress(); got != tt.want {
				t.Errorf("getRegisterAddress() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestValueInfo_setIp(t *testing.T) {
	type fields struct {
		Kind        string
		Ip          string
		Port        uint32
		Status      uint32
		RequestFlow uint32
		UpdatedAt   uint32
	}
	type args struct {
		ip string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "setIp",
			args: args{ip: "192.168.0.1"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ValueInfo{
				Kind:        tt.fields.Kind,
				Ip:          tt.fields.Ip,
				Port:        tt.fields.Port,
				Status:      tt.fields.Status,
				RequestFlow: tt.fields.RequestFlow,
				UpdatedAt:   tt.fields.UpdatedAt,
			}
			i.setIp(tt.args.ip)
			if i.Ip != "192.168.0.1"{
				t.Errorf("setIp error ip not eq 192.168.0.1")
			}
		})
	}
}

func TestValueInfo_setPort(t *testing.T) {
	type fields struct {
		Kind        string
		Ip          string
		Port        uint32
		Status      uint32
		RequestFlow uint32
		UpdatedAt   uint32
	}
	type args struct {
		port uint32
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "setPort",
			args: args{port: 3307},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			i := &ValueInfo{
				Kind:        tt.fields.Kind,
				Ip:          tt.fields.Ip,
				Port:        tt.fields.Port,
				Status:      tt.fields.Status,
				RequestFlow: tt.fields.RequestFlow,
				UpdatedAt:   tt.fields.UpdatedAt,
			}
			i.setPort(tt.args.port)
			if i.Port != 3307 {
				t.Errorf("setPort error port not eq 3307")
			}
		})
	}
}
