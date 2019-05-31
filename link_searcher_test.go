package main

import (
	"net/url"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestAbsoluteUrl(t *testing.T) {
	type args struct {
		prev *url.URL
		base *url.URL
		path string
	}
	tests := []struct {
		name string
		args func(t *testing.T) args

		want1 *url.URL
	}{
		{"No Path", func(t *testing.T) args {
			prev, _ := url.Parse("http://localhost")
			return args{prev: prev}
		}, &url.URL{Scheme: "http", Host: "localhost"}},
		{"Path=#", func(t *testing.T) args {
			prev, _ := url.Parse("http://localhost")
			return args{prev: prev}
		}, &url.URL{Scheme: "http", Host: "localhost"}},
		{"No Base", func(t *testing.T) args {
			prev, _ := url.Parse("http://localhost")
			return args{prev: prev, path: "path"}
		}, &url.URL{Scheme: "http", Host: "localhost", Path: "path"}},
		{"With Base=base/", func(t *testing.T) args {
			prev, _ := url.Parse("http://localhost")
			base, err := url.Parse("base/")
			assert.NoError(t, err)
			return args{prev: prev, path: "path", base: base}
		}, &url.URL{Scheme: "http", Host: "localhost", Path: "base/path"}},
		{"With Base=/base", func(t *testing.T) args {
			prev, _ := url.Parse("http://localhost")
			base, err := url.Parse("/base/")
			assert.NoError(t, err)
			return args{prev: prev, path: "path", base: base}
		}, &url.URL{Scheme: "http", Host: "localhost", Path: "base/path"}},
		{"With Base to different host", func(t *testing.T) args {
			prev, _ := url.Parse("http://localhost")
			base, err := url.Parse("https://base")
			assert.NoError(t, err)
			return args{prev: prev, path: "path", base: base}
		},&url.URL{Scheme: "https", Host: "base", Path: "path"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tArgs := tt.args(t)

			got1 := AbsoluteUrl(tArgs.prev, tArgs.base, tArgs.path)

			if !reflect.DeepEqual(got1.String(), tt.want1.String()) {
				t.Errorf("AbsoluteUrl got1 = %v, want1: %v", got1, tt.want1)
			}
		})
	}
}
