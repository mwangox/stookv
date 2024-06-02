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
          "value": "kivyao*2024"
        }'
```

###### Set Secret key
```shell
curl -X POST --location "http://localhost:9098/stoo-kv/secrets/my-app/prod" \
    -H "Content-Type: application/json" \
    -d '{
          "key": "database.password",
          "value": "kivyetu*2024"
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
General stookv configurations are stored in `stoo_kv.json` and storage provider specific configurations are stored in `provider.json`. 

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


Sample configurations for each of the supported storage providers are shown in [provider.json](./conf/provider.json). 
You may remove the configurations for provider(s) which you don't need in your setup.
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
The following is the list of currently supported storage types, with more to be added in future releases.
- Redis
- Mongo
- Etcd
- MySQL
- MariaDB
- Postgres
- Memory

Storage type is specified in the configuration file under key `storage_type`. In case the storage type is not specified explicitly it will default to `memory`.
If you want to add your storage implementation that is not available in the list above, just implement the [store](./internal/store/store.go) interface and add it accordingly.

### Available StooKV SDKs
You can use `StooKV` without use of these SDKs by simply calling the REST or gRPC APIs using any tool of your choice. If you don't want to bother
with the underlying low-level implementations, you can use any of these based on your preferred language:

- Go: [stogo](https://github.com/mwangox/stogo)
- Java: [stoja](https://github.com/mwangox/stoja)
- Rust: [storus](https://github.com/mwangox/stogo)
- Spring Boot Starter: [stoja-spring-boot-starter](https://github.com/mwangox/stoja-spring-boot-starter) (Java framework)

For those who want to implement their own SDK(s), I recommend using gRPC APIs instead of REST APIs due to their associated benefits. All the SDKs mentioned above
use gRPC APIs to interact with the `stookv` instance. Please visit the respective repository for more details on the SDK usage. n

### Secrets Encryption
`StooKV` supports encryption of values if one needs to store configurations that should not be plain/visible like passwords, tokens, auth keys etc.
Stookv is using **AES** with **GCM** mode to encrypt values. To use this feature, firstly, you need to set encryption key `encrypt_key` in the configuration file. The key length should be
either `128`, `192`, or `256` bits. Then you can use any of your preferred APIs to set your secrets.
Also, one can opt to use manual encryption or decryption REST API endpoints as:
###### Encrypt
```shell
curl -X POST --location "http://localhost:9098/stoo-kv/encrypt" \
    -H "Content-Type: text/plain" \
    -d 'sote*2024'
```

###### Decrypt
```shell
curl -X POST --location "http://localhost:9098/stoo-kv/decrypt" \
    -H "Content-Type: text/plain" \
    -d '48fa702f0614a5550a4ebf98e2541e8708afe23bce365d14c100d1b7d1c455534e433ed32867ffdfdf'
```
The decryption endpoint is not enabled by default, you need to enable it in the configuration file before using it.


### Installation
To use `stookv` you need to download binary/archive from release page based on your target operating system. Optionally, you can build stookv from
sources as:
```shell 
go build -o stookv
```
After you have downloaded the binary or build from sources, you can run stookv using this command:
```shell
./stookv --config.file=/path/stoo_kv.json
```
### License

The project is licensed under [MIT license](./LICENSE).

### Contribution

Unless you explicitly state otherwise, any contribution intentionally submitted
for inclusion in `stookv` by you, shall be licensed as MIT, without any additional
terms or conditions.