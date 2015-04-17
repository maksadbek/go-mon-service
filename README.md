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
