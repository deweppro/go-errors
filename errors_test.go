package errors

import (
	e "errors"
	"testing"
)

func TestNew(t *testing.T) {
	type args struct {
		message string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{name: "Case1", args: args{message: "hello"}, want: "hello", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := New(tt.args.message)
			if (err != nil) != tt.wantErr {
				t.Errorf("New() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if err.Error() != tt.want {
				t.Errorf("New() error = %v, want %v", err.Error(), tt.want)
				return
			}
		})
	}
}

func TestWrap(t *testing.T) {
	type args struct {
		msg []error
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name:    "Case1",
			args:    args{msg: nil},
			want:    "",
			wantErr: false,
		},
		{
			name:    "Case2",
			args:    args{msg: []error{New("hello"), e.New("world")}},
			want:    "hello: world",
			wantErr: true,
		},
		{
			name:    "Case3",
			args:    args{msg: []error{New("err1"), e.New("err2"), nil, e.New("err3")}},
			want:    "err1: err2: err3",
			wantErr: true,
		},
		{
			name:    "Case4",
			args:    args{msg: []error{WrapMessage(New("err1"), "err1 message"), WrapMessage(e.New("err2"), "err2 message"), WrapMessage(e.New("err3"), "err3 message")}},
			want:    "err1 message: err1\n\t[trace] github.com/deweppro/go-errors.TestWrap:65\n: err2 message: err2\n\t[trace] github.com/deweppro/go-errors.TestWrap:65\n: err3 message: err3\n\t[trace] github.com/deweppro/go-errors.TestWrap:65\n",
			wantErr: true,
		},
		{
			name:    "Case5",
			args:    args{msg: []error{nil, nil, nil}},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := Wrap(tt.args.msg...)
			if (err != nil) != tt.wantErr {
				t.Errorf("Wrap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.want {
				t.Errorf("Wrap() error = %v, want %v", err.Error(), tt.want)
				return
			}
		})
	}
}

func TestWrapMessage(t *testing.T) {
	type args struct {
		cause   error
		message string
		args    []interface{}
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Case1",
			args: args{
				cause:   nil,
				message: "err context",
				args:    nil,
			},
			want:    "",
			wantErr: false,
		},
		{
			name: "Case2",
			args: args{
				cause:   e.New("err1"),
				message: "err context",
				args:    nil,
			},
			want:    "err context: err1\n\t[trace] github.com/deweppro/go-errors.TestWrapMessage.func1:136\n",
			wantErr: true,
		},
		{
			name: "Case3",
			args: args{
				cause:   e.New("err1"),
				message: "bad ip %s",
				args:    []interface{}{"127.0.0.1"},
			},
			want:    "bad ip 127.0.0.1: err1\n\t[trace] github.com/deweppro/go-errors.TestWrapMessage.func1:136\n",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := WrapMessage(tt.args.cause, tt.args.message, tt.args.args...)
			if (err != nil) != tt.wantErr {
				t.Errorf("WrapMessage() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.want {
				t.Errorf("WrapMessage() error = %v, want %v", err.Error(), tt.want)
				return
			}
		})
	}
}

func Test_errMessage_CauseUnwrap(t *testing.T) {
	type fields struct {
		cause   error
		message string
	}
	tests := []struct {
		name    string
		fields  fields
		want    string
		wantErr bool
	}{
		{
			name: "Case1",
			fields: fields{
				cause:   e.New("err1"),
				message: "context",
			},
			want:    "err1",
			wantErr: true,
		},
		{
			name: "Case2",
			fields: fields{
				cause:   nil,
				message: "context",
			},
			want:    "err1",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &errMessage{
				cause:   tt.fields.cause,
				message: tt.fields.message,
			}
			err := v.Cause()
			if (err != nil) != tt.wantErr {
				t.Errorf("Cause() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.want {
				t.Errorf("Cause() error = %v, want %v", err.Error(), tt.want)
				return
			}
			err = v.Unwrap()
			if (err != nil) != tt.wantErr {
				t.Errorf("Unwrap() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if tt.wantErr && err.Error() != tt.want {
				t.Errorf("Unwrap() error = %v, want %v", err.Error(), tt.want)
				return
			}
		})
	}
}

func Test_Is(t *testing.T) {
	err0 := New("test")
	type args struct {
		err    error
		target error
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "Case1", args: args{err: err0, target: err0}, want: true},
		{name: "Case2", args: args{err: WrapMessage(err0, "ttt"), target: err0}, want: true},
		{name: "Case3", args: args{err: New("hello"), target: err0}, want: false},
		{name: "Case4", args: args{err: nil, target: err0}, want: false},
		{name: "Case5", args: args{err: New("hello"), target: nil}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Is(tt.args.err, tt.args.target); got != tt.want {
				t.Errorf("Is() = %v, want %v", got, tt.want)
			}
		})
	}
}
