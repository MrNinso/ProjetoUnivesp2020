package main

import (
	"ProjetoUnivesp2020/managers/database"
	"bufio"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"

	jsoniter "github.com/json-iterator/go"
	"github.com/urfave/cli/v2"
)

type handleType func(Conn *database.DB, args []string) error

type BackupInstruction struct {
	name   string
	handle handleType
}

type BackupInstructionList []*BackupInstruction

func (b BackupInstructionList) exec(Conn *database.DB, line string) error {
	for _, instruction := range b {
		if strings.HasPrefix(line, instruction.name) {
			args := strings.Split(line, "|")
			return instruction.handle(Conn, args[1:])
		}
	}
	if strings.HasPrefix(line, "-") {
		return nil
	}
	return errors.New("Invalid line:" + line)
}

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

					var m map[string]interface{}
					pm := make(map[string]interface{})
					_ = jsoniter.Unmarshal(doc, &m)

					if fields != nil {
						for _, f := range fields {
							pm[f] = m[f]
						}
					} else {
						pm = m
					}

					pm["id"] = strconv.Itoa(id)

					l, _ := jsoniter.Marshal(pm)
					fmt.Print(string(l))

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
		&cli.Command{
			Name:        "update",
			Description: "update a entry in database",
			Category:    Categorys[0],
			Action: func(c *cli.Context) error {
				Conn := database.InitDataBase(c.String("database-path"), false)

				col := Conn.Use(c.String("collection"))

				if col == nil {
					return errors.New("collection not found")
				}

				var m map[string]interface{}
				err := jsoniter.Unmarshal([]byte(c.String("json")), &m)

				if err != nil {
					return err
				}

				id := c.Int("id")

				oldM, err := col.Read(id)

				if err != nil {
					return err
				}

				for key, value := range m {
					oldM[key] = value
				}

				return col.Update(id, oldM)
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "collection",
					Aliases:  []string{"c"},
					Usage:    "Choose the colletion to update",
					Required: true,
					Value:    "",
				},
				&cli.IntFlag{
					Name:     "id",
					Usage:    "id of target entry",
					Required: true,
				},
				&cli.StringFlag{
					Name:     "json",
					Aliases:  []string{"j"},
					Usage:    "new values in json format",
					Required: true,
					Value:    "",
					EnvVars:  []string{"UPDATE_JSON"},
				},
			},
		},
		&cli.Command{
			Name:        "insert",
			Description: "insert a entry in database",
			Category:    Categorys[0],
			Action: func(c *cli.Context) error {
				Conn := database.InitDataBase(c.String("database-path"), false)

				col := Conn.Use(c.String("collection"))

				if col == nil {
					return errors.New("collection not found")
				}

				var m map[string]interface{}
				err := jsoniter.Unmarshal([]byte(c.String("json")), &m)

				if err != nil {
					return err
				}

				if err != nil {
					return err
				}

				id, err := col.Insert(m)

				m["id"] = strconv.Itoa(id)

				if err != nil {
					return err
				}

				b, _ := jsoniter.Marshal(m)
				log.Println(string(b))

				return nil
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "collection",
					Aliases:  []string{"c"},
					Usage:    "Choose the colletion to insert",
					Required: true,
					Value:    "",
				},
				&cli.StringFlag{
					Name:     "json",
					Aliases:  []string{"j"},
					Usage:    "new values in json format",
					Required: true,
					Value:    "",
					EnvVars:  []string{"INSERT_JSON"},
				},
			},
		},
		&cli.Command{
			Name:        "delete",
			Description: "delete a entry in database",
			Category:    Categorys[0],
			Action: func(c *cli.Context) error {
				Conn := database.InitDataBase(c.String("database-path"), false)

				col := Conn.Use(c.String("collection"))

				if col == nil {
					return errors.New("collection not found")
				}

				return col.Delete(c.Int("id"))
			},
			Flags: []cli.Flag{
				&cli.StringFlag{
					Name:     "collection",
					Aliases:  []string{"c"},
					Usage:    "Choose the colletion to update",
					Required: true,
					Value:    "",
				},
				&cli.IntFlag{
					Name:     "id",
					Usage:    "id of target entry",
					Required: true,
				},
			},
		},
		&cli.Command{
			Name:        "backup",
			Description: "backup database",
			Category:    Categorys[1],
			Action: func(c *cli.Context) error {
				Conn := database.InitDataBase(c.String("database-path"), false)
				colls := Conn.AllCols()

				for _, colName := range colls {
					colNameB64 := b64.StdEncoding.EncodeToString([]byte(colName))
					fmt.Printf("CREATE|%s\n", colNameB64)
				}

				for _, colName := range colls {
					col := Conn.Use(colName)

					col.ForEachDoc(func(id int, doc []byte) bool {
						fmt.Printf("INSERT|%s|%s\n",
							b64.StdEncoding.EncodeToString([]byte(colName)),
							b64.StdEncoding.EncodeToString(doc),
						)
						return true
					})
				}
				return nil
			},
		},
		&cli.Command{
			Name:        "restore",
			Description: "restore database",
			Category:    Categorys[1],
			Action: func(c *cli.Context) error {
				Conn := database.InitDataBase(c.String("database-path"), false)
				file, err := os.Open(c.Path("backup-file"))
				handles := loadBackupInstructions()

				if err != nil {
					log.Fatal(err)
				}
				defer file.Close()

				scanner := bufio.NewScanner(file)
				for scanner.Scan() {
					err := handles.exec(Conn, scanner.Text())
					if err != nil {
						return err
					}
				}
				Conn.Dump()
				return scanner.Err()
			},
			Flags: []cli.Flag{
				&cli.PathFlag{
					Name:     "backup-file",
					Aliases:  []string{"b"},
					Required: true,
				},
			},
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func loadBackupInstructions() *BackupInstructionList {
	return &BackupInstructionList{
		{"CREATE", func(Conn *database.DB, args []string) error {
			fmt.Println(args)
			col, err := b64.StdEncoding.DecodeString(args[0])

			if err != nil {
				return err
			}

			return Conn.Create(string(col))
		}},
		{"INSERT", func(Conn *database.DB, args []string) error {
			col, err := b64.StdEncoding.DecodeString(args[0])

			if err != nil {
				return err
			}

			doc, err := b64.StdEncoding.DecodeString(args[0])

			if err != nil {
				return err
			}

			c := Conn.Use(string(col))

			if c == nil {
				return errors.New("Collection " + string(col) + " not found")
			}

			var d map[string]interface{}
			err = jsoniter.Unmarshal(doc, &d)

			if err != nil {
				return err
			}

			_, err = c.Insert(d)

			return err
		}},
	}
}
