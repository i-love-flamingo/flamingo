package infrastructure

import (
	"path/filepath"
	"reflect"
	"testing"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

func TestRedirectData_Get(t *testing.T) {
	type fields struct {
		filePath string
		logger   flamingo.Logger
	}
	tests := []struct {
		name   string
		fields fields
		want   []CsvContent
	}{
		{
			name: "with file",
			fields: fields{
				filePath: filepath.Join("testdata", "redirect.csv"),
				logger:   flamingo.NullLogger{},
			},
			want: []CsvContent{
				{
					HTTPStatusCode: 301,
					OriginalPath:   "foo",
					RedirectTarget: "bar",
				},
				{
					HTTPStatusCode: 302,
					OriginalPath:   "baz",
					RedirectTarget: "bam",
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rd := NewRedirectData(
				&struct {
					RedirectsCsv string `inject:"config:redirects.csv"`
				}{RedirectsCsv: tt.fields.filePath},
				tt.fields.logger,
			)
			if got := rd.Get(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RedirectData.Get() = %v, want %v", got, tt.want)
			}
		})
	}
}
