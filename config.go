package config

import (
	"database/sql"
	"errors"
	_ "github.com/mattn/go-sqlite3"
	"log"
)

var db *sql.DB

type Config struct {
	Id                int
	Host              string
	BackendServers    []BackendServer
	NextBackendServer chan BackendServer
	done              chan struct{}
}

type BackendServer struct {
	Id       int
	Host     string
	ConfigId int
}

func (config *Config) NextBackendServerRoutine() {
	for {
		for _, next := range config.BackendServers {
			select {
			case config.NextBackendServer <- next:

			case <-config.done:
				return
			}
		}

	}
}

func (config *Config) Destroy() {
	config.done <- struct{}{}
}

func NewConfig(host string, backendServers []BackendServer) (Config, error) {
	var config Config
	if _, err := GetConfigId(host); err == nil {
		log.Println(err)
		return config, errors.New("Config already exists")
	}
	config.Host = host
	config.BackendServers = backendServers
	config.NextBackendServer = make(chan BackendServer)
	config.done = make(chan struct{})
	go config.NextBackendServerRoutine()
	return config, nil
}

func NewBackendServer(host string) BackendServer {
	var backendServer BackendServer
	backendServer.Host = host
	return backendServer
}

func (config *Config) Save() error {
	if config.Id == 0 {

		tx, err := db.Begin()
		if err != nil {
			log.Println(err)
			return err
		}
		stmt, err := tx.Prepare("insert into Config(host) values(?)")
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return err
		}
		defer stmt.Close()
		if _, err = stmt.Exec(config.Host); err != nil {
			log.Println(err)
			tx.Rollback()
			return err
		}
		tx.Commit()

		id, err := GetConfigId(config.Host)
		if err != nil {
			log.Println(err)
			return err
		}

		tx, err = db.Begin()
		if err != nil {
			log.Println(err)
			return err
		}
		stmt, err = tx.Prepare("delete from BackendServer where config_id = ?")
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return err
		}
		defer stmt.Close()
		if _, err = stmt.Exec(id); err != nil {
			log.Println(err)
			tx.Rollback()
			return err
		}
		tx.Commit()

		for _, backendServer := range config.BackendServers {
			backendServer.ConfigId = id
			if err = backendServer.Save(); err != nil {
				log.Println(err)
				return err
			}
		}
		return nil
	} else {
		tx, err := db.Begin()
		if err != nil {
			log.Println(err)
			return err
		}
		stmt, err := tx.Prepare("delete from BackendServer where config_id = ?")
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return err
		}
		defer stmt.Close()
		if _, err = stmt.Exec(config.Id); err != nil {
			log.Println(err)
			tx.Rollback()
			return err
		}
		tx.Commit()

		for _, backendServer := range config.BackendServers {
			backendServer.ConfigId = config.Id
			if err = backendServer.Save(); err != nil {
				log.Println(err)
				return err
			}
		}
		return nil
	}
}

func (backendServer *BackendServer) Save() error {
	if backendServer.Id == 0 {

		tx, err := db.Begin()
		if err != nil {
			log.Println(err)
			return err
		}
		stmt, err := tx.Prepare("insert into BackendServer(host, config_id) values(?,?)")
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return err
		}
		defer stmt.Close()
		_, err = stmt.Exec(backendServer.Host, backendServer.ConfigId)
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return err
		}
		tx.Commit()
		return nil

	} else {

		tx, err := db.Begin()
		if err != nil {
			log.Println(err)
			return err
		}
		stmt, err := tx.Prepare("insert into BackendServer(id, host, config_id) values(?,?,?)")
		if err != nil {
			log.Println(err)
			tx.Rollback()
			return err
		}
		defer stmt.Close()
		if _, err = stmt.Exec(backendServer.Id, backendServer.Host, backendServer.ConfigId); err != nil {
			log.Println(err)
			tx.Rollback()
			return err
		}
		tx.Commit()
		return nil
	}
}

func GetConfigId(host string) (int, error) {
	stmt, err := db.Prepare("select id from Config where host = ?")
	if err != nil {
		log.Println(err)
		return 0, err
	}
	defer stmt.Close()
	row := stmt.QueryRow(host)

	var id int
	err = row.Scan(&id)
	return id, err
}

func GetConfig(host string) (Config, error) {
	var config Config
	var err error
	config.Id, err = GetConfigId(host)
	if err != nil {
		log.Println(err)
		return config, err
	}
	config.Host = host

	stmt, err := db.Prepare("select * from BackendServer where config_id = ?")
	if err != nil {
		log.Println(err)
		return config, err
	}
	defer stmt.Close()

	rows, err := stmt.Query(config.Id)
	if err != nil {
		log.Println(err)
		return config, err
	}
	defer rows.Close()
	for rows.Next() {
		var backendServer BackendServer
		if err = rows.Scan(&backendServer.Id, &backendServer.Host, &backendServer.ConfigId); err != nil {
			return config, err
		}
		config.BackendServers = append(config.BackendServers, backendServer)
	}
	return config, err
}
