package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// parsePomDependency 解析 a.txt 文件中的 dependency 部分，返回包含 groupId, artifactId, version 的 map 和 error
func parsePomDependency(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("open file err: %w", err)
	}
	defer file.Close()

	dependencyMap := make(map[string]string)
	var currentKey string

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		switch {
		case strings.Contains(line, "<groupId>"):
			currentKey = "groupId"
			dependencyMap[currentKey] = extractValue(line)
		case strings.Contains(line, "<artifactId>"):
			currentKey = "artifactId"
			dependencyMap[currentKey] = extractValue(line)
		case strings.Contains(line, "<version>"):
			currentKey = "version"
			dependencyMap[currentKey] = extractValue(line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("read file err: %w", err)
	}

	// 检查是否所有的 key 都存在
	expectedKeys := []string{"groupId", "artifactId", "version"}
	for _, key := range expectedKeys {
		if _, exists := dependencyMap[key]; !exists {
			return nil, fmt.Errorf("missing required key: %s", key)
		}
	}

	return dependencyMap, nil
}

// extractValue 从XML标签中提取值
func extractValue(line string) string {
	start := strings.Index(line, ">") + 1
	end := strings.Index(line[start:], "<")
	if end == -1 {
		return line[start:]
	}
	return line[start : start+end]
}

// findUniqueFile 在当前目录下查找唯一的指定后缀文件
func findUniqueFile(pattern string) (string, error) {
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to scan directory: %v", err)
	}

	switch len(matches) {
	case 0:
		return "", fmt.Errorf("no %s files found in the current directory", pattern)
	case 1:
		return matches[0], nil
	default:
		return "", fmt.Errorf("multiple %s files found in the current directory: %v", pattern, matches)
	}
}

// installMavenDependency 使用 mvn install:install-file 命令安装 JAR 文件
func installMavenDependency(dependencyMap map[string]string, jarFilePath string) error {
	groupId := dependencyMap["groupId"]
	artifactId := dependencyMap["artifactId"]
	version := dependencyMap["version"]

	// 构建 Maven 命令
	cmd := exec.Command("mvn", "install:install-file",
		"-Dfile="+jarFilePath,
		"-DgroupId="+groupId,
		"-DartifactId="+artifactId,
		"-Dversion="+version,
		"-Dpackaging=jar")

	// 打印命令以便调试
	fmt.Println("Executing command:", cmd.String())

	// 执行命令并捕获输出
	output, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("failed to execute maven command: %v\nOutput: %s", err, string(output))
	}

	// 打印命令输出
	fmt.Println(string(output))

	return nil
}

func main() {
	// 解析命令行参数
	txtFlag := flag.String("txt", "", "The path to the TXT file containing GAV coordinates (optional)")
	jarFlag := flag.String("jar", "", "The path to the JAR file (optional)")
	flag.Parse()

	var txtFile, jarFile string
	var err error

	// 处理 TXT 文件路径
	if *txtFlag != "" {
		txtFile = *txtFlag
	} else {
		// 否则在当前目录下查找唯一的 .txt 文件
		txtFile, err = findUniqueFile("*.txt")
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// 处理 JAR 文件路径
	if *jarFlag != "" {
		jarFile = *jarFlag
	} else {
		// 否则在当前目录下查找唯一的 .jar 文件
		jarFile, err = findUniqueFile("*.jar")
		if err != nil {
			fmt.Println(err)
			return
		}
	}

	// 检查文件是否存在
	if _, err := os.Stat(txtFile); os.IsNotExist(err) {
		fmt.Println("TXT file not found:", txtFile)
		return
	}

	if _, err := os.Stat(jarFile); os.IsNotExist(err) {
		fmt.Println("JAR file not found:", jarFile)
		return
	}

	// 解析 a.txt 文件中的依赖信息
	dependencyMap, err := parsePomDependency(txtFile)
	if err != nil {
		fmt.Println("Error parsing pom file:", err)
		return
	}

	// 安装 JAR 文件到本地 Maven 仓库
	err = installMavenDependency(dependencyMap, jarFile)
	if err != nil {
		fmt.Println("Error installing JAR file:", err)
		return
	}

	fmt.Println("JAR file installed successfully.")
}
