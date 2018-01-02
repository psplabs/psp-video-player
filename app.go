package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"time"

	"github.com/djherbis/times"
	"github.com/gobuffalo/packr"
)

func main() {
	isExpired() // self destruct if video has expired.

	box := packr.NewBox("./static") // packr compiles all these files and bundles them up into a single executable
	//media := packr.NewBox("./media") // packr compiles all these files and bundles them up into a single executable

	http.Handle("/", http.FileServer(box)) // serve the index.html inside the ./app directory

	if !box.Has("media") {
		log.Println("Missing /static/media, serving /media from filesystem")
		http.Handle("/media/", http.StripPrefix("/media/", http.FileServer(http.Dir("./media")))) // serve the index.html inside the ./app directory
	}

	http.HandleFunc("/download/", download) // serve the index.html inside the ./app directory

	http.HandleFunc("/exit", exit) // index.html will call this route when the tab is closed, so the executable doesn't remain running.

	log.Println("Server ready, opening browser...")

	go open("http://localhost:3000/")        // open the browser window automatically
	panic(http.ListenAndServe(":3000", nil)) // start the server
}

func download(w http.ResponseWriter, r *http.Request) {
	log.Println("Copying...")
	CopyDir("./tmp", "./static")
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

// exists returns whether the given file or directory exists or not
func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

// open opens the specified URL in the default browser of the user.
func open(url string) error {
	var cmd string
	var args []string

	switch runtime.GOOS {
	case "windows":
		cmd = "cmd"
		args = []string{"/c", "start"}
	case "darwin":
		cmd = "open"
	default: // "linux", "freebsd", "openbsd", "netbsd"
		cmd = "xdg-open"
	}
	args = append(args, url)
	return exec.Command(cmd, args...).Start()
}

// checks to see if this file is expired, and self destructs if it is.
func isExpired() {
	t, err := times.Stat(os.Args[0])
	if err != nil {
		log.Fatal(err.Error())
	}

	expiresOn := t.ChangeTime().Add(time.Minute * time.Duration(1440*2)) // Expires after 48 hours
	isExpired := time.Now().After(expiresOn) || time.Now().Before(t.ChangeTime())

	if isExpired {
		selfDestruct()
	}
}

func selfDestruct() {
	log.Println("ERROR: Video expired! Please generate a new one at https://perspectives.org!")

	err := os.Remove(os.Args[0])
	if err != nil {
		panic(err)
	} else {
		os.Exit(0)
	}
}

// CopyFile copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file. The file mode will be copied from the source and
// the copied data is synced/flushed to stable storage.
func CopyFile(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		if e := out.Close(); e != nil {
			err = e
		}
	}()

	_, err = io.Copy(out, in)
	if err != nil {
		return
	}

	err = out.Sync()
	if err != nil {
		return
	}

	si, err := os.Stat(src)
	if err != nil {
		return
	}
	err = os.Chmod(dst, si.Mode())
	if err != nil {
		return
	}

	return
}

// CopyDir recursively copies a directory tree, attempting to preserve permissions.
// Source directory must exist, destination directory must *not* exist.
// Symlinks are ignored and skipped.
func CopyDir(src string, dst string) (err error) {
	src = filepath.Clean(src)
	dst = filepath.Clean(dst)

	si, err := os.Stat(src)
	if err != nil {
		return err
	}
	if !si.IsDir() {
		return fmt.Errorf("source is not a directory")
	}

	_, err = os.Stat(dst)
	if err != nil && !os.IsNotExist(err) {
		return
	}
	if err == nil {
		return fmt.Errorf("destination already exists")
	}

	err = os.MkdirAll(dst, si.Mode())
	if err != nil {
		return
	}

	entries, err := ioutil.ReadDir(src)
	if err != nil {
		return
	}

	for _, entry := range entries {
		srcPath := filepath.Join(src, entry.Name())
		dstPath := filepath.Join(dst, entry.Name())

		if entry.IsDir() {
			err = CopyDir(srcPath, dstPath)
			if err != nil {
				return
			}
		} else {
			// Skip symlinks.
			if entry.Mode()&os.ModeSymlink != 0 {
				continue
			}

			err = CopyFile(srcPath, dstPath)
			if err != nil {
				return
			}
		}
	}

	return
}

func exit(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}
