package main

import (
	"covid19kalteng/covid19"
	"covid19kalteng/migration"
	"covid19kalteng/router"
	"database/sql"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/labstack/echo/middleware"
	"github.com/pressly/goose"
)

var (
	flags = flag.NewFlagSet("goose", flag.ExitOnError)
)

func main() {
	defer covid19.App.Close()

	flags.Usage = usage
	flags.Parse(os.Args[1:])
	args := flags.Args()

	migrationDir := "migration" // migration directory

	switch args[0] {
	default:
		flags.Usage()
		break
	case "run":
		e := router.NewRouter()
		if covid19.App.Config.GetBool("react_cors") {
			e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
				AllowOrigins: []string{"*"},
				AllowMethods: []string{"*"},
				AllowHeaders: []string{"*"},
			}))
		}
		e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
			Format: `{"time":"${time_rfc3339_nano}","id":"${id}","remote_ip":"${remote_ip}","host":"${host}",` +
				`"method":"${method}","uri":"${uri}","status":${status},"error":"${error}","latency":${latency},` +
				`"latency_human":"${latency_human}","bytes_in":${bytes_in},` +
				`"bytes_out":${bytes_out}}` + "\n",
		}))
		e.Logger.Fatal(e.Start(":" + covid19.App.Port))
		os.Exit(0)
		break
	case "seed":
		migration.Seed()
		os.Exit(0)
		break
	case "truncate":
		err := migration.Truncate(args[1:])
		if err != nil {
			log.Fatalf("%v", err)
			flags.Usage()
		}
		os.Exit(0)
		break
	case "create":
		if err := goose.Run("create", nil, migrationDir, args[1:]...); err != nil {
			log.Fatalf("goose create: %v", err)
			flags.Usage()
		}
		return
	case "migrate": // command example : [app name] migrate up
		if err := goose.SetDialect("postgres"); err != nil {
			log.Fatalf("goose set dialect : %v", err)
			flags.Usage()
		}

		dbconf := covid19.App.Config.GetStringMap(fmt.Sprintf("%s.database", covid19.App.ENV))
		connectionString := fmt.Sprintf("host=%s user=%s dbname=%s sslmode=%s password=%s",
			dbconf["host"].(string), dbconf["username"].(string), dbconf["table"].(string),
			dbconf["sslmode"].(string), dbconf["password"].(string),
		)

		db, err := sql.Open("postgres", connectionString)
		if err != nil {
			log.Fatalf("-connectionString=%q: %v\n", connectionString, err)
			flags.Usage()
		}

		if err := goose.Run(args[1], db, migrationDir, args[2:]...); err != nil {
			log.Fatalf("goose run: %v", err)
			flags.Usage()
		}
		break
	}
}

func usage() {
	usagestring := `
to run the app :
	[app_name] run
	example : covid19kalteng run

to update db :
	[app_name] migrate [goose_command]
	example : covid19kalteng migrate up
	goose command lists:
		up                   Migrate the DB to the most recent version available
		up-by-one            Migrate the DB up by 1
		up-to VERSION        Migrate the DB to a specific VERSION
		down                 Roll back the version by 1
		down-to VERSION      Roll back to a specific VERSION
		redo                 Re-run the latest migration
		reset                Roll back all migrations
		status               Dump the migration status for the current DB
		version              Print the current version of the database
		create NAME [sql|go] Creates new migration file with the current timestamp
		fix                  Apply sequential ordering to migrations

database seeding : (development environment only)
	[app_name] seed
	example : covid19kalteng seed

database truncate : (development environment only)
	[app_name] truncate [table(s)]
	example : covid19kalteng truncate borrowers | covid19kalteng truncate borrowers loans | covid19kalteng truncate all
	replace [table] with 'all' to truncate all tables
	`

	log.Print(usagestring)
}
