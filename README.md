# GO Monitoring Service
monitoring service in go

## Conf package ##
conf package is for configurations.
config file is saved in toml format.
original config file is saved only in master branch.
the sample config file:
```
[ds]
	[ds.redis]
		port = "6379"
		chan = "orders"
[srv]
	port = "1234"
```

To add any sub-info into config file

1. Create struct with the same subtree name in conf/conf.go
    i.e: 
  
```
#!go

    type Datastore struct {
        Redis struct {
            Port string
            Chan string
        }

        Postgres struct{
            login string
            passwd string
        }
    }
```

2. Add info in config.toml
    
```

[ds]
    [ds.redis]
        port = "6379"
        chan = "orders"
    [ds.postgres]
         login = "superuser"
         passwd = "p@$$"
[srv]
    port = "1234"

```

## Rcache(redis cache) ##
Each fleet is save in the following structure

```
    #!go

type Pos struct {
	Id        int
	Latitude  string
	Longitude string
	Time      string
}

type Fleet struct {
	Id     string
	Update map[string]Pos
}
```

This stucture becomes to the following JSON format


```
{
	"Id": "fleet_202",
	"Update": {
		"106206": {
			"Id": 106206,
			"Latitude": "41.260006",
			"Longitude": "69.145340",
			"Time": "2015-04-21 17:59:59"
		},
		"107749": {
			"Id": 107749,
			"Latitude": "41.293964",
			"Longitude": "69.245811",
			"Time": "2015-03-26 13:29:06"
		}
	}
}
```
