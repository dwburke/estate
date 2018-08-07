# Lode
Preferences or configuration database, with the ability to specify
global settings, and per-user or per-server / region /data center/etc overrides.

## Pretty much fully functional
- delete isn't safe in some ways yet, and has no tests committed

For example, say you have an application (or microservice) that you you want all
users to have the same value for the default, but each user could override it
with their own value.  A flag, a server address, a plugin version...

Configuration is handled by Viper so file formats are flexible...

Sample yaml config:

```
  lode:
    port: 4441
    https: true
    search:
    - "{context}.{app}.{customer_id}.{key}"
    - "{context}.{app}.{key}"
    storage:
      type: "mysql"
      dns: "dbuser:password@/lode?charset=utf8"
      table: "lode"
```

Example usage:

  `GET /lode/dev/foo?app=someapp&customer_id=123456`

From the above config, it will search the keys in top down priority.  If 
customer_id is not in the query, or the derrived key does not have a value,
it will move on to the next key down.

":context" and ":key" are part of the http routing ('dev' and 'foo').
Everything else is optional.

# Adding to your own servce

```
import (
    lode_routes "github.com/dwburke/lode/api/key"
)
...
    lode_routes.SetupRoutes(r)
...

```

# Environment Variables

Configuration can be set with environment variables.  This is useful for containers where
you don't always have easy access to a config file.


```
  export LODE_STORAGE_TYPE=memory
```


# Supported Storage Types
- memory (mainly for testing)
- mysql

# Roadmap
- multiple storage types:
  - postgres
  - etcd
  - redis
