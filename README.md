# mvn-repo-helper
帮助用maven install from file命令安装jar包

## 用法
坐标写在mvn-txt里，确保有<groupId>、<artifactId>、<version>，出现在单独的行，参考示例写法

jar包放到程序运行目录。

使用 flag.String 定义两个可选的命令行参数：

--txt：用于指定包含 GAV 坐标的 .txt 文件路径。

--jar：用于指定 JAR 文件的路径。

如果用户提供了 --txt 参数，则直接使用该参数提供的文件路径；否则调用 findUniqueFile 函数在当前目录下查找唯一的 .txt 文件。

如果用户提供了 --jar 参数，则直接使用该参数提供的文件路径；否则调用 findUniqueFile 函数在当前目录下查找唯一的 .jar 文件。

go run main.go --txt=a.txt --jar=b.jar

go run main.go --txt=a.txt

go run main.go --txt=a.txt

如果程序运行目录只有唯一的jar包和唯一的.txt文件（txt文件里是gav坐标信息）：

go run main.go

## windows下路径怎么写
txtFile := `E:\abc\def\xyz.txt`

txtFile := "E:\\abc\\def\\xyz.txt"

或者：

filepath.Join 是 Go 标准库中的一个函数，它可以跨平台地构建路径。它会根据当前操作系统自动选择正确的路径分隔符（Windows 使用 \，Unix/Linux/macOS 使用 /）。虽然 filepath.Join 通常用于拼接相对路径，但它也可以用于绝对路径。

import "path/filepath"

// 构建路径
txtFile := filepath.Join("E:", "abc", "def", "xyz.txt")

### 启动方式（windows下）
go run main.go --txt="E:\abc\def\xyz.txt" --jar="E:\abc\def\b.jar"

go run main.go --txt="E:\abc\def\xyz.txt"


