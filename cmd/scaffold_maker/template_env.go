package main

const envTemplate = `# Basic env configuration

# Service Handlers
## Server
USE_REST=true
USE_GRPC=false
USE_GRAPHQL=false

## Worker
USE_KAFKA_CONSUMER=false
USE_CRON_SCHEDULER=false
USE_REDIS_SUBSCRIBER=false
USE_TASK_QUEUE_WORKER=false
USE_JEAGER_TRACING = false

## Setting Port 
REST_HTTP_PORT=8000
GRPC_PORT=8002


## Setting Basic Auth
BASIC_AUTH_USERNAME=user
BASIC_AUTH_PASS=pass

## Setting Mongodb Connection
MONGODB_HOST_WRITE=mongodb://localhost:27017
MONGODB_HOST_READ=mongodb://localhost:27017
MONGODB_DATABASE_NAME={{.ServiceName}}


## Setting SQL Connection
SQL_DRIVER_NAME=[string]
SQL_DB_READ_HOST=[string]
SQL_DB_READ_USER=[string]
SQL_DB_READ_PASSWORD=[string]
SQL_DB_WRITE_HOST=[string]
SQL_DB_WRITE_USER=[string]
SQL_DB_WRITE_PASSWORD=[string]
SQL_DATABASE_NAME=[string]

## Setting Redis Connection
REDIS_READ_HOST=localhost
REDIS_READ_PORT=6379
REDIS_READ_AUTH=
REDIS_WRITE_HOST=localhost
REDIS_WRITE_PORT=6379
REDIS_WRITE_AUTH=

## Setting Kafka Connection
KAFKA_BROKERS=localhost:9092
KAFKA_CLIENT_ID={{.ServiceName}}-service-client
KAFKA_CONSUMER_GROUP={{.ServiceName}}-service-consumer-group


## Setting Tracing data with jeager, Graphql, and JsonShema
JAEGER_TRACING_HOST=127.0.0.1:5775
GRAPHQL_SCHEMA_DIR="api/graphql/"
JSON_SCHEMA_DIR="api/jsonschema/"

## Additional env (if you have another env will be set in here)

`
