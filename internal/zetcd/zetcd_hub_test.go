package zetcd

import (
	clientv3 "go.etcd.io/etcd/client/v3"
	"reflect"
	"testing"
)

func TestHub_AddEtcdServerAddress(t *testing.T) {
	type fields struct {
		Hosts         []string
		DialTimeout   int64
		DialKeepAlive int64
		client        *clientv3.Client
	}
	type args struct {
		addr []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "AddEtcdServerAddress",
			fields: fields{
				Hosts:         []string{"http://etcd-server:2379"},
				DialTimeout:   1,
				DialKeepAlive: 1,
				client:        nil,
			},
			args: args{addr: []string{"http://etcd-server:2379"}},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := &Hub{
				Etcd: &Etcd{
					Hosts:         tt.fields.Hosts,
					DialTimeout:   tt.fields.DialTimeout,
					DialKeepalive: tt.fields.DialKeepAlive,
				},
				client:        tt.fields.client,
			}
			z.connect()
			z.AddEtcdServerAddress([]string{"127.0.0.1:2379"})
			wat := []string{"http://etcd-server:2379","127.0.0.1:2379"}
			if !reflect.DeepEqual(z.Etcd.Hosts, wat){
				t.Errorf("NewHub error want: %+v, got: %+v",wat,z.Etcd.Hosts)
			}
			got,err := NewHub(z.Etcd)
			got.client = z.client
			if err != nil ||  !reflect.DeepEqual(z, got) {
				t.Errorf("NewHub error want: %+v, got: %+v, err:%v",z,got,err)
			}
		})
	}
}

