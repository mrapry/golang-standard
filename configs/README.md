
# `Register Additional Environtment`
You can set additional Environtment in file `environtment.go` for example :
```
package configs

import "os"

// Environment additional in this service
type Environment struct {
	KafkaTopicXXX string

var env Environment

// GetEnv get global additional environment
func GetEnv() Environment {
	return env
}

func loadAdditionalEnv() {
	var ok bool

	env.KafkaTopicXXX, ok = os.LookupEnv("KAFKA_TOPIC_XXX")
	if !ok {
		panic("missing KAFKA_TOPIC_XXX environment")
	}

}

```

For call in another method / file go, you just call method config like :
```
configs.GetEnv().KafkaTopicXXX
```



# `/configs`

Configuration file templates or default configs.

Put your `confd` or `consul-template` template files here.