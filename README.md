# system model

The System Model component provides a source of truth for the main entities of the system. It is intended usage is that
high-level managers will perform CRUD operations, while lower level components may have read/update access if required.

## Server

To launch the system model execute:

```
$ ./bin/system-model run
{"level":"info","time":"2018-09-28T15:16:30+02:00","message":"Launching API!"}
{"level":"info","port":8800,"time":"2018-09-28T15:16:30+02:00","message":"gRPC port"}
{"level":"info","UseInMemoryProviders":true,"time":"2018-09-28T15:16:30+02:00","message":"Using in-memory providers"}
{"level":"info","port":8800,"time":"2018-09-28T15:16:30+02:00","message":"Launching gRPC server"}
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
