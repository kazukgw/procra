package main

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/kazukgw/procra"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "procra"

	app.Commands = []cli.Command{
		{
			Name:  "migrate",
			Usage: "migrate mysql tables",
			Subcommands: []cli.Command{
				{
					Name:   "create",
					Usage:  "create tables",
					Flags:  migrateFlags,
					Action: ActionMigrateCreate,
				},
				{
					Name:   "drop",
					Usage:  "drop tables",
					Flags:  migrateFlags,
					Action: ActionMigrateDrop,
				},
			},
		},
	}
	app.OnUsageError = func(c *cli.Context, err error, isSubcommand bool) error {
		fmt.Fprintf(c.App.Writer, "WRONG: %#v\n", err)
		return nil
	}
	app.Run(os.Args)
}

var migrateFlags []cli.Flag = []cli.Flag{
	cli.StringFlag{Name: "host", Value: "127.0.0.1"},
	cli.StringFlag{Name: "port", Value: "3306"},
	cli.StringFlag{Name: "user", Value: "root"},
	cli.StringFlag{Name: "password", Value: "password"},
	cli.StringFlag{Name: "database", Value: "procra"},
}

func dbConnStrFromCtx(ctx *cli.Context) string {
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s",
		ctx.String("user"),
		ctx.String("password"),
		ctx.String("host"),
		ctx.String("port"),
		ctx.String("database"),
	)
}

func ActionMigrateCreate(ctx *cli.Context) error {
	connstr := dbConnStrFromCtx(ctx)
	db, err := gorm.Open("mysql", connstr)
	if err != nil {
		panic(err.Error())
	}
	db.CreateTable(&procra.TargetURL{})
	db.CreateTable(&procra.TargetURLStats{})
	db.CreateTable(&procra.Attempt{})
	return nil
}

func ActionMigrateDrop(ctx *cli.Context) error {
	connstr := dbConnStrFromCtx(ctx)
	db, err := gorm.Open("mysql", connstr)
	if err != nil {
		panic(err.Error())
	}
	db.DropTable(&procra.TargetURL{})
	db.DropTable(&procra.TargetURLStats{})
	db.DropTable(&procra.Attempt{})
	return nil
}
