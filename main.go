package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
)

var escapeDir = flag.String("escape", "evraklar", "Root klasör içerisinde işleme alınmayacak klasör adı. Örn: 'evrak'")
var path = flag.String("path", "", "Kopyalama yapılacak ana dosya dizini. Örn:`C:\\Users\\IO\\Desktop\\storage`")

func readDirectoryFn(path string, c func(dir os.FileInfo)) {

	if exists(path) {

		files, err := ioutil.ReadDir(path)
		if err != nil {
			panic(err)
		}

		for _, dir := range files {
			c(dir)
		}
	}

}

func mapDirectory(path string, to string) {

	readDirectoryFn(path, func(dir os.FileInfo) {
		if dir.Name() != *escapeDir {
			src := path + "\\" + dir.Name()
			dest := to + "\\" + dir.Name()

			err := Copy(src, dest)
			if err != nil {
				panic(err)
			}
		}

	})
}

func exists(src string) bool {
	if _, err := os.Lstat(src); err != nil {
		return false
	}

	return true
}

func exit() {
	fmt.Print("Press 'Enter' to continue...")
	_, _ = bufio.NewReader(os.Stdin).ReadBytes('\n')
}
func main() {
	flag.Parse()

	mapDirectory(*path, *path+"\\"+*escapeDir)

	fmt.Println("Dosyalar başarıyla kopyalandı ve silindi.")

	exit()
}
