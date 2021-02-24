package generator

import (
	"github.com/dave/jennifer/jen"
	"github.com/kujtimiihoxha/kit/fs"
	"path"
	"testing"

	"runtime"

	"github.com/kujtimiihoxha/kit/parser"
	"github.com/kujtimiihoxha/kit/utils"
	"github.com/spf13/viper"
)

func setDefaults() {
	viper.SetDefault("gk_service_path_format", path.Join("%s", "pkg", "service"))
	viper.SetDefault("gk_cmd_service_path_format", path.Join("%s", "cmd", "service"))
	viper.SetDefault("gk_cmd_path_format", path.Join("%s", "cmd"))
	viper.SetDefault("gk_endpoint_path_format", path.Join("%s", "pkg", "endpoint"))
	viper.SetDefault("gk_http_path_format", path.Join("%s", "pkg", "http"))
	viper.SetDefault("gk_http_client_path_format", path.Join("%s", "client", "http"))
	viper.SetDefault("gk_grpc_client_path_format", path.Join("%s", "client", "grpc"))
	viper.SetDefault("gk_client_cmd_path_format", path.Join("%s", "cmd", "client"))
	viper.SetDefault("gk_grpc_path_format", path.Join("%s", "pkg", "grpc"))
	viper.SetDefault("gk_grpc_pb_path_format", path.Join("%s", "pkg", "grpc", "pb"))

	viper.SetDefault("gk_service_file_name", "service.go")
	viper.SetDefault("gk_service_middleware_file_name", "middleware.go")
	viper.SetDefault("gk_endpoint_base_file_name", "endpoint_gen.go")
	viper.SetDefault("gk_endpoint_file_name", "endpoint.go")
	viper.SetDefault("gk_endpoint_middleware_file_name", "middleware.go")
	viper.SetDefault("gk_http_file_name", "handler.go")
	viper.SetDefault("gk_http_base_file_name", "handler_gen.go")
	viper.SetDefault("gk_cmd_base_file_name", "service_gen.go")
	viper.SetDefault("gk_cmd_svc_file_name", "service.go")
	viper.SetDefault("gk_http_client_file_name", "http.go")
	viper.SetDefault("gk_grpc_client_file_name", "grpc.go")
	viper.SetDefault("gk_grpc_pb_file_name", "%s.proto")
	viper.SetDefault("gk_grpc_base_file_name", "handler_gen.go")
	viper.SetDefault("gk_grpc_file_name", "handler.go")
	if runtime.GOOS == "windows" {
		viper.SetDefault("gk_grpc_compile_file_name", "compile.bat")
	} else {
		viper.SetDefault("gk_grpc_compile_file_name", "compile.sh")
	}
	viper.SetDefault("gk_service_struct_prefix", "basic")
	viper.Set("gk_testing", true)

}
func createTestMethod(name string, param []parser.NamedTypeValue, result []parser.NamedTypeValue) parser.Method {
	param = append(param, parser.NewNameType("ctx", "context.Context"))
	return parser.Method{
		Name:       name,
		Parameters: param,
		Results:    result,
	}
}
func getTestServiceInterface(name string) parser.Interface {
	n := utils.ToCamelCase(name + "_Service")
	return parser.NewInterface(n, []parser.Method{
		createTestMethod(
			"Foo",
			[]parser.NamedTypeValue{
				parser.NewNameType("s", "string"),
			},
			[]parser.NamedTypeValue{
				parser.NewNameType("r", "string"),
				parser.NewNameType("err", "error"),
			},
		),
	})
}

// fix #7
func TestBaseGenerator_AddImportsToFile(t *testing.T) {
	type fields struct {
		srcFile *jen.File
		code    *PartialGenerator
		fs      *fs.KitFs
	}
	type args struct {
		imp []parser.NamedTypeValue
		src string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name:   "test gen dmw",
			fields: fields{},
			args: args{
				imp: []parser.NamedTypeValue{
					{
						Name:  "",
						Type:  "\"context\"",
						Value: "",
					}, {
						Name:  "log",
						Type:  "\"github.com/go-kit/kit/log\"",
						Value: "",
					},
				},
				src: "package service\n\ntype Middleware func(TestService) TestService\n\n\ttype loggingMiddleware struct {\n\t\tlogger log.Logger\n\t\tnext   TestService\n\t}\n\n\t// LoggingMiddleware takes a logger as a dependency\n\t// and returns a TestService Middleware.\n\tfunc LoggingMiddleware(logger log.Logger) Middleware {\n\t\treturn func(next TestService) TestService {\n\t\t\treturn &loggingMiddleware{logger, next}\n\t\t}\n\n\t}\n\n\tfunc (l loggingMiddleware) Foo(ctx context.Context, s string) (rs string, err error) {\n\t\tdefer func() {\n\t\t\tl.logger.Log(\"method\", \"Foo\", \"s\", s, \"rs\", rs, \"err\", err)\n\t\t}()\n\t\treturn l.next.Foo(ctx, s)\n\t} \n",
			},
			want: "package service\n\nimport (\n\t\"context\"\n\tlog \"github.com/go-kit/kit/log\"\n)\n\ntype Middleware func(TestService) TestService\n\ntype loggingMiddleware struct {\n\tlogger log.Logger\n\tnext   TestService\n}\n\nfunc LoggingMiddleware(logger log.Logger) Middleware {\n\treturn func(next TestService) TestService {\n\t\treturn &loggingMiddleware{logger, next}\n\t}\n\n}\n\nfunc (l loggingMiddleware) Foo(ctx context.Context, s string) (rs string, err error) {\n\tdefer func() {\n\t\tl.logger.Log(\"method\", \"Foo\", \"s\", s, \"rs\", rs, \"err\", err)\n\t}()\n\treturn l.next.Foo(ctx, s)\n}\n",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BaseGenerator{
				srcFile: tt.fields.srcFile,
				code:    tt.fields.code,
				fs:      tt.fields.fs,
			}
			got, err := b.AddImportsToFile(tt.args.imp, tt.args.src)
			if (err != nil) != tt.wantErr {
				t.Errorf("AddImportsToFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("AddImportsToFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestBaseGenerator_CreateFolderStructure(t *testing.T) {
	type fields struct {
		srcFile *jen.File
		code    *PartialGenerator
		fs      *fs.KitFs
	}
	type args struct {
		path string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "foo",
			fields: fields{
				srcFile: nil,
				code:    nil,
				fs:      fs.Get(),
			},
			args: args{
				path: "test2/pkg/service",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			b := &BaseGenerator{
				srcFile: tt.fields.srcFile,
				code:    tt.fields.code,
				fs:      tt.fields.fs,
			}
			if err := b.CreateFolderStructure(tt.args.path); (err != nil) != tt.wantErr {
				t.Errorf("CreateFolderStructure() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
