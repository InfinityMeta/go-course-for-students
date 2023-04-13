package homework

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestValidate(t *testing.T) {
	type args struct {
		v any
	}
	tests := []struct {
		name     string
		args     args
		wantErr  bool
		checkErr func(err error) bool
	}{
		{
			name: "invalid struct: interface",
			args: args{
				v: new(any),
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "invalid struct: map",
			args: args{
				v: map[string]string{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "invalid struct: string",
			args: args{
				v: "some string",
			},
			wantErr: true,
			checkErr: func(err error) bool {
				return errors.Is(err, ErrNotStruct)
			},
		},
		{
			name: "valid struct with no fields",
			args: args{
				v: struct{}{},
			},
			wantErr: false,
		},
		{
			name: "valid struct with untagged fields",
			args: args{
				v: struct {
					f1 string
					f2 string
				}{},
			},
			wantErr: false,
		},
		{
			name: "valid struct with unexported fields",
			args: args{
				v: struct {
					foo string `validate:"len:10"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				e := &ValidationErrors{}
				return errors.As(err, e) && e.Error() == ErrValidateForUnexportedFields.Error()
			},
		},
		{
			name: "invalid validator syntax",
			args: args{
				v: struct {
					Foo string `validate:"len:abcdef"`
				}{},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				e := &ValidationErrors{}
				return errors.As(err, e) && e.Error() == ErrInvalidValidatorSyntax.Error()
			},
		},
		{
			name: "valid struct with tagged fields",
			args: args{
				v: struct {
					Len       string `validate:"len:20"`
					LenZ      string `validate:"len:0"`
					InInt     int    `validate:"in:20,25,30"`
					InNeg     int    `validate:"in:-20,-25,-30"`
					InStr     string `validate:"in:foo,bar"`
					MinInt    int    `validate:"min:10"`
					MinIntNeg int    `validate:"min:-10"`
					MinStr    string `validate:"min:10"`
					MinStrNeg string `validate:"min:-1"`
					MaxInt    int    `validate:"max:20"`
					MaxIntNeg int    `validate:"max:-2"`
					MaxStr    string `validate:"max:20"`
				}{
					Len:       "abcdefghjklmopqrstvu",
					LenZ:      "",
					InInt:     25,
					InNeg:     -25,
					InStr:     "bar",
					MinInt:    15,
					MinIntNeg: -9,
					MinStr:    "abcdefghjkl",
					MinStrNeg: "abc",
					MaxInt:    16,
					MaxIntNeg: -3,
					MaxStr:    "abcdefghjklmopqrst",
				},
			},
			wantErr: false,
		},
		{
			name: "wrong length",
			args: args{
				v: struct {
					Lower    string `validate:"len:24"`
					Higher   string `validate:"len:5"`
					Zero     string `validate:"len:3"`
					BadSpec  string `validate:"len:%12"`
					Negative string `validate:"len:-6"`
				}{
					Lower:    "abcdef",
					Higher:   "abcdef",
					Zero:     "",
					BadSpec:  "abc",
					Negative: "abcd",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong in",
			args: args{
				v: struct {
					InA     string `validate:"in:ab,cd"`
					InB     string `validate:"in:aa,bb,cd,ee"`
					InC     int    `validate:"in:-1,-3,5,7"`
					InD     int    `validate:"in:5-"`
					InEmpty string `validate:"in:"`
				}{
					InA:     "ef",
					InB:     "ab",
					InC:     2,
					InD:     12,
					InEmpty: "",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong min",
			args: args{
				v: struct {
					MinA string `validate:"min:12"`
					MinB int    `validate:"min:-12"`
					MinC int    `validate:"min:5-"`
					MinD int    `validate:"min:"`
					MinE string `validate:"min:"`
				}{
					MinA: "ef",
					MinB: -22,
					MinC: 12,
					MinD: 11,
					MinE: "abc",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong max",
			args: args{
				v: struct {
					MaxA string `validate:"max:2"`
					MaxB string `validate:"max:-7"`
					MaxC int    `validate:"max:-12"`
					MaxD int    `validate:"max:5-"`
					MaxE int    `validate:"max:"`
					MaxF string `validate:"max:"`
				}{
					MaxA: "efgh",
					MaxB: "ab",
					MaxC: 22,
					MaxD: 12,
					MaxE: 11,
					MaxF: "abc",
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 6)
				return true
			},
		},
		{
			name: "valid struct with slices",
			args: args{
				v: struct {
					MinInt []int    `validate:"min:1"`
					MaxInt []int    `validate:"max:10"`
					MinStr []string `validate:"min:3"`
					MaxStr []string `validate:"max:5"`
					InInt  []int    `validate:"in:1,2,3,4,5,6,7,8"`
					InStr  []string `validate:"in:ab,cd,ef,gn,hm"`
					Len    []string `validate:"len:4"`
				}{
					MinInt: []int{1, 2, 3, 4, 5},
					MaxInt: []int{6, 7, 8, 9, 10},
					MinStr: []string{"abc", "cdef", "words"},
					MaxStr: []string{"a", "cde", "weqwq"},
					InInt:  []int{2, 4, 6},
					InStr:  []string{"ab", "cd", "hm"},
					Len:    []string{"abcd", "efdg", "gfdg"},
				},
			},
			wantErr: false,
		},
		{
			name: "wrong max slice",
			args: args{
				v: struct {
					MaxSliceA []int    `validate:"max:10"`
					MaxSliceB []int    `validate:"max:7-"`
					MaxSliceC []int    `validate:"max:"`
					MaxSliceD []string `validate:"max:5"`
					MaxSliceE []string `validate:"max:-2"`
					MaxSliceF []string `validate:"max:"`
				}{
					MaxSliceA: []int{5, 10, 15},
					MaxSliceB: []int{1, 2, 3},
					MaxSliceC: []int{4, 8},
					MaxSliceD: []string{"abc", "dfehtj"},
					MaxSliceE: []string{"rew", "tyu"},
					MaxSliceF: []string{"as", "gd", "ty"},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 6)
				return true
			},
		},
		{
			name: "wrong min slice",
			args: args{
				v: struct {
					MinSliceA []int    `validate:"min:5"`
					MinSliceB []int    `validate:"min:7-"`
					MinSliceC []int    `validate:"min:"`
					MinSliceD []string `validate:"min:3"`
					MinSliceE []string `validate:"min:"`
				}{
					MinSliceA: []int{1, 10, 15},
					MinSliceB: []int{1, 2, 3},
					MinSliceC: []int{4, 8},
					MinSliceD: []string{"ab", "dfehtj"},
					MinSliceE: []string{"rew", "tyu"},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong length slice",
			args: args{
				v: struct {
					Lower    []string `validate:"len:5"`
					Higher   []string `validate:"len:5"`
					Zero     []string `validate:"len:3"`
					BadSpec  []string `validate:"len:%12"`
					Negative []string `validate:"len:-6"`
				}{
					Lower:    []string{"abc", "gfdsd", "dadad"},
					Higher:   []string{"asdasdfaf", "fsdsd", "weree"},
					Zero:     []string{"", ""},
					BadSpec:  []string{"asf", "sdg", "wqe"},
					Negative: []string{"fd", "re"},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
		{
			name: "wrong in slice",
			args: args{
				v: struct {
					InSliceA     []string `validate:"in:ab,cd"`
					InSliceB     []string `validate:"in:aa,bb,cd,ee"`
					InSliceC     []int    `validate:"in:-1,-3,5,7"`
					InSliceD     []int    `validate:"in:5-"`
					InSliceEmpty []string `validate:"in:"`
				}{
					InSliceA:     []string{"ab", "ef"},
					InSliceB:     []string{"ab", "bc"},
					InSliceC:     []int{2, 4},
					InSliceD:     []int{3, 5, 7},
					InSliceEmpty: []string{""},
				},
			},
			wantErr: true,
			checkErr: func(err error) bool {
				assert.Len(t, err.(ValidationErrors), 5)
				return true
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Validate(tt.args.v)
			if tt.wantErr {
				assert.Error(t, err)
				assert.True(t, tt.checkErr(err), "test expect an error, but got wrong error type")
			} else {
				assert.NoError(t, err)
			}
		})
	}

}
