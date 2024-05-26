![](./conf/stooKV.png)

`StooKv` is a key-value datastore written in `Go` that is language agnostic and not limited to one backend storage type.

### Features of StooKV
- Supports multiple storage backends such as MySQL, MongoDB, Postgres, etcd, Redis, In-Memory, etc.
- Out-of-the-box encryption of data with automatic decryption upon retrieval.
- Rest or grpc-based APIs for clients.
- Each key-value pair is organized using the concept of namespace and profile.

### Available API Operations

| Rest Endpoint                                           | Rest Method | Grpc Method                     | Description                                                   |
|:--------------------------------------------------------|:------------|---------------------------------|---------------------------------------------------------------|
| {host:port}/stoo-kv/{namespace}/{profile}/{key}         | GET         | GetService                      | Reads value from a key.                                       |
| {host:port}/stoo-kv/{namespace}/{profile}               | GET         | GetServiceByNamespaceAndProfile | Get all key-value pairs from the given namespace and profile. |
| {host:port}/stoo-kv/{namespace}/{profile}               | POST        | SetKeyService                   | Sets a value to a given key.                                  |
| {host:port}/stoo-kv/secrets/{namespace}                 | POST        | SetSecretKeyService             | Sets value as secret to a given key.                          |
| {host:port}/stoo-kv/{namespace}/{profile}?{key}={value} | DELETE      | DeleteKeyService                | Removes a key from the datastore.                             | 
| {host:port}/stoo-kv/encrypt	                            | POST	       | -                               | Manual encrypt data.                                          |
| {host:port}/stoo-kv/decrypt	                            | POST	       | -                               | Manual decrypt data.                                          |

### Rest API USAGE Examples

###### Set plain key
```shell 
curl -X POST --location "http://localhost:9098/stoo-kv/my-app/prod" \
    -H "Content-Type: application/json" \
    -d '{
          "key": "database.password",
          "value": "123456aaa*"
        }'
```

###### Set Secret key
```shell
curl -X POST --location "http://localhost:9098/stoo-kv/secrets/my-app/prod" \
    -H "Content-Type: application/json" \
    -d '{
          "key": "database.password",
          "value": "123456aaa*"
        }'
```

###### Get Key
```shell
GET http://localhost:9098/stoo-kv/my-app/prod/database.password
```
###### Delete Key
```shell
curl -X DELETE --location "http://localhost:9098/stoo-kv/my-app/prod?key=database.password"
```

###### Get All Keys by Namespace and Profile
```shell
curl -X GET --location "http://localhost:9098/stoo-kv/my-app/prod"
```
### Configurations
There are general stookv configurations in `stoo_kv.json` and storage provider specific configurations in `provider.json`. 

General configurations definitions as in [stoo_kv.json](./conf/stoo_kv.json).

| Key                     | Example                               | Description                         |
|-------------------------|---------------------------------------|-------------------------------------|
| `storage_type`          | `mysql`                               | Type of storage used (e.g., mysql)  |
| `server_port`           | `9098`                                | Port number for the server          |
| `server_log_level`      | `debug`                               | Log level for the server            |
| `grpc_port`             | `50051`                               | Port number for gRPC server         |
| `encrypt_key`           | `abcdefghijklmnopqrstuvwxyzaaaaaa`    | Key used for encryption of keys.    |
| `enable_decrypt_endpoint` | `true`                                | Flag to enable decrypt endpoint     |
| `rdbms_default_table`   | `kv_store`                            | Default table name in the RDBMS     |
| `encrypt_prefix`        | `{ENC} `                              | Prefix used for encrypted values    |
| `provider_path`         | `./conf/provider.json`                | Path to provider configuration file |
| `grpc_use_tls`          | `true`                                | Flag to enable TLS for gRPC         |
| `grpc_server_cert`      | `/stoo-kv/grpc/certs/server_cert.pem` | Path to the gRPC server certificate |
| `grpc_server_key`       | `/stoo-kv/grpc/certs/server_key.pem`  | Path to the gRPC server key         |


