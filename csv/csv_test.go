package csv

import (
	"io"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestNewCSV(t *testing.T) {
	os.Setenv(headersEnvKey, "header1,header2")
	defer os.Unsetenv(headersEnvKey)

	tests := []struct {
		name    string
		file    io.Reader
		want    []Record
		wantErr bool
	}{
		{
			name:    "csv file input is empty",
			file:    strings.NewReader(""),
			want:    []Record{},
			wantErr: true,
		},
		{
			name:    "csv file input has only header row",
			file:    strings.NewReader("header1,header2\n"),
			want:    []Record{},
			wantErr: true,
		},
		{
			name:    "csv file input has different format",
			file:    strings.NewReader("header1,header2,header3\n"),
			want:    []Record{},
			wantErr: true,
		},
		{
			name: "happy path",
			file: strings.NewReader("header1,header2\ndata1,data2\n"),
			want: []Record{{
				Data: map[string]string{
					"header1": "data1",
					"header2": "data2",
				},
			}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := NewCSV(tt.file)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewCSV() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewCSV() = %v, want %v", got, tt.want)
			}
		})
	}
}
