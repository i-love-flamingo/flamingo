package gotemplate

import (
	"bytes"
	"context"
	"io"
	"reflect"
	"testing"
	"time"

	"flamingo.me/flamingo/v3/framework/flamingo"
)

var noAdditionalTemplateFuncs = func() map[string]flamingo.TemplateFunc {
	return make(map[string]flamingo.TemplateFunc)
}

var additionalTemplateFuncs = func() map[string]flamingo.TemplateFunc {
	funcs := make(map[string]flamingo.TemplateFunc)
	funcs["customTemplateFunc"] = CustomTemplateFunc{}
	return funcs
}

type CustomTemplateFunc struct{}

var _ flamingo.TemplateFunc = CustomTemplateFunc{}

func (CustomTemplateFunc) Func(ctx context.Context) interface{} {
	return func() interface{} {
		return "test-abc"
	}
}

func Test_engine_Render(t *testing.T) {
	type fields struct {
		templatesBasePath  string
		layoutTemplatesDir string
		debug              bool
		tplFuncs           func() map[string]flamingo.TemplateFunc
	}
	type args struct {
		name string
		data interface{}
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    io.Reader
		wantErr bool
	}{
		{
			name: "Template base path not found",
			fields: fields{
				templatesBasePath: "non-existing-dir/",
				tplFuncs:          noAdditionalTemplateFuncs,
			},
			wantErr: true,
		},
		{
			name: "Layout path not found",
			fields: fields{
				templatesBasePath:  "testdata/test-simple",
				layoutTemplatesDir: "non-existing-layout-dir",
				tplFuncs:           noAdditionalTemplateFuncs,
			},
			wantErr: true,
		},
		{
			name: "Template not found",
			fields: fields{
				templatesBasePath:  "testdata/test-simple",
				layoutTemplatesDir: "",
				debug:              false,
				tplFuncs:           noAdditionalTemplateFuncs,
			},
			args: args{
				name: "non-existing-template",
				data: nil,
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "Simple template found and working",
			fields: fields{
				templatesBasePath:  "testdata/test-simple",
				layoutTemplatesDir: "",
				debug:              false,
				tplFuncs:           noAdditionalTemplateFuncs,
			},
			args: args{
				name: "simple",
				data: nil,
			},
			want:    bytes.NewBuffer([]byte("Hello World!")),
			wantErr: false,
		},
		{
			name: "Built in template functions work",
			fields: fields{
				templatesBasePath:  "testdata/test-simple",
				layoutTemplatesDir: "",
				debug:              false,
				tplFuncs:           noAdditionalTemplateFuncs,
			},
			args: args{
				name: "built-in-template-funcs",
				data: struct {
					Time time.Time
				}{time.Unix(0, 0)},
			},
			want: bytes.NewBuffer([]byte(
				`Upper: HELLO WORLD!
formatDate: 1970-01-01
map (invalid params): map[]
map (valid params): map[a:b x:y]`)),
			wantErr: false,
		},
		{
			name: "Additional template functions work",
			fields: fields{
				templatesBasePath:  "testdata/test-simple",
				layoutTemplatesDir: "",
				debug:              false,
				tplFuncs:           additionalTemplateFuncs,
			},
			args: args{
				name: "additional-template-funcs",
				data: struct {
					Time time.Time
				}{time.Unix(0, 0)},
			},
			want: bytes.NewBuffer([]byte(
				`Upper: HELLO WORLD!
formatDate: 1970-01-01
map (invalid params): map[]
map (valid params): map[a:b x:y]
customTemplateFunc: test-abc`)),
			wantErr: false,
		},
		{
			name: "Nested layouts/templates should work",
			fields: fields{
				templatesBasePath:  "testdata/test-nested-dirs",
				layoutTemplatesDir: "layouts",
				debug:              false,
				tplFuncs:           additionalTemplateFuncs,
			},
			args: args{
				name: "dir-a/sub-dir-a/main",
				data: nil,
			},
			want: bytes.NewBuffer([]byte(
				`<!DOCTYPE html><html lang="en"><head><meta charset="UTF-8"><title>Title</title></head><body><div>Hello World!</div></body></html>

`)),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &engine{}
			e.Inject(tt.fields.tplFuncs, flamingo.NullLogger{}, &struct {
				TemplatesBasePath  string `inject:"config:gotemplates.engine.templates.basepath"`
				LayoutTemplatesDir string `inject:"config:gotemplates.engine.layout.dir"`
				Debug              bool   `inject:"config:debug.mode"`
			}{
				tt.fields.templatesBasePath,
				tt.fields.layoutTemplatesDir,
				tt.fields.debug,
			})
			got, err := e.Render(context.Background(), tt.args.name, tt.args.data)
			if (err != nil) != tt.wantErr {
				t.Errorf("Render() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Render() got = %v, want %v", got, tt.want)
			}
		})
	}
}
