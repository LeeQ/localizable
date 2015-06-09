// readFile project readFile.go
package readFile

import (
	"bufio"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"strings"
)

const (
	SEM  = "\""
	CRLF = "\n"
)

//国际化的xml文件中的每个string的结构
type String struct {
	XMLName xml.Name `xml:"string"`
	Name    string   `xml:"name,attr"`
	Content string   `xml:",chardata"`
}

type Resource struct {
	XMLName   xml.Name `xml:"resources"`
	Resources []String `xml:"string"`
}

// 读取文件的函数调用大多数都需要检查错误
func check(e error) {
	if e != nil {
		fmt.Println(e.Error())
	}
}

//读取苹果系统的国际化文件，将其转为xml格式的文件
func ReadIOS(srcFilePath string, desFilePath string) {
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		fmt.Println("Failed to open the input file: ", srcFilePath)

	}
	defer srcFile.Close()

	srcReader := bufio.NewReader(srcFile)

	v := &Resource{} //生成的xml文件流
	for {
		str, err := srcReader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				k, val := getKeyVal(strings.TrimSpace(str)) //获取健值
				v.Resources = append(v.Resources, String{Name: k, Content: val})
				break
			} else {
				check(err)
			}
		}
		//条件判断
		if len(str) == 0 || !strings.HasPrefix(str, "\"") {
			continue
		}
		key, value := getKeyVal(strings.TrimSpace(str)) //获取健值
		v.Resources = append(v.Resources, String{Name: key, Content: value})
	}
	output, errors := xml.MarshalIndent(v, " ", "  ")
	check(errors)

	desFile, errs := os.Create(desFilePath)
	defer desFile.Close()
	if errs != nil {
		fmt.Println("Failed to create the output file: ", desFile)
	}

	desFile.WriteString(xml.Header)
	desFile.Write(output)
}

//提取每一行的key value值
func getKeyVal(content string) (string, string) {
	keyValues := strings.Split(content, "=")
	key := strings.TrimSpace(keyValues[0])
	key = strings.Trim(key, "\"")
	temp := strings.TrimSpace(keyValues[1])
	temp = strings.TrimLeft(temp, SEM)
	value := strings.TrimRight(temp, "\";")
	return key, value
}

//解析xml文件并解析为ios的国际化文件格式
func ParseXML(srcFilePath string, desFilePath string) {
	content, error := ioutil.ReadFile(srcFilePath)
	check(error)
	resources := Resource{}
	err := xml.Unmarshal(content, &resources)
	check(err)
	var result string
	for _, content := range resources.Resources {
		result = result + SEM + content.Name + SEM + "=" + SEM + content.Content + SEM + ";" + CRLF
	}

	//生成目标文件
	desFile, errs := os.Create(desFilePath)
	defer desFile.Close()
	if errs != nil {
		fmt.Println("Failed to create the output file: ", desFile)
		return
	}
	desFile.WriteString(result)
}
