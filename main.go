package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
)

func readDirectoryFn(path string, c func(dir os.FileInfo)) {

	files, err := ioutil.ReadDir(path)
	if err != nil {
		panic(err)
	}

	for _, dir := range files {
		c(dir)
	}

}

var escapeDir = flag.String("escape", "evraklar", "Root klasör içerisinde işleme alınmayacak klasör adı. Bkz: 'evrak'")

const dirPermission = os.FileMode(0755)

// copy dispatches copy-funcs according to the mode.
// Because this "copy" could be called recursively,
// "info" MUST be given here, NOT nil.
func copy(src, dest string, info os.FileInfo) error {
	if info.Mode()&os.ModeSymlink != 0 {
		return lcopy(src, dest, info)
	}
	if info.IsDir() {
		return dcopy(src, dest, info)
	}
	return fcopy(src, dest, info)
}

// fcopy, dosya kopyalama için kullanılıyor.
// Dosya açarak yetkileri öngörebiliyor.
func fcopy(src, dest string, info os.FileInfo) error {

	if err := os.MkdirAll(filepath.Dir(dest), os.ModePerm); err != nil {
		return err
	}

	f, err := os.Create(dest)
	if err != nil {
		return err
	}
	defer f.Close()

	if err = os.Chmod(f.Name(), info.Mode()); err != nil {
		return err
	}

	s, err := os.Open(src)
	if err != nil {
		return err
	}
	defer s.Close()

	_, err = io.Copy(f, s)
	return err
}

// dcopy, sadece klasör kopyalama işlemi yapar.
func dcopy(srcdir, destdir string, info os.FileInfo) error {
	if exists(srcdir) {


		originalMode := info.Mode()

		// Make dest dir with 0755 so that everything writable.
		if err := os.MkdirAll(destdir, dirPermission); err != nil {
			return err
		}
		// Recover dir mode with original one.
		defer os.Chmod(destdir, originalMode)

		contents, err := ioutil.ReadDir(srcdir)
		if err != nil {
			return err
		}

		err = removeAll(srcdir)
		if err != nil {
			return err
		}


		for _, content := range contents {
			cs, cd := filepath.Join(srcdir, content.Name()), filepath.Join(destdir, content.Name())

			if err := copy(cs, cd, content); err != nil {
				// If any error, exit immediately
				return err
			}
		}
	}


	return nil
}

// lcopy is for a symlink,
// with just creating a new symlink by replicating src symlink.
func lcopy(src, dest string, info os.FileInfo) error {
	src, err := os.Readlink(src)
	if err != nil {
		return err
	}
	return os.Symlink(src, dest)
}

func removeAll(src string) error {
	if _, err := os.Lstat(src); err != nil {
		return err
	}
	err := os.RemoveAll(src)
	if err != nil {
		return err
	}
	return nil
}

func Copy(src, dest string) error {
	info, err := os.Lstat(src)
	if err != nil {
		return err
	}
	return copy(src, dest, info)
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
func main()  {
	flag.Parse()
	path := "C:\\Users\\user\\Desktop\\storage"

	mapDirectory(path, path + "\\" + *escapeDir)
	fmt.Println("Dosyalar başarıyla kopyalandı ve silindi.")

	exit()

}
