package main

import (
	"fmt"
	"io"
	"log"
	"net/http"
	"os"

	"github.com/thatisuday/commando"
)

func check(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func server(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	dir, err := flags["dir"].GetString()
	check(err)

	port, err := flags["port"].GetString()
	check(err)

	fmt.Printf("Starting server for directory %s on port %s\n\n", dir, port)
	fs := http.FileServer(http.Dir(dir))
	fmt.Println(http.ListenAndServe(":"+port, fs))

}

func download(args map[string]commando.ArgValue, flags map[string]commando.FlagValue) {
	url := args["url"].Value

	outfile, err := flags["outfile"].GetString()
	check(err)

	// Send the request
	resp, err := http.Get(url)
	check(err)
	defer resp.Body.Close()

	// Make the file
	out, err := os.Create(outfile)
	check(err)
	defer out.Close()

	_, err = io.Copy(out, resp.Body)
	check(err)
}

func main() {

	wd, err := os.Getwd()
	check(err)

	name, err := os.Executable()
	check(err)

	commando.
		SetExecutableName(name).
		SetVersion("0.0.1").
		SetDescription("This tool provides utilities to facilitate penetration testing on multiple architectures.")

	// File server
	commando.
		Register("server").
		SetShortDescription("starts a server").
		SetDescription("This command starts a server with the specified options.").
		AddFlag("dir,d", "directory to serve", commando.String, wd).
		AddFlag("port,p", "port to serve on", commando.String, "8080").
		SetAction(server)

	// Download file command
	commando.
		Register("download").
		SetShortDescription("downloads a file from a URL").
		SetDescription("This command downloads the file at a provided URL").
		AddArgument("url", "target URL", "").
		AddFlag("path,p", "output file", commando.String, wd+"/outfile").
		SetAction(download)

	commando.Parse(nil)

}
