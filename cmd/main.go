package main

import (
	"fmt"
	"os"

	"github.com/arfan21/vocagame/cmd/api"
	migration "github.com/arfan21/vocagame/cmd/migrate"
	"github.com/arfan21/vocagame/config"
	"github.com/urfave/cli/v2"
)

// @title vocagame
// @version 1.0
// @description This is a sample server cell for vocagame.
// @termsOfService http://swagger.io/terms/
// @contact.name API Support
// @contact.url http://www.synapsis.id
// @contact.email
// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html
// @host localhost:8080
// @BasePath /
func main() {
	_, err := config.LoadConfig()
	if err != nil {
		fmt.Println("error load config", err)
		panic(err)
	}

	_, err = config.ParseConfig(config.GetViper())
	if err != nil {
		fmt.Println("error parse config", err)
		panic(err)
	}

	appCli := cli.NewApp()
	appCli.Name = config.Get().Service.Name
	appCli.Commands = []*cli.Command{
		migration.Root(),
		api.Serve(),
	}

	if err := appCli.Run(os.Args); err != nil {
		panic(err)
	}
}