Storage providers specific configurations as in [provider.json](./conf/provider.json).
###### Redis Configuration
| Key                     | Example     | Description                           |
|-------------------------|-------------|---------------------------------------|
| `host`                  | `localhost` | Hostname for Redis                    |
| `port`                  | `6379`      | Port number for Redis                 |
| `password`              | `""`        | Password for Redis (empty if not set) |
| `database`              | `0`         | Database index for Redis              |
| `store_name`            | `kv_store`  | Name of the store(map) in Redis       |
| `connection_pool_size`  | `10`        | Connection pool size for Redis        |

###### MySQL Configuration
| Key                     | Example     | Description                                    |
|-------------------------|-------------|------------------------------------------------|
| `host`                  | `localhost` | Hostname for MySQL                             |
| `port`                  | `3306`      | Port number for MySQL                          |
| `username`              | `root`      | Username for MySQL                             |
| `password`              | `root`      | Password for MySQL                             |
| `database_name`         | `key_value` | Name of the database in MySQL                  |

###### Postgres Configuration
| Key                     | Example          | Description                                    |
|-------------------------|------------------|------------------------------------------------|
| `host`                  | `localhost`      | Hostname for Postgres                          |
| `port`                  | `3306`           | Port number for Postgres                       |
| `username`              | `root`           | Username for Postgres                          |
| `password`              | `root`           | Password for Postgres                          |
| `database_name`         | `key_value`      | Name of the database in Postgres               |
| `ssl_mode`              | `disabled`       | SSL mode for Postgres                          |
| `timezone`              | `Africa/Nairobi` | Timezone for Postgres                          |

###### MongoDB Configuration

| Key                     | Example                     | Description                                    |
|-------------------------|-----------------------------|------------------------------------------------|
| `mongo_uri`             | `mongodb://localhost:27017` | MongoDB connection URI                         |
| `database_name`         | `key_value`                 | Name of the database in MongoDB                |
| `collection_name`       | `kv_store`                  | Name of the collection in MongoDB              |

###### Etcd Configuration
| Key                     | Example                        | Description                                    |
|-------------------------|--------------------------------|------------------------------------------------|
| `endpoints`             | `["http://155.12.30.14:2379"]` | List of Etcd endpoints                         |
| `username`              | `admin`                        | Username for Etcd                              |
| `password`              | `admin`                        | Password for Etcd                              |
| `dial_timeout`          | `20`                           | Dial timeout for Etcd (in seconds)             |

### Supported Backend Storages
The following is the list of current supported storage types, more to be added in the future releases.
- Redis
- Mongo
- Etcd
- MySQL
- MariaDB
- Postgres
- Memory

Storage type is specified in the configuration file under key `storage_type`. In case if storage type is not
specified explicit it will default to `memory`.
If you want to add your own storage implementation that is not available in the list above, just implement the [store](./internal/store/store.go) interface and add it accordingly.

### Secrets Encryption
StooKV supports encryption of properties if one needs to store configurations that should not be naked like passwords, auth keys etc.
Stookv is using AES with GCM mode to encrypt values. To use this feature, first you need to set encryption key `encrypt_key` in the configuration file. The key length should be
either of 128, 192 or 256 bits. Then you can use any of your preferred API to set your secrets.
Also, one can opt to use manual encryption or decryption REST APIs as:
###### Encrypt
```shell
curl -X POST --location "http://localhost:9098/stoo-kv/encrypt" \
    -H "Content-Type: text/plain" \
    -d '123456a*'
```

###### Decrypt
```shell
curl -X POST --location "http://localhost:9098/stoo-kv/decrypt" \
    -H "Content-Type: text/plain" \
    -d '48fa702f0614a5550a4ebf98e2541e8708afe23bce365d14c100d1b7d1c455534e433ed32867ffdfdf'
```
**NOTE**: Decryption endpoint is not enabled by default, you need to enable it in the configuration file before using it.


### License

The project is licensed under [MIT license](./LICENSE).

### Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted
for inclusion in `stookv` by you, shall be licensed as MIT, without any additional
terms or conditions.