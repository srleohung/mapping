package structure

import (
	"fmt"
	"log"
	"reflect"
	"testing"
	"time"
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
				t.Errorf("got = %v, want = %v", got, tt.want)
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
				t.Errorf("got = %v, want = %v", got, tt.want)
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
				t.Errorf("err = %v, wantErr = %v", err, tt.wantErr)
			}
			if tt.args.i.ID != tt.want.ID {
				t.Errorf("i = %v, want = %v", tt.args.i, tt.want)
			}
			log.Printf("After: %+v", tt.args.i)
		})
	}
}

func TestStructToMap(t *testing.T) {
	type Args struct {
		ID   int
		IDs  []int
		Time time.Time
	}
	type args struct {
		i Args
	}
	tests := []struct {
		args args
		want map[string]interface{}
	}{
		{
			args: args{
				i: Args{ID: 0, IDs: []int{1, 2, 3}, Time: time.Unix(0, 0)},
			},
			want: map[string]interface{}{"ID": 0, "IDs": []int{1, 2, 3}, "Time": time.Unix(0, 0)},
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			got := StructToMap(&tt.args.i)
			for k, v := range got {
				if fmt.Sprintf("%v", v) != fmt.Sprintf("%v", tt.want[k]) {
					t.Errorf("\n%v got = %v, \nwant = %v", k, v, tt.want[k])
				}
			}
		})
	}
}

func TestStructToStruct(t *testing.T) {
	type Z struct {
		ID   int
		Name string
		Time time.Time
	}
	type Y struct {
		ID   int
		Name string
		Time time.Time
	}
	type X struct {
		ID   int
		Y    Y
		Zs   []Z
		Time time.Time
	}
	type C struct {
		ID   int       `struct:"ID"`
		Name string    `struct:"Name"`
		Time time.Time `struct:"Time"`
	}
	type B struct {
		ID   int       `struct:"Y.ID"`
		Name string    `struct:"Y.Name"`
		Time time.Time `struct:"Y.Time"`
	}
	type A struct {
		ID    string    `struct:"ID"`
		Time  time.Time `struct:"Time"`
		YID   string    `struct:"Y.ID"`
		YName string    `struct:"Y.Name"`
		YTime time.Time `struct:"Y.Time"`
		B     B         `struct:"Y"`
		Cs    []C       `struct:"Zs"`
	}
	type args struct {
		i X
		j A
	}
	tests := []struct {
		args    args
		wantErr bool
	}{
		{
			args: args{
				i: X{
					ID:   9,
					Y:    Y{ID: 8, Name: "Y", Time: time.Now().Add(-1 * 24 * time.Hour)},
					Zs:   []Z{{ID: 7, Name: "Z1", Time: time.Now().Add(-2 * 24 * time.Hour)}, {ID: 6, Name: "Z2", Time: time.Now().Add(-3 * 24 * time.Hour)}, {ID: 5, Name: "Z3", Time: time.Now().Add(-4 * 24 * time.Hour)}},
					Time: time.Now().Add(-5 * 24 * time.Hour),
				},
				j: A{},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run("", func(t *testing.T) {
			if err := StructToStruct(tt.args.i, &tt.args.j); (err != nil) != tt.wantErr {
				t.Errorf("\ngot = %v, \nwant = %v", err, tt.wantErr)
			}
			if tt.args.j.ID != fmt.Sprintf("%v", tt.args.i.ID) {
				t.Errorf("\nID got = %v, \nwant = %v", tt.args.j.ID, tt.args.i.ID)
			}
			if tt.args.j.Time != tt.args.i.Time {
				t.Errorf("\nTime got = %v, \nwant = %v", tt.args.j.Time, tt.args.i.Time)
			}
			if tt.args.j.YID != fmt.Sprintf("%v", tt.args.i.Y.ID) {
				t.Errorf("\nYID got = %v, \nwant = %v", tt.args.j.YID, tt.args.i.Y.ID)
			}
			if tt.args.j.YName != tt.args.i.Y.Name {
				t.Errorf("\nYName got = %v, \nwant = %v", tt.args.j.YName, tt.args.i.Y.Name)
			}
			if tt.args.j.YTime != tt.args.i.Y.Time {
				t.Errorf("\nYTime got = %v, \nwant = %v", tt.args.j.YTime, tt.args.i.Y.Time)
			}
			if tt.args.j.B.ID != tt.args.i.Y.ID {
				t.Errorf("\nB.ID got = %v, \nwant = %v", tt.args.j.B.ID, tt.args.i.Y.ID)
			}
			if tt.args.j.B.Name != tt.args.i.Y.Name {
				t.Errorf("\nB.Name got = %v, \nwant = %v", tt.args.j.B.Name, tt.args.i.Y.Name)
			}
			if tt.args.j.B.Time != tt.args.i.Y.Time {
				t.Errorf("\nB.Time got = %v, \nwant = %v", tt.args.j.B.Time, tt.args.i.Y.Time)
			}
			if len(tt.args.j.Cs) != len(tt.args.i.Zs) {
				t.Errorf("\nlen(Cs) got = %v, \nwantLen = %v", len(tt.args.j.Cs), len(tt.args.i.Zs))
			}
			for i, c := range tt.args.j.Cs {
				if c.ID != tt.args.i.Zs[i].ID {
					t.Errorf("\nCs.%v.ID got = %v, \nwant = %v", i, c.ID, tt.args.i.Zs[i].ID)
				}
				if c.Name != tt.args.i.Zs[i].Name {
					t.Errorf("\nCs.%v.ID got = %v, \nwant = %v", i, c.Name, tt.args.i.Zs[i].Name)
				}
				if c.Time != tt.args.i.Zs[i].Time {
					t.Errorf("\nCs.%v.ID got = %v, \nwant = %v", i, c.Time, tt.args.i.Zs[i].Time)
				}
			}
		})
	}
}
