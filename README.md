# ginject

field's value inject tool for goframe

## Installation

```shell
go get -u github.com/wildchun/ginject
```

## Usage

config.yaml

```yaml
app:
  name: "MyApp"
  version: "1.0.0"
  number: 10
  mqtt:
    broker: "tcp://127.0.0.1:1883"
    clientId: "client"
    username: "admin"
    password: "public"
```


```go
package main

import (
	"github.com/wildchun/ginject"
	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/util/gmeta"
)

type Object struct {
	gmeta.Meta `prefix:"app"`
	appName    string `inject:"name" def:"app_default"`
	version    string `inject:"version" def:"1.0.0"`
	mqtt       struct {
		gmeta.Meta `prefix:"mqtt"`
		broker     string `inject:"broker" def:"tcp://10.147.198.110:1883"`
		clientId   string `inject:"clientId" def:""`
		username   string `inject:"username" def:""`
		password   string `inject:"password" def:""`
	}
	number struct {
		number1 int   `inject:"number" def:"1"`
		number2 int8  `inject:"number" def:"2"`
		number3 int16 `inject:"number" def:"3"`
		number4 int32 `inject:"number" def:"4"`
		number5 int64 `inject:"number" def:"5"`
	}
}

func main(){
	obj := &Object{}
	inj := ginject.New(g.Cfg())
	if err := inj.Apply(obj); err != nil {
		t.Error("injector failed", err)
	}
}

```