func TestHub_GetAll(t *testing.T) {
	type fields struct {
		Hosts         []string
		DialTimeout   int64
		DialKeepAlive int64
		client        *clientv3.Client
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *[]string
		wantErr bool
	}{
		{
			name:"GetAll",
			fields: fields{
				Hosts:         nil,
				DialTimeout:   0,
				DialKeepAlive: 0,
				client:        nil,
			},
			args: args{key: "prod/book:grpc"},
			want: &[]string{},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := &Hub{
				Etcd: &Etcd{
					Hosts:         tt.fields.Hosts,
					DialTimeout:   tt.fields.DialTimeout,
					DialKeepalive: tt.fields.DialKeepAlive,
				},
				client:        tt.fields.client,
			}
			got, err := z.GetAll(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAll() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAll() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHub_GetClient(t *testing.T) {
	type fields struct {
		Hosts         []string
		DialTimeout   int64
		DialKeepAlive int64
		client        *clientv3.Client
	}
	tests := []struct {
		name   string
		fields fields
		want   *clientv3.Client
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := &Hub{
				Etcd: &Etcd{
					Hosts:         tt.fields.Hosts,
					DialTimeout:   tt.fields.DialTimeout,
					DialKeepalive: tt.fields.DialKeepAlive,
				},
				client:        tt.fields.client,
			}
			if got := z.GetClient(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetClient() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHub_GetOne(t *testing.T) {
	type fields struct {
		Hosts         []string
		DialTimeout   int64
		DialKeepAlive int64
		client        *clientv3.Client
	}
	type args struct {
		key     string
		balance Balancer
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := &Hub{
				Etcd: &Etcd{
					Hosts:         tt.fields.Hosts,
					DialTimeout:   tt.fields.DialTimeout,
					DialKeepalive: tt.fields.DialKeepAlive,
				},
				client:        tt.fields.client,
			}
			got, err := z.GetOne(tt.args.key, tt.args.balance)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetOne() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetOne() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHub_SetDialKeepAlive(t *testing.T) {
	type fields struct {
		Hosts         []string
		DialTimeout   int64
		DialKeepAlive int64
		client        *clientv3.Client
	}
	type args struct {
		dialKeepAlive int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Hub
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := &Hub{
				Etcd: &Etcd{
					Hosts:         tt.fields.Hosts,
					DialTimeout:   tt.fields.DialTimeout,
					DialKeepalive: tt.fields.DialKeepAlive,
				},
				client:        tt.fields.client,
			}
			if got := z.SetDialKeepAlive(tt.args.dialKeepAlive); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetDialKeepAlive() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHub_SetDialTimeout(t *testing.T) {
	type fields struct {
		Hosts         []string
		DialTimeout   int64
		DialKeepAlive int64
		client        *clientv3.Client
	}
	type args struct {
		dialTimeout int64
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   *Hub
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := &Hub{
				Etcd: &Etcd{
					Hosts:         tt.fields.Hosts,
					DialTimeout:   tt.fields.DialTimeout,
					DialKeepalive: tt.fields.DialKeepAlive,
				},
				client:        tt.fields.client,
			}
			if got := z.SetDialTimeout(tt.args.dialTimeout); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("SetDialTimeout() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHub_SetEtcdServerAddress(t *testing.T) {
	type fields struct {
		Hosts         []string
		DialTimeout   int64
		DialKeepAlive int64
		client        *clientv3.Client
	}
	type args struct {
		etcdHosts []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		// TODO: Add test cases.
		{
			name: "SetEtcdServerAddress",
			fields: fields{
				Hosts:         nil,
				DialTimeout:   0,
				DialKeepAlive: 0,
				client:        nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := &Hub{
				Etcd: &Etcd{
					Hosts:         tt.fields.Hosts,
					DialTimeout:   tt.fields.DialTimeout,
					DialKeepalive: tt.fields.DialKeepAlive,
				},
				client:        tt.fields.client,
			}
			z.SetEtcdServerAddress([]string{"127.0.0.1:2379"})
			if !reflect.DeepEqual(z.Etcd.Hosts,[]string{"127.0.0.1:2379"}){
				t.Errorf("want=%v, got=%v",[]string{"127.0.0.1:2379"},z.Etcd.Hosts)
			}
		})
	}
}

func TestHub_connect(t *testing.T) {
	type fields struct {
		Hosts         []string
		DialTimeout   int64
		DialKeepAlive int64
		client        *clientv3.Client
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := &Hub{
				Etcd: &Etcd{
					Hosts:         tt.fields.Hosts,
					DialTimeout:   tt.fields.DialTimeout,
					DialKeepalive: tt.fields.DialKeepAlive,
				},
				client:        tt.fields.client,
			}
			if err := z.connect(); (err != nil) != tt.wantErr {
				t.Errorf("connect() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHub_grant(t *testing.T) {
	type fields struct {
		Hosts         []string
		DialTimeout   int64
		DialKeepAlive int64
		client        *clientv3.Client
	}
	tests := []struct {
		name    string
		fields  fields
		want    *clientv3.LeaseGrantResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := &Hub{
				Etcd: &Etcd{
					Hosts:         tt.fields.Hosts,
					DialTimeout:   tt.fields.DialTimeout,
					DialKeepalive: tt.fields.DialKeepAlive,
				},
				client:        tt.fields.client,
			}
			got, err := z.grant()
			if (err != nil) != tt.wantErr {
				t.Errorf("grant() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("grant() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHub_keepAlive(t *testing.T) {
	type fields struct {
		Hosts         []string
		DialTimeout   int64
		DialKeepAlive int64
		client        *clientv3.Client
	}
	type args struct {
		id clientv3.LeaseID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    <-chan *clientv3.LeaseKeepAliveResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := &Hub{
				Etcd: &Etcd{
					Hosts:         tt.fields.Hosts,
					DialTimeout:   tt.fields.DialTimeout,
					DialKeepalive: tt.fields.DialKeepAlive,
				},
				client:        tt.fields.client,
			}
			got, err := z.keepAlive(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("keepAlive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("keepAlive() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHub_put(t *testing.T) {
	type fields struct {
		Hosts         []string
		DialTimeout   int64
		DialKeepAlive int64
		client        *clientv3.Client
	}
	type args struct {
		key   string
		value string
		id    clientv3.LeaseID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "put",
			fields: fields{
				Hosts:         []string{"http://etcd-server:2379"},
				DialTimeout:   1,
				DialKeepAlive: 1,
				client:        nil,
			},
			args: args{
				key: "prod/book:grpc",
				value: "192.168.2.1",
				id: 33443,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := &Hub{
				Etcd: &Etcd{
					Hosts:         tt.fields.Hosts,
					DialTimeout:   tt.fields.DialTimeout,
					DialKeepalive: tt.fields.DialKeepAlive,
				},
				client:        tt.fields.client,
			}
			if err := z.put(tt.args.key, tt.args.value, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("put() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHub_revoke(t *testing.T) {
	type fields struct {
		Hosts         []string
		DialTimeout   int64
		DialKeepAlive int64
		client        *clientv3.Client
	}
	type args struct {
		id clientv3.LeaseID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "revoke",
			fields: fields{
				Hosts:         []string{"http://etcd-server:2379"},
				DialTimeout:   1,
				DialKeepAlive: 1,
				client:        nil,
			},
			args: args{id: 23333},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := &Hub{
				Etcd: &Etcd{
					Hosts:         tt.fields.Hosts,
					DialTimeout:   tt.fields.DialTimeout,
					DialKeepalive: tt.fields.DialKeepAlive,
				},
				client:        tt.fields.client,
			}
			if err := z.revoke(tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("revoke() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestHub_timeToLive(t *testing.T) {
	type fields struct {
		Hosts         []string
		DialTimeout   int64
		DialKeepAlive int64
		client        *clientv3.Client
	}
	type args struct {
		id clientv3.LeaseID
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *clientv3.LeaseTimeToLiveResponse
		wantErr bool
	}{
		{
			name: "timeToLive",
			fields: fields{
				Hosts:         []string{"http://etcd-server:2379"},
				DialTimeout:   1,
				DialKeepAlive: 1,
				client:        nil,
			},
			args: args{id: 123222},
			want: nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := &Hub{
				Etcd: &Etcd{
					Hosts:         tt.fields.Hosts,
					DialTimeout:   tt.fields.DialTimeout,
					DialKeepalive: tt.fields.DialKeepAlive,
				},
				client:        tt.fields.client,
			}
			got, err := z.timeToLive(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("timeToLive() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("timeToLive() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestHub_watch(t *testing.T) {
	type fields struct {
		Hosts         []string
		DialTimeout   int64
		DialKeepAlive int64
		client        *clientv3.Client
	}
	type args struct {
		key string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
	}{
		{
			name: "watch",
			fields: fields{
				Hosts:         []string{"http://etcd-server:2379"},
				DialTimeout:   1,
				DialKeepAlive: 1,
				client:        nil,
			},
			args: args{key: "prod/book:grpc"},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			z := &Hub{
				Etcd: &Etcd{
					Hosts:         tt.fields.Hosts,
					DialTimeout:   tt.fields.DialTimeout,
					DialKeepalive: tt.fields.DialKeepAlive,
				},
				client:        tt.fields.client,
			}
			z.watch("prod/book:grpc")
			t.Parallel()
		})
	}
}

func TestNewHub(t *testing.T) {
	type args struct {
		hosts         []string
		dialTimeout   int64
		dialKeepAlive int64
	}
	etcd := &Etcd{
		Hosts:         []string{"127.0.0.1:2379"},
		DialTimeout:   1,
		DialKeepalive: 1,
	}
	hub,_ := NewHub(etcd)
	tests := []struct {
		name    string
		args    args
		want    *Hub
		wantErr bool
	}{
		{
			name: "NewHub",
			args: args{
				hosts:         []string{"127.0.0.1:2379"},
				dialTimeout:   1,
				dialKeepAlive: 1,
			},
			want: hub,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			etcd := &Etcd{
				tt.args.hosts, tt.args.dialTimeout, tt.args.dialKeepAlive,
			}
			got, err := NewHub(etcd)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewHub() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewHub() got = %v, want %v", got, tt.want)
			}
		})
	}
}
