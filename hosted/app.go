package main

import (
	"net/http"
	"os/exec"
	"runtime"
)

func main() {
	//box := packr.NewBox("./") // packr compiles all these files and bundles them up into a single executable

	http.Handle("/", http.FileServer(http.Dir("./"))) // serve the index.html inside the ./app directory

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
