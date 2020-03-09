# dsn

We often work with configuration strings like pq database connection:

```
user=mylogin password=mypass database=mydb host=127.0.0.1 port=5432 sslmode=true
```

The simplest way to parse them in Golang is regular expressions:

```golang
package main

import (
	"fmt"
	"regexp"
)

var rex = regexp.MustCompile("(\\w+)=(\\w+)")

func main() {
	conn := `user=mylogin password=mypass database=mydb
                 host=127.0.0.1 port=5432 sslmode=true`

	data := rex.FindAllStringSubmatch(conn, -1)

	res := make(map[string]string)
	for _, kv := range data {
		k := kv[1]
		v := kv[2]
		res[k] = v
	}

	fmt.Println(res)
}
```

But if the value can contain a space character, then you need to add escaping support in a regular expression. And so on, each additional action makes processing harder. But there is an easier way (wmentor/dsn).

Install package:

```
go get github.com/wmentor/dsn
```

Usage:

```golang
package main

import (
  "fmt"

  "github.com/wmentor/dsn"
)

func main() {

  str := `user=mylogin passwd=mypass database=mydb
          port=5432 sslmode=true`

  ds, err := dsn.New(str)
  if err != nil {
    panic("invalid string")
  }n

  // print user=mylogin
  fmt.Printf( "user=%s\n", ds.GetString("user","unknown") )

  // print passwd=mypass
  fmt.Printf( "passwd=%s\n", ds.GetString("passwd","nopass") )

  // host is not exists, print host=127.0.0.1
  fmt.Printf( "host=%s\n", ds.GetString("host","127.0.0.1") )

  // get int value and print port=5432
  fmt.Printf( "port=%d\n", ds.GetInt("port", 4321) )

  // print sslmode=true
  fmt.Printf( "sslmode=%t\n", ds.GetBool("sslmode", false) )

  // print keepalive=false
  fmt.Printf( "keepalive=%t\n", ds.GetBool("keepalive", false) )
}
```

*dns.New* returns object *dsn.DSN* or an error. All get methods (GetString,GetBool,GetInt,GetInt64,GetFloat) take 2 arguments - key name and default value. The default value is used when the key is missing or contains a invalid value.

Moreover, dsn support escape some characters in key name and value (\s,\t,\r,\n,\=,\\,\",\'). See example below:

```golang
package main

import (
  "fmt"

  "github.com/wmentor/dsn"
)

func main() {

  str := `message=Hello,\sWorld! calc=1+1\=2`

  ds, err := dsn.New(str)
  if err != nil {
    panic("invalid string")
  }

  // print message=Hello, World!
  fmt.Printf( "message=%s\n", ds.GetString("message","") )

  // print calc=1+2=2
  fmt.Printf( "calc=%s\n", ds.GetString("calc","") )
}
```

Object dsn.DSN support set methods (SetString,SetBool,SetInt,SetInt64,SetFloat). They take 2 arguments - key name and value:

```golang
ds.SetString("host", "192.168.1.1")
```

Stringer interface implemention makes simple convert dsn.DSN to string:

```golang
str := ds.String()
// or
fmt.Println(ds)
```