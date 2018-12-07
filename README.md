# system model

The System Model component provides a source of truth for the main entities of the system. It is intended usage is that
high-level managers will perform CRUD operations, while lower level components may have read/update access if required.

## Server

To launch the system model execute:

```
"level":"info","time":"2018-12-03T10:53:57Z","message":"Launching API!"}
{"level":"info","app":"v0.1.0","commit":"d92f8385efaebc6fa75316bb9aed9994ed03fee9","time":"2018-12-03T10:53:57Z","message":"Version"}
{"level":"info","port":8800,"time":"2018-12-03T10:53:57Z","message":"gRPC port"}
{"level":"info","UseDBScyllaProviders":true,"time":"2018-12-03T10:53:57Z","message":"using dbScylla providers"}
{"level":"info","URL":"scylladb.nalej","KeySpace":"nalej","Port":9042,"time":"2018-12-03T10:53:57Z","message":"ScyllaDB"}
{"level":"info","port":8800,"time":"2018-12-03T10:53:57Z","message":"Launching gRPC server"}
```

## CLI

A CLI has been added for convenience, use:

```
$ ./bin/system-model-cli
```

## Kubernetes deploy

Before creating the system model tables, we should deploy scyllaDb with kubernetes (see scylla-deploy project)

Create configMap:
```
$ create -f ./components/system-model/mngtcluster/systemmodel-scylla.configmap.yaml
```
and create the job responsible for the creation of tables
```
kubectl create -f components/system-model/mngtcluster/systemmodel-scylla.job.yaml
```
