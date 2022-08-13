package compress

import (
	"io/fs"
	"os"
	"path/filepath"
	"testing"
)

func TestCompress(t *testing.T) {
	// 在当前目录创建一个临时文件夹、文件夹下创建两个临时文件
	fileDir, err := os.MkdirTemp("", "files*")
	zipDir, err := os.MkdirTemp("", "zip*")
	_, err = os.Create(filepath.Join(fileDir, "temp1.txt"))
	_, err = os.Create(filepath.Join(fileDir, "temp2.txt"))
	// 调用压缩函数
	zipFilePath := filepath.Join(zipDir, "temp.zip")
	cw, err := os.Create(zipFilePath)
	err = Compress(cw, fileDir)
	if err != nil {
		t.Error("压缩时出现错误", err)
	}
	// 预期得到一个压缩文件
	fi, err := cw.Stat()
	if fi.Size() == 0 {
		t.Errorf("压缩后的内容不正确")
	}
	// 读取压缩文件的内容，需要完全符合预期
	extractFileDir, err := os.MkdirTemp("", "extractFiles*")
	file, err := os.Open(zipFilePath)
	if err != nil {
		t.Error("读取压缩文件错误", err)
	}
	var bytes []byte
	_, err = file.Read(bytes)
	err = DeCompress(extractFileDir, bytes)
	var k int
	err = filepath.WalkDir(extractFileDir, func(path string, d fs.DirEntry, err error) error {
		k++
		return nil
	})
	if k != 3 {
		t.Errorf("解压缩后的文件数量不正确")
	}
	defer func() {
		err = os.RemoveAll(fileDir)
		err = os.RemoveAll(zipDir)
		err = os.RemoveAll(extractFileDir)
	}()
	if err != nil {
		t.Error("测试出现问题", err)
	}
}
