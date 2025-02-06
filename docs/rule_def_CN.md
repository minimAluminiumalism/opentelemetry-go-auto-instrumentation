# 规则定义的字段
## 插桩函数

ImportPath: 包含要插桩函数的包的导入路径。例如 net/http。
Function: 要插桩的函数名称。
ReceiverType: 要插桩函数的接收器类型。例如 *http.Server、http.Server。
OnEnter: 在插桩函数被调用时要调用的函数名称。例如 clientOnEnter。
OnExit: 在插桩函数返回时要调用的函数名称。例如 clientOnExit。
Order: 插桩函数中探针代码的顺序。例如 0、1、2。
Path: 包含探针代码的目录路径。路径可以是 go 模块 URL 或本地文件系统路径，例如 github.com/foo/bar 或 /path/to/probe/code。
Version: 包含要插桩函数的包的版本。例如 [1.0.0,1.1.0)，版本范围是 [1.0.0,1.1.0)，表示版本大于等于 1.0.0 且小于 1.1.0。

## 在编译包期间添加新文件
ImportPath: 包含要插桩函数的包的导入路径。
FileName: 要添加的文件名。
Path: 包含探针代码的目录路径。


## 向结构体添加新字段
ImportPath: 包含要插桩结构体的包的导入路径。
StructType: 要插桩的结构体名称。
FieldName: 要添加的字段名称。
FieldType: 要添加的字段类型。