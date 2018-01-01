package main

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"time"

	"github.com/djherbis/times"
	"github.com/gobuffalo/packr"
)

func main() {
	isExpired() // self destruct if video has expired.

	box := packr.NewBox("./") // packr compiles all these files and bundles them up into a single executable

	http.Handle("/", http.FileServer(box)) // serve the index.html inside the ./app directory
	http.HandleFunc("/exit", exit)         // index.html will call this route when the tab is closed, so the executable doesn't remain running.

	log.Println("Server ready, opening browser...")

	go open("http://localhost:3000/")        // open the browser window automatically
	panic(http.ListenAndServe(":3000", nil)) // start the server
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

func exit(w http.ResponseWriter, r *http.Request) {
	os.Exit(0)
}
