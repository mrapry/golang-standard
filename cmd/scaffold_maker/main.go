package main

import (
	"bytes"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

const (
	packageName    = "mrapry"
	libraryAddress = "github.com/mrapry/go-lib"
)

type param struct {
	PackageName    string
	GoModules      string
	LibraryAddress string
	ServiceName    string
	Modules        []string
}

// FileStructure model
type FileStructure struct {
	TargetDir    string
	IsDir        bool
	FromTemplate bool
	DataSource   interface{}
	Source       string
	FileName     string
	Skip         bool
	Childs       []FileStructure
}

var (
	tpl *template.Template
)

func main() {

	var scope string
	var serviceName string
	var modulesFlag string
	var goMod string

	flag.StringVar(&scope, "scope", "initservice", "set scope")
	flag.StringVar(&serviceName, "servicename", "", "set service name")
	flag.StringVar(&modulesFlag, "modules", "", "set all modules from service")
	flag.StringVar(&goMod, "gomod", "", "set go modules in this project")

	flag.Usage = func() {
		fmt.Println("-scope | --scope => set scope (initservice or addmodule), example: --scope initservice")
		fmt.Println("-servicename | --servicename => set service name, example: --servicename master-service")
		fmt.Println("-modules | --modules => set modules name, example: --modules user,auth")
		fmt.Println("-gomod | --gomod => set init go modules, example: --gomod master-service")
	}

	flag.Parse()

	var data param
	data.PackageName = packageName
	data.GoModules = goMod
	data.LibraryAddress = libraryAddress
	data.ServiceName = serviceName

	tpl = template.New(packageName)

	modules := strings.Split(modulesFlag, ",")
	if modulesFlag == "" {
		modules = []string{"module"} // default module name
	}

	sort.Slice(modules, func(i, j int) bool {
		return modules[i] < modules[j]
	})

	//init gomodules
	gomodInitStructure := FileStructure{
		// TargetDir: "", IsDir: false, DataSource: data,
		Childs: []FileStructure{
			{FromTemplate: true, DataSource: data, Source: gomodTemplate, FileName: "go.mod"},
		},
	}

	// init api
	// apiStructure := FileStructure{
	// 	TargetDir: "api/", IsDir: true, DataSource: data,
	// }

	// init configs
	configsStructure := FileStructure{
		TargetDir: "configs/", IsDir: true, DataSource: data,
		Childs: []FileStructure{
			{FromTemplate: true, DataSource: data, Source: configTemplate, FileName: "config.go"},
			{FromTemplate: true, DataSource: data, Source: configLoadEnvTemplate, FileName: "environment.go"},
		},
	}

	// init cmd
	cmdStructure := FileStructure{
		TargetDir: "cmd/{{.ServiceName}}/", IsDir: true, DataSource: data,
		Childs: []FileStructure{
			{FromTemplate: true, DataSource: data, Source: cmdMainTemplate, FileName: "main.go"},
			{FromTemplate: true, DataSource: data, Source: envTemplate, FileName: ".env"},
			{FromTemplate: true, DataSource: data, Source: envTemplate, FileName: ".env.development"},
		},
	}

	//init service internal
	serviceStructure := FileStructure{
		TargetDir: "internal/", IsDir: true, DataSource: data,
		Childs: []FileStructure{
			{FromTemplate: true, DataSource: data, Source: serviceMainTemplate, FileName: "service.go"},
		},
	}

	//init modules
	var moduleStructure = FileStructure{
		TargetDir: "modules/", IsDir: true, DataSource: data,
	}

	for _, moduleName := range modules {
		moduleName = strings.TrimSpace(moduleName)
		data.Modules = append(data.Modules, moduleName)
		dataSource := map[string]string{"PackageName": data.PackageName, "ServiceName": data.ServiceName, "module": moduleName, "GoModules": data.GoModules, "LibraryAddress": data.LibraryAddress}

		//init clean architecture module directory
		cleanArchModuleDir := []FileStructure{
			{
				TargetDir: "delivery/", IsDir: true,
				Childs: []FileStructure{
					{TargetDir: "resthandler/", IsDir: true, Childs: []FileStructure{
						{FromTemplate: true, DataSource: dataSource, Source: deliveryRestTemplate, FileName: "resthandler.go"},
					}},
				},
			},
			{
				TargetDir: "domain/", IsDir: true,
				Childs: []FileStructure{
					{FromTemplate: true, FileName: "domain.go"},
				},
			},
			{
				TargetDir: "repository/", IsDir: true,
				Childs: []FileStructure{
					{TargetDir: "interfaces/", IsDir: true},
					{FromTemplate: true, FileName: "repository.go"},
				},
			},
			{
				TargetDir: "usecase/", IsDir: true,
				Childs: []FileStructure{
					{FromTemplate: true, FileName: "usecase.go"},
					{FromTemplate: true, FileName: "usecase_impl.go"},
				},
			},
		}

		moduleStructure.Childs = append(moduleStructure.Childs, []FileStructure{
			{
				TargetDir: moduleName + "/", IsDir: true,
				Childs: append(cleanArchModuleDir,
					FileStructure{
						FromTemplate: true, DataSource: dataSource, Source: moduleMainTemplate, FileName: "module.go",
					},
				),
			},
		}...)
	}

	serviceStructure.Childs = append(serviceStructure.Childs, moduleStructure)

	var baseDirectoryFile FileStructure
	switch scope {
	case "initservice":
		baseDirectoryFile.Childs = []FileStructure{
			gomodInitStructure, configsStructure, cmdStructure, serviceStructure,
		}

	// case "addmodule":
	// 	moduleStructure.Skip = true
	// 	serviceStructure.Skip = true
	// 	serviceStructure.Childs = []FileStructure{
	// 		moduleStructure,
	// 		{FromTemplate: true, DataSource: data, Source: serviceMainTemplate, FileName: "service.go"},
	// 	}

	// 	// apiStructure.Skip = true
	// 	// apiProtoStructure.Skip, apiGraphQLStructure.Skip = true, true
	// 	// apiStructure.Childs = []FileStructure{
	// 	// 	apiProtoStructure, apiGraphQLStructure,
	// 	// }

	// 	baseDirectoryFile.Childs = []FileStructure{apiStructure, serviceStructure}
	// 	baseDirectoryFile.Skip = true

	default:
		panic("invalid scope parameter")
	}

	exec(baseDirectoryFile)

}

//exec for generate file
func exec(fl FileStructure) {
	dirBuff := loadTemplate(fl.TargetDir, fl.DataSource)
	dirName := string(dirBuff)

	if fl.Skip {
		goto execChild
	}

	if _, err := os.Stat(dirName); os.IsExist(err) {
		panic(err)
	}

	if fl.IsDir {
		_, err := os.Stat(dirName)
		if os.IsNotExist(err) {
			fmt.Printf("creating %s...\n", dirName)
			if errDir := os.Mkdir(dirName, 0700); errDir != nil {
				fmt.Println("mkdir err:", errDir)
				panic(errDir)
			}
		}

	}

	if fl.FileName != "" {
		var buff []byte
		if fl.FromTemplate {
			if fl.Source != "" {
				buff = loadTemplate(fl.Source, fl.DataSource)
			} else {
				lastDir := filepath.Dir(fl.TargetDir)
				buff = defaultDataSource(lastDir[strings.LastIndex(lastDir, "/")+1:])
			}
		} else {
			buff = []byte(fl.Source)
		}
		if len(dirName) > 0 {
			dirName = strings.TrimSuffix(dirName, "/")
			if err := ioutil.WriteFile(dirName+"/"+fl.FileName, buff, 0644); err != nil {
				panic(err)
			}
		} else {
			if err := ioutil.WriteFile(fl.FileName, buff, 0644); err != nil {
				panic(err)
			}
		}

	}

execChild:
	for _, child := range fl.Childs {
		child.TargetDir = dirName + child.TargetDir
		exec(child)
	}
}

//loadTemplate function for mapping data template
func loadTemplate(source string, sourceData interface{}) []byte {
	var byteBuff = new(bytes.Buffer)
	defer byteBuff.Reset()

	tmpl, err := tpl.Funcs(formatTemplate()).Parse(source)
	if err != nil {
		panic(err)
	}

	if err := tmpl.Execute(byteBuff, sourceData); err != nil {
		panic(err)
	}

	return byteBuff.Bytes()
}

//formatTemplate function remapping variable in template
func formatTemplate() template.FuncMap {
	replacer := strings.NewReplacer("-", "", "*", "", "/", "", ":", "")
	return template.FuncMap{

		"clean": func(v string) string {
			return replacer.Replace(v)
		},

		"upper": func(str string) string {
			return strings.Title(str)
		},
	}
}
