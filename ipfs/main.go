package ipfs

import (
	"fmt"
	shell "github.com/ipfs/go-ipfs-api"
	"github.com/joho/godotenv"
	"log"
	"os"
)

const (
	ServerAddr = "localhost"
	ServerPort = "5001"
)

var Api *shell.Shell

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	Api = shell.NewShell(fmt.Sprintf("%s:%s", os.Getenv("ipfs_addr"), os.Getenv("ipfs_port")))
}
