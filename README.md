# System-model

The System Model component provides a source of truth for the main entities of the system. Its intended usage is that
high-level managers can perform CRUD operations, while lower-level components may have read/update access if required.

## Getting Started

This component is divided into three big sections: entities, providers and server.
* **Entities**: contains all entity definitions, functions for converting them to `grpc structs` and viceversa, and request validations.
* **Providers**: contains all the providers required. Two types of providers exist: _mockup_ (or memory) providers and _scylladb_ providers. Each Scylla provider has an equivalent Mockup provider, and both must pass the same tests, so integration tests on the servers are not needed.
* **Server**: contains all the logic of the component.

### Prerequisites

To run system-model, we need a [**ScyllaDB**](https://www.scylladb.com/open-source/) installation.

### Build and compile

In order to build and compile this repository use the provided Makefile:

```
make all
```

This operation generates the binaries for this repo, downloads the required dependencies, runs existing tests and generates ready-to-deploy Kubernetes files.

### Run tests

Tests are executed using Ginkgo. To run all the available tests:

```
make test
```

### Update dependencies

Dependencies are managed using Godep. For an automatic dependencies download use:

```
make dep
```

In order to have all dependencies up-to-date run:

```
dep ensure -update -v
```

## Integration test
Some integration tests are included. The execution of integration tests may have collateral effects on the state of the platform. **DO NOT execute those tests in production**. After each test, the tables are truncated.

​The following table contains the variables that activate the integration tests:

| Variable  | Example Value | Description |
 | ------------- | ------------- |------------- |
 | RUN_INTEGRATION_TEST  | true | Run integration tests |
 | IT_SCYLLA_HOST  | 127.0.0.1 | Scylla Address |
 | IT_SCYLLA_PORT | 9042 | Scylla Port |
 | IT_NALEJ_KEYSPACE | nalej | Keyspace name |

The database must be created to run the integration test. There is a file (`scripts/database.cql`) that contains all the 
sentences to create the keyspace and the tables needed for it.

## Known Issues

All the operations related to accounts and projects are not available yet. They will be ready in future releases. 

## Contributing

Please read [contributing.md](contributing.md) for details on our code of conduct, and the process for submitting pull requests to us.


## Versioning

We use [SemVer](http://semver.org/) for versioning. For the available versions, see the [tags on this repository](https://github.com/nalej/system-model/tags). 

## Authors

See also the list of [contributors](https://github.com/nalej/system-model/contributors) who participated in this project.

## License
This project is licensed under the Apache 2.0 License - see the [LICENSE-2.0.txt](LICENSE-2.0.txt) file for details.
