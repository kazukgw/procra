package main

import (
	"fmt"
	"os"

	"bitbucket.org/monotaro/monotaro_crawler2"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"gopkg.in/urfave/cli.v1"
)

func main() {
	app := cli.NewApp()
	app.Name = "monocra2"

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
				{
					Name:   "example",
					Usage:  "insert example data",
					Flags:  migrateFlags,
					Action: ActionMigrateExampleData,
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
	cli.StringFlag{Name: "database", Value: "monocra2"},
}

func appendFlag(flags []cli.Flag, flag cli.Flag) []cli.Flag {
	return append(flags, flag)
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
	db.CreateTable(&monocra2.TargetURL{})
	db.CreateTable(&monocra2.TargetURLStats{})
	db.CreateTable(&monocra2.Attempt{})
	return nil
}

func ActionMigrateDrop(ctx *cli.Context) error {
	connstr := dbConnStrFromCtx(ctx)
	db, err := gorm.Open("mysql", connstr)
	if err != nil {
		panic(err.Error())
	}
	db.DropTable(&monocra2.TargetURL{})
	db.DropTable(&monocra2.TargetURLStats{})
	db.DropTable(&monocra2.Attempt{})
	return nil
}

func ActionMigrateExampleData(ctx *cli.Context) error {
	connstr := dbConnStrFromCtx(ctx)
	db, err := gorm.Open("mysql", connstr)
	if err != nil {
		panic(err.Error())
	}

	targs := []string{
		"http://jp.misumi-ec.com/vona2/detail/110300488540/?Tab=CodeList",
		"http://jp.misumi-ec.com/vona2/detail/110300489240/?Tab=codeList",
		"http://jp.misumi-ec.com/vona2/detail/110302376530/?Tab=codeList",
		"http://jp.misumi-ec.com/vona2/detail/223004966754/?Tab=codeList",
		"http://jp.misumi-ec.com/vona2/detail/110100265220/?Tab=codeList",
		"http://jp.misumi-ec.com/vona2/detail/222000373278/?Tab=codeList",
		"http://jp.misumi-ec.com/vona2/detail/110400161750/?Tab=codeList",
	}

	for _, t := range targs {
		targ, _ := monocra2.NewTargetURLFromRawURL(t)
		db.Create(targ)
		if db.Error != nil {
			panic(db.Error.Error())
		}
	}
	return nil
}
