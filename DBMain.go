package main

import (
	"ProjetoUnivesp2020/managers/database"
	"encoding/json"
	"errors"
	"fmt"
	jsoniter "github.com/json-iterator/go"
	"github.com/urfave/cli/v2"
	"os"
)

var (
	Categorys = []string{"Manage", "Backup"}
)

func main() {
	app := cli.NewApp()

	app.Name = "Gerenciador do banco de dados do Projeto integrador"
	app.EnableBashCompletion = true
	app.Usage = "Realize consultas, alterações e backup"

	app.Flags = []cli.Flag{
		&cli.PathFlag{
			Name:     "database-path",
			Aliases:  []string{"d"},
			Usage:    "Caminho da pasta do banco de dados",
			EnvVars:  []string{"DATABASE_PATH"},
			Required: true,
		},
	}

	app.Commands = []*cli.Command{
		&cli.Command{
			Name:        "select",
			Description: "make a select in database",
			Category:    Categorys[0],
			Action: func(c *cli.Context) error {
				Conn := database.InitDataBase(c.String("database-path"), false)

				if Conn == nil {
					return errors.New("Database not found")
				}

				selectedCollection := c.String("collection")

				if selectedCollection == "" {
					l, _ := jsoniter.Marshal(Conn.AllCols())
					fmt.Println(string(l))
					return nil
				}

				fields := c.StringSlice("fields")
				collection := Conn.Use(selectedCollection)

				if collection == nil {
					return errors.New("collection not found")
				}

				first := true
				fmt.Print("[")
				collection.ForEachDoc(func(id int, doc []byte) (moveOn bool) {
					if first {
						first = false
					} else {
						fmt.Print(",")
					}

					if fields == nil {
						fmt.Print(string(doc))
					} else {
						var m map[string]interface{}
						pm := make(map[string]interface{})
						_ = json.Unmarshal(doc, &m)

						for _, f := range fields {
							pm[f] = m[f]
						}
						l, _ := json.Marshal(pm)
						fmt.Print(string(l))
					}

					return true
				})
				fmt.Print("]")

				return nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:    "collection",
					Aliases: []string{"c"},
					Usage:   "Choose the colletion to do the select, if not set will show all colletions",
					Value:   "",
				},
				&cli.StringSliceFlag{
					Name:    "fields",
					Aliases: []string{"f"},
					Usage:   "Filter the showing fields",
					Value:   nil,
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		panic(err)
	}
}
