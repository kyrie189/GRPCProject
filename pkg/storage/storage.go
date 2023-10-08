package storage

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
)


func Store(file *File) error {
	workDir, _ := os.Getwd() // 当前目录
	// 创建保存文件的目标路径

	dstPath := filepath.Join(workDir, "uploads",file.Name)   // 同名直接覆盖
	fmt.Println(dstPath)
	dst, err := os.Create(dstPath)  // 创建一个文件
	if err != nil {
		return err
	}
	defer dst.Close()
	if _,err := io.Copy(dst, file.buffer); err != nil {
		return err // 文件复制失败
	}
	return nil
}
