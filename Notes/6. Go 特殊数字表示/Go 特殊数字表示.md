1. 二进制
```go
num := 0b10 // 4
```

2. 八进制
```go
num := 0o10 // 8
```

3. 十六进制
```go
num := 0x10 // 16
```

4. 科学计数法
```go
num := 1e3 // 1000
```

只有十进制可以用科学计数法表示

5. 虚数
```go
num := 0o123i
num := 0x1e3i
```

末尾加上i即可

6. 分隔
```
num := 10_000
```

分隔符不可用于开头和结尾