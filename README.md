# bulter

![gopher-butler-256](docs/img/gopher-butler-256.png)

a simple goroute pool manager

![workflow](https://www.plantuml.com/plantuml/proxy?fmt=svg&src=https://raw.githubusercontent.com/tangx/butler/main/docs/workflow.puml)


## Usage

[examples/main.go](__examples__/main.go)

#### get butler
```
go get -u github.com/tangx/butler
```

### demo code

```go
func demo() {
	javis := butler.Default()

	go func() {
		for {
			javis.AddJobs(func() {})
			time.Sleep(1 * time.Second)
		}
	}()

	javis.Work()
}
```

## todo

+ [x] context cancel 通知退出
+ [x] channel 分发信道
+ [x] signal 通知退出
+ [x] sync.WaitGroup 安全退出
+ [ ] errorgroup