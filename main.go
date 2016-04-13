package main

import (
	"bytes"
	"fmt"
	"github.com/drone/drone-plugin-go/plugin"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

type Config struct {
	Package string `json:"package"`
	Output  string `json:"output"`
}

var (
	buildDate string
)

func main() {
	fmt.Printf("Drone Go Compiler Plugin built at %s\n", buildDate)

	workspace := plugin.Workspace{}
	vargs := Config{}

	plugin.Param("workspace", &workspace)
	plugin.Param("vargs", &vargs)
	plugin.MustParse()

	if len(vargs.Output) == 0 {
		vargs.Output = "."
	}

	path := filepath.Join(workspace.Path, vargs.Package)
	output := filepath.Join(workspace.Path, vargs.Output, vargs.Package)

	var cmd *exec.Cmd
	var out bytes.Buffer
	var err error

	cmd = exec.Command("go", "version")
	cmd.Stdout = &out
	cmd.Stderr = os.Stderr
	trace(cmd)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}

	fmt.Printf("Building with %s\n", out.String())

	cmd = exec.Command("godep", "go", "build", "-o", output)
	cmd.Env = append(cmd.Env, "PATH=/usr/local/go/bin:$PATH")
	cmd.Env = append(cmd.Env, "GOPATH=/drone")
	cmd.Env = append(cmd.Env, "CGO_ENABLED=0")
	cmd.Env = append(cmd.Env, "GO15VENDOREXPERIMENT=0")
	cmd.Env = append(cmd.Env, "LDFLAGS='-d -w -s'")
	cmd.Dir = path
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	trace(cmd)
	err = cmd.Run()
	if err != nil {
		fmt.Println("Error: ", err.Error())
		os.Exit(1)
	}
}

// Trace writes each command to standard error (preceded by a ‘$ ’) before it
// is executed. Used for debugging your build.
func trace(cmd *exec.Cmd) {
	fmt.Println("$", strings.Join(cmd.Args, " "))
}
