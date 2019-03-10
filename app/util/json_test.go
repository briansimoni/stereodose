package util

import (
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

var wantString = `[
"testarray"
]`

func TestJSON(t *testing.T) {
	type args struct {
		data interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantW   string
		wantErr bool
	}{
		{
			name: "Valid data",
			args: args{
				data: []string{"testarray"},
			},
			wantW:   wantString,
			wantErr: false,
		},
		{
			name: "Invalid data",
			args: args{
				data: make(chan int),
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			if err := JSON(w, tt.args.data); (err != nil) != tt.wantErr {
				t.Errorf("JSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			result := w.Result()
			defer result.Body.Close()
			data, err := ioutil.ReadAll(result.Body)
			if err != nil {
				t.Error("Reading JSON result failed")
			}

			gotW := string(data)
			// I was having trouble getting the test case to actually have the right indentation
			gotWnoWhiteSpace := strings.Replace(gotW, "\n", "", -1)
			gotWnoWhiteSpace = strings.Replace(gotW, "	", "", -1)
			if gotWnoWhiteSpace != tt.wantW {
				t.Errorf("JSON() = \n%v\n, want \n%v\n", gotWnoWhiteSpace, tt.wantW)
			}
		})
	}
}
