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
