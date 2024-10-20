# roci

`roci` (**R**educed **OCI** Runtime) is an experimental container runtime that implements only a subset of the OCI runtime specification, yet strives to maintain interoperability with container managers such as Podman.
This project is part of my Bachelor's thesis and is **not intended for use in a production environment**.

The bachelor's thesis was graded as "very good". You can access the thesis at [/thesis](/thesis/thesis_m4schini.pdf).

## Table of Contents

- [Installation](#installation)
- [Usage](#usage)
- [Building the Project](#building-the-project)
- [Benchmarking](#benchmarking)
- [License](#license)

## Installation

1. **Go**: Ensure that Go is installed on your system. 

2. **Buf CLI**: The project uses the `buf` CLI to generate protobuf code. Install it by following the instructions on the [buf.build](https://buf.build/) website.

For other platforms, please refer to the official installation guide.

## Usage
### Generating Code with buf

To generate Go code from your .proto files, use the following command:

```shell
buf generate
```

This command will read the configuration in your buf.yaml file and generate Go code into the appropriate directory.
## Building the Project
### Compiling

To compile the project without the benchmark tag, simply use the go build command:

```shell
go build -o roci
```

### Compiling (with verbose output)

To include benchmarking code in the compilation, use the -tags flag:

```shell
go build -tags=verbose
```

This will compile the project with any additional code or optimizations enabled by the benchmark tag.
## Benchmarking

To run benchmarks, use the provided run-benchmark.sh script located in the benchmark/ directory. This script executes the benchmark tests and outputs the results.

Run the script with:

```shell
./benchmark/run-benchmark.sh
```

Ensure that you have compiled the project with the benchmark tag enabled to get accurate benchmarking results.
## License

This project is licensed under the MIT License. See the LICENSE file for details.