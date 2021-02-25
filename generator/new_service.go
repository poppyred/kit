package generator

import (
	"bytes"
	"fmt"
	"os/exec"
	"path"
	"strings"

	"github.com/dave/jennifer/jen"
	"github.com/kujtimiihoxha/kit/fs"
	"github.com/kujtimiihoxha/kit/utils"
	"github.com/spf13/viper"
)

// NewService implements Gen and is used to create a new service.
type NewService struct {
	BaseGenerator
	name          string
	interfaceName string
	destPath      string
	filePath      string
}

// NewNewService returns a initialized and ready generator.
//
// The name parameter is the name of the service that will be created
// this name should be without the `Service` suffix
func NewNewService(name string) Gen {
	gs := &NewService{
		name:          name,
		interfaceName: utils.ToCamelCase(name + "Service"),
		destPath:      fmt.Sprintf(viper.GetString("gk_service_path_format"), utils.ToLowerSnakeCase(name)),
	}
	gs.filePath = path.Join(gs.destPath, viper.GetString("gk_service_file_name"))
	gs.srcFile = jen.NewFilePath(strings.Replace(gs.destPath, "\\", "/", -1))
	gs.InitPg()
	gs.fs = fs.Get()
	return gs
}

// Generate will run the generator.
func (g *NewService) Generate() error {
	err := g.CreateFolderStructure(g.destPath)
	if err != nil {
		println(err.Error())
		return err
	}
	err = g.genModule()
	if err != nil {
		println(err.Error())
		return err
	}

	comments := []string{
		"Add your methods here",
		"e.x: Foo(ctx context.Context,s string)(rs string, err error)",
	}
	partial := NewPartialGenerator(nil)
	partial.appendMultilineComment(comments)
	g.code.Raw().Commentf("%s describes the service.", g.interfaceName).Line()
	g.code.appendInterface(
		g.interfaceName,
		[]jen.Code{jen.Id("Status").Params(
			jen.Id("ctx").Qual("context", "Context"),
		).Params(jen.Id("err").Error()),
			partial.raw},
	)

	return g.fs.WriteFile(g.filePath, g.srcFile.GoString(), false)
}

func (g *NewService) genModule() error {
	prjName := utils.ToLowerSnakeCase(g.name)
	exist, _ := g.fs.Exists(prjName + "/go.mod")
	if exist {
		return nil
	}

	moduleName := prjName
	if viper.GetString("n_s_module") != "" {
		moduleName = viper.GetString("n_s_module")
		moduleNameSlice := strings.Split(moduleName, "/")
		moduleNameSlice[len(moduleNameSlice)-1] = utils.ToLowerSnakeCase(moduleNameSlice[len(moduleNameSlice)-1])
		moduleName = strings.Join(moduleNameSlice, "/")
	}
	cmdStr := "cd " + prjName + " && go mod init " + moduleName + " && go mod vendor"
	cmd := exec.Command("sh", "-c", cmdStr)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	_, err := cmd.Output()
	// return cmd.Stderr to debug (err here provides nothing useful, only `exit status 1`)
	if err != nil {
		return fmt.Errorf("genModule: sh -c %s => err:%v", cmdStr, err.Error()+" , "+stderr.String())
	}
	return nil
}
