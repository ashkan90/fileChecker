package main

import (
	"io/ioutil"
	"os"
	"reflect"
	"testing"
)

func load(path string) []os.FileInfo {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	return files
}

func look(where []os.FileInfo, i func(os.FileInfo)) {
	for _, info := range where {
		i(info)
	}
}

func compare(s interface{}, t interface{}) bool {
	sv := reflect.ValueOf(s)
	tv := reflect.ValueOf(t)

	if (sv.IsValid() && tv.IsValid()) && sv.Len() == tv.Len() {
		return reflect.DeepEqual(sv.Interface(), tv.Interface())
	}

	return false
}

func TestCopy(t *testing.T) {

	source := `C:\Users\IO\Desktop\express_rtc`
	destination := `C:\Users\IO\Desktop\storage`

	sourceFiles := load(source) // files of express_rtc

	var wanted []os.FileInfo

	err := Copy(source, destination)
	if err != nil {
		t.Fatal(err)
	}

	readDirectoryFn(destination, func(dir os.FileInfo) {
		look(sourceFiles, func(info os.FileInfo) {
			if dir.Name() == info.Name() {
				wanted = append(wanted, info)
			}
		})
	})

	if compare(wanted, sourceFiles) {
		t.Skipped()
	}

}
