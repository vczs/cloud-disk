package helper

import (
	"crypto/md5"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"strconv"
)

const partSize = 3 * 1024 * 1024 // 3M

// 文件分块
func GeneratePart() {
	fileInfo, err := os.Stat("music.mp3")
	if err != nil {
		fmt.Println(err)
	}
	// 分块个数
	partNum := math.Ceil(float64(fileInfo.Size()) / float64(partSize))
	myFile, err := os.OpenFile("music.mp3", os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	b := make([]byte, partSize)
	for i := 0; i < int(partNum); i++ {
		// 指定读取文件的起始位置
		myFile.Seek(int64(i*partSize), 0)
		if partSize > fileInfo.Size()-int64(i*partSize) {
			b = make([]byte, fileInfo.Size()-int64(i*partSize))
		}
		myFile.Read(b)
		f, err := os.OpenFile("./"+strconv.Itoa(i+1)+".part", os.O_CREATE|os.O_WRONLY, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
		f.Write(b)
		f.Close()
	}
	myFile.Close()
}

// 合并分块
func MergePart() {
	myFile, err := os.OpenFile("merage.mp3", os.O_CREATE|os.O_WRONLY|os.O_APPEND, os.ModePerm)
	if err != nil {
		fmt.Println(err)
	}
	fileInfo, err := os.Stat("music.mp3")
	if err != nil {
		fmt.Println(err)
	}
	// 分片的个数
	partNum := math.Ceil(float64(fileInfo.Size()) / float64(partSize))
	for i := 0; i < int(partNum); i++ {
		f, err := os.OpenFile("./"+strconv.Itoa(i+1)+".part", os.O_RDONLY, os.ModePerm)
		if err != nil {
			fmt.Println(err)
		}
		b, err := ioutil.ReadAll(f)
		if err != nil {
			fmt.Println(err)
		}
		myFile.Write(b)
		f.Close()
	}
	myFile.Close()
}

// 文件一致性校验
func TestCheck() {
	// 获取第一个文件的信息
	file1, err := os.OpenFile("music.mp3", os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	b1, err := ioutil.ReadAll(file1)
	if err != nil {
		fmt.Println(err)
	}
	// 获取第二个文件的信息
	file2, err := os.OpenFile("merage.mp3", os.O_RDONLY, 0666)
	if err != nil {
		fmt.Println(err)
	}
	b2, err := ioutil.ReadAll(file2)
	if err != nil {
		fmt.Println(err)
	}
	s1 := fmt.Sprintf("%x", md5.Sum(b1))
	s2 := fmt.Sprintf("%x", md5.Sum(b2))
	fmt.Println(s1, s2)
}
