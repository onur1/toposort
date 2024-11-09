# toposort

**toposort** provides functionality to perform topological sorting on directed acyclic graphs (DAGs). Topological sorting is the linear ordering of vertices such that for every directed edge `U → V`, vertex `U` comes before vertex `V` in the ordering. This is particularly useful in scenarios like task scheduling, resolving symbol dependencies in linkers, and determining compilation order in programming languages.

## Features

- **Topological Sorting**: Generates a linear ordering of vertices in a DAG.
- **Cycle Detection**: Identifies cycles in the graph and returns an error if a cycle is detected, as topological sorting is only possible for acyclic graphs.

## Installation

To install the `toposort` package, use the following command:

```sh
go get github.com/onur1/toposort
```

## Usage

```go
package main

import (
    "fmt"
    "github.com/onur1/toposort"
)

func main() {
    // Define the relationships in the graph
    relations := map[string]string{
        "Barbara": "Nick",
        "Nick":    "Sophie",
        "Sophie":  "Jonas",
    }

    // Perform topological sorting
    sorted, err := toposort.Sort(relations)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }

    fmt.Println("Sorted order:", sorted)
}
```

**Output:**

```
Sorted order: [Jonas Sophie Nick Barbara]
```

In this example, the `relations` map defines a set of dependencies where each key depends on its corresponding value. The `Sort` function processes these relationships and returns a slice of strings representing the topologically sorted order.

## Error Handling

If the graph contains cycles, the `Sort` function will return an error indicating the presence of a cycle. For example:

```go
relations := map[string]string{
    "Jonas": "Jonas",
}

_, err := toposort.Sort(relations)
if err != nil {
    fmt.Println("Error:", err)
}
```

**Output:**

```
Error: cyclic: [Jonas Jonas]
```

This indicates that a cycle exists in the graph, making topological sorting impossible.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.

