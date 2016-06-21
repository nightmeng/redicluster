# Redis Cluster

RediCluster is a redis cluster package based on redigo.

## Installation

Install RediCluster using the "go get" command:

```golang
go get githup.com/nightmeng/redicluster
```

## Example

The following is a example:

```golang
package main

import(
	"githup.com/nightment/redicluster"
)

func main() {
	cluster := redicluster.NewRediCluster([]string{"redis://127.0.0.1"})
	
	conn, err := cluster.GetConn("xxx")
	if err ! = nil {
		fmt.Printf("Get redis connection failed, %s", err)
		return
	}
	defer conn.Close()
	
	_, err = conn.Do("INCR", "test-incr")
	if err != nil {
		fmt.Printf("INCR failed, %s", err)
		return
	}
}
```
