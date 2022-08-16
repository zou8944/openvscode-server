package compress

import (
	"archive/zip"
	"bytes"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
	"strings"
)

// Compress 将文件夹或文件压缩，写入到输出流
func Compress(dst io.Writer, src string) error {
	zw := zip.NewWriter(dst)
	defer zw.Close()
	return filepath.Walk(src, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}
		// 不写入当前层
		if path == src {
			return nil
		}
		fh, err := zip.FileInfoHeader(info)
		if err != nil {
			return err
		}
		// 文件路径不包含根目录
		fh.Name = strings.TrimPrefix(path, src)
		fh.Name = strings.TrimPrefix(fh.Name, string(filepath.Separator))
		if info.IsDir() {
			fh.Name += "/"
		}
		fw, err := zw.CreateHeader(fh)
		if err != nil {
			return err
		}
		// 如果是文件夹就只写入头信息
		if !fh.Mode().IsRegular() {
			return nil
		}
		fr, err := os.Open(path)
		if err != nil {
			return err
		}
		n, err := io.Copy(fw, fr)
		if err != nil {
			return err
		}
		log.Printf("成功压缩文件： %s, 共写入了 %d 个字符的数据\n", path, n)
		return nil
	})
}

func DeCompress(dstDir string, b []byte) error {
	readAt := bytes.NewReader(b)
	reader, err := zip.NewReader(readAt, int64(len(b)))
	if err != nil {
		return err
	}
	return deCompress(dstDir, reader)
}

func deCompress(dstDir string, reader *zip.Reader) error {
	// 删除原本的数据再创建
	err := os.RemoveAll(dstDir)
	err = os.MkdirAll(dstDir, os.ModePerm)
	if err != nil {
		return err
	}
	// 遍历zip reader中的文件，依次写入目标文件夹
	for _, file := range reader.File {
		path := filepath.Join(dstDir, file.Name)
		// 如果是目录就创建目录
		if file.FileInfo().IsDir() {
			if err := os.Mkdir(path, file.Mode()); err != nil {
				return err
			}
			log.Printf("成功解压目录: %s", path)
			continue
		}
		// 如果是文件则创建文件
		fr, err := file.Open()
		if err != nil {
			return err
		}
		fw, err := os.OpenFile(path, os.O_CREATE|os.O_WRONLY, file.Mode())
		if err != nil {
			return err
		}
		_, err = io.Copy(fw, fr)
		if err != nil {
			return err
		}
		fr.Close()
		fw.Close()
		log.Printf("成功解压文件: %s", path)
	}
	return nil
}
