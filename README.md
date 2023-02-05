# toposort

sort a graph topologically.

## Example

```golang
relations := map[string]string{
  "Barbara": "Nick",
  "Nick":    "Sophie",
  "Sophie":  "Jonas",
}

sorted, _ := toposort.Sort(relations)

fmt.Println(sorted)

relations = map[string]string{
  "Jonas": "Jonas",
}

_, err := toposort.Sort(relations)

fmt.Println(err.Error())
```

Output:

```
[Jonas Sophie Nick Barbara]
cyclic: [Jonas Jonas]
```
