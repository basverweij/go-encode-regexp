# go-encode-regexp

This library provides Unmarshalling functionality using a [regular expression](https://golang.org/pkg/regexp/).

Decoders are initialized by providing a regular expression (compiled or as string).

When unmarshalling the following mapping from the regular expression to the fields in the struct is used:
1. Named capture groups are mapped to fields with the same name in the 'encre' tag.
2. Indexed capture groups are mapped to fields with the corresponding index in the 'encre' tag.
3. Named capture groups are mapped to fields with the same name.
