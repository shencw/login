package main

import (
	"github.com/shencw/login/internal/apiserver"
	"math/rand"
	"os"
	"runtime"
	"time"
)

func main() {
	rand.Seed(time.Now().UTC().UnixNano())
	if len(os.Getenv("GOMAXPROCS")) ==  0 {
		runtime.GOMAXPROCS(runtime.NumCPU())
	}
	apiserver.NewApp("login-apiserver").Run()
}
