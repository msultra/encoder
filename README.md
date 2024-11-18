<h1 align="center">
    <br>
    <img style="border-radius: 50%;" src="https://github.com/msultra.png" width="200px" alt="msultra/encoder">
    <br>
    Encoder - Part of MsUltra
</h1>

<h4 align="center">Library that implements binary encoding and decoding, and some other utilities.</h4>

<p align="center">
    <img src="https://img.shields.io/github/go-mod/go-version/msultra/encoder">
    <img src="https://github.com/msultra/encoder/actions/workflows/test.yml/badge.svg">
    <a href="https://goreportcard.com/report/github.com/msultra/encoder"><img src="https://goreportcard.com/badge/msultra/encoder"></a>
    <a href="https://pkg.go.dev/github.com/msultra/encoder"><img src="https://pkg.go.dev/badge/github.com/msultra/encoder.svg"></a>
</p>

---

This library is used to encode and decode binary data. Basic usage can be summarized as:

```go
struct A {
    Field1 [4]byte
    Field2 uint32
}

a := A{
    Field1: [4]byte{1, 2, 3, 4},
    Field2: 1234,
}

bs, err := encoder.Marshal(a)
if err != nil {
    panic(err)
}

fmt.Println(bs) // Output: 01020304d2040000
```

