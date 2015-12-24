### go-config

更轻松的读取配置文件

### 用法

示例配置文件test.conf内容如下：

> # this is comment line 1
> ; this is comment line 2
> USERNAME = jiazhoulvke
> PORT     = 1984
> VERSION  = 1.1
> HOST     = localhost

读取配置文件:

```go
    cfg,err:=goconfig.Parse("test.conf")
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println(cfg.GetInt("PORT")) //1984
    fmt.Println(cfg.GetString("USERNAME")) //jiazhoulvke
```


获取字段方法一:

```go
    username, err := cfg.GetString("USERNAME")
    if err != nil {
        panic(err)
    }
    fmt.Println(username) //jiazhoulvke
```

获取字段方法二:

```go
    username := cfg.StringDefault("NOKEY", "jiazhoulvke") //如果没有找到对应的字段则使用默认值
    fmt.Println(username) //jiazhoulvke
```


设置字段,支持代码链:

```go
    cfg.Set("Url", "http://www.jiazhoulvke.com").Set("Tags", "vim,python,linux,go")
```


定义一个struct，利用goconfig初始化:

```go
    type MyConfig struct {
        Username `cfgname:"USERNAME"` //cfgname表示在配置文件中的实际字段名称
        Password `default:"hello,world"` //default表示当未在配置文件中找到对应的字段时赋予变量的默认值
    }

    var myconfig MyConfig
    cfg.Init(&myconfig) //注意是传指针
    fmt.Println(myconfig.UserName) //输出jiazhoulvke hello,world
```

