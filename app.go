package main

import (
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

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

	go open("http://localhost:3000/")              // open the browser window automatically
	log.Println(http.ListenAndServe(":3000", nil)) // start the server
}

func download(w http.ResponseWriter, r *http.Request) {
	urlPart := strings.Split(r.URL.Path, "/")
	video := urlPart[2]
	platform := ""

	if len(urlPart) > 3 {
		platform = urlPart[3]
	}

	fmt.Printf("Downloading %s\n", video)

	expires := time.Now().Add(time.Minute * time.Duration(1440*2)).Unix() // 2 days
	var args []string
	args = []string{video, strconv.FormatInt(expires, 10), platform}
	out, err := exec.Command("./build", args...).Output()

	if err != nil {
		fmt.Printf("output is %s\n", err)
		log.Fatal(err)
	}

	fmt.Printf("%s", out)

	Filename := "psp-" + video + ".zip"
	fpath := "./dist/" + Filename

	//Check if file exists and open
	Openfile, err := os.Open(fpath)
	defer Openfile.Close() //Close after function return
	if err != nil {
		//File not found, send 404
		http.Error(w, "File not found: "+fpath, 404)
		return
	}

	//File is found, create and send the correct headers

	//Get the Content-Type of the file
	//Create a buffer to store the header of the file in
	FileHeader := make([]byte, 512)
	//Copy the headers into the FileHeader buffer
	Openfile.Read(FileHeader)
	//Get content type of file
	FileContentType := http.DetectContentType(FileHeader)

	//Get the file size
	FileStat, _ := Openfile.Stat()                     //Get info from file
	FileSize := strconv.FormatInt(FileStat.Size(), 10) //Get file size as a string

	//Send the headers
	w.Header().Set("Content-Disposition", "attachment; filename="+string(Filename))
	w.Header().Set("Content-Type", FileContentType)
	w.Header().Set("Content-Length", FileSize)

	//Send the file
	//We read 512 bytes from the file already so we reset the offset back to 0
	Openfile.Seek(0, 0)
	io.Copy(w, Openfile) //'Copy' the file to the client
	return

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
	s, err := ioutil.ReadFile("expires") // just pass the file name
	if err != nil {
		fmt.Print(err)

	} else {

		expires, err := strconv.ParseInt(strings.Replace(string(s), "\n", "", -1), 10, 64)

		if err != nil {
			panic(err)
		}

		expiresOn := time.Unix(expires, 0) // Expires after 48 hours
		isExpired := time.Now().After(expiresOn)

		if isExpired {
			selfDestruct()
		}
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

func exit(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}
