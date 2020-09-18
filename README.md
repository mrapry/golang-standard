# Golang Standart

## Made with
<p align="center">
  <img src="https://storage.googleapis.com/agungdp/static/logo/golang.png" width="80" alt="golang logo" />
  <img src="https://storage.googleapis.com/agungdp/static/logo/docker.png" width="80" hspace="10" alt="docker logo" />
  <img src="https://storage.googleapis.com/agungdp/static/logo/rest.png" width="80" hspace="10" alt="rest logo" />
  <img src="https://storage.googleapis.com/agungdp/static/logo/graphql.png" width="80" alt="graphql logo" />
  <img src="https://storage.googleapis.com/agungdp/static/logo/grpc.png" width="160" hspace="15" vspace="15" alt="grpc logo" />
  <img src="https://storage.googleapis.com/agungdp/static/logo/kafka.png" height="80" alt="kafka logo" />
</p>

This repository explain implementation of Go for building multiple microservices using a single codebase. Using [Standard Golang Project Layout](https://github.com/golang-standards/project-layout) and [Clean Architecture](https://blog.cleancoder.com/uncle-bob/2012/08/13/the-clean-architecture.html)


### Go-Lib
This service need [go-lib](https://github.com/mrapry/go-lib), and need add to private library in go
```
go env -w GOPRIVATE="github.com/mrapry/go-lib"
```

Install dependency:
```
$ env GIT_TERMINAL_PROMPT=1 go get github.com/mrapry/go-lib  //or
$ go get github.com/mrapry/go-lib
```
### Create new service
We need someting parameter for create this service

`service` : name of service (required)

`module` : name of module (example: book, book-supplier) use `dash (-)`, dont use space

`gomod` : name init go module (required)
```
make init service={{service_name}} modules={{module_a}},{{module_b}} gomod={{name_init_go_module}}
```

### Run service
```
make run service={{service_name}} gomod={{name_init_go_module}}
```

### Clear service
```
make clear service={{service_name}} gomod={{name_init_go_module}}
```

### Test service
For using script test. We need you initiate git in your work. Because in script test for get `jsonschema root` we check in repository use library :

"github.com/integralist/go-findroot/find"
```
// Repo uses git via the console to locate the top level directory
func Repo() (Stat, error) {
	path, err := rootPath()
	if err != nil {
		return Stat{
			"Unknown",
			"./",
		}, err
	}

	gitRepo, err := exec.Command("basename", path).Output()
	if err != nil {
		return Stat{}, err
	}

	return Stat{
		strings.TrimSpace(string(gitRepo)),
		path,
	}, nil
}
```

Just add `git init` in your repo, and you can run this :)

```
make test service={{service_name}} gomod={{name_init_go_module}}
```

### Note 
Dont forget to register module in internal service `/internal/service.go`
```
// NewService in this service
func NewService(serviceName string, cfg *config.Config) factory.ServiceFactory {
	deps := configs.LoadConfigs(cfg)

	modules := []factory.ModuleFactory{
		// name_of_module.NewModule(deps),   <<- REGISTER YOUR MODULE IN HERE
	}

	return &Service{
		deps:    deps,
		modules: modules,
		name:    types.Service(serviceName),
	}
}
```