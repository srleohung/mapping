package structure

import (
	"log"
	"reflect"
	"testing"
)

func TestGetType(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		args args
		want reflect.Kind
	}{
		{
			args: args{i: args{}},
			want: reflect.Struct,
		},
		{
			args: args{i: &args{}},
			want: reflect.Struct,
		},
		{
			args: args{i: ""},
			want: reflect.String,
		},
		{
			args: args{i: 0},
			want: reflect.Int,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := GetType(tt.args.i); got.Kind() != tt.want {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetTypeName(t *testing.T) {
	type args struct {
		i interface{}
	}
	tests := []struct {
		args args
		want string
	}{
		{
			args: args{i: args{}},
			want: "args",
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if got := GetTypeName(tt.args.i); got != tt.want {
				t.Errorf("got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSetFieldValue(t *testing.T) {
	type Args struct {
		ID int
	}
	type args struct {
		i Args
		f string
		n interface{}
	}
	tests := []struct {
		args    args
		wantErr bool
		want    Args
	}{
		{
			args: args{
				i: Args{ID: 0},
				f: "ID",
				n: 1,
			},
			wantErr: false,
			want:    Args{ID: 1},
		},
		{
			args: args{
				i: Args{ID: 0},
				f: "ID",
				n: "1",
			},
			wantErr: false,
			want:    Args{ID: 1},
		},
		{
			args: args{
				i: Args{ID: 0},
				f: "ID",
				n: 1.1,
			},
			wantErr: false,
			want:    Args{ID: 1},
		},
		{
			args: args{
				i: Args{ID: 0},
				f: "ID",
				n: true,
			},
			wantErr: false,
			want:    Args{ID: 1},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			log.Printf("Before: %+v", tt.args.i)
			if err := SetFieldValue(&tt.args.i, tt.args.f, tt.args.n); (err != nil) != tt.wantErr {
				t.Errorf("err = %v, wantErr %v", err, tt.wantErr)
			}
			if tt.args.i.ID != tt.want.ID {
				t.Errorf("i = %v, want %v", tt.args.i, tt.want)
			}
			log.Printf("After: %+v", tt.args.i)
		})
	}
}
