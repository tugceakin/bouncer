package config

import (
	"database/sql"
	"fmt"
	"github.com/stretchr/testify/assert"
	"log"
	"os"
	"testing"
)

func TestConfig(t *testing.T) {
	assert := assert.New(t)
	fmt.Println("----Test for Config Module----")
	os.Remove("./bouncer.db")
	var err error
	db, err = sql.Open("sqlite3", "./bouncer.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	sqlStmt := `
    create table BackendServer (id integer not null primary key autoincrement, host text, config_id integer, foreign key(config_id) references Config);
    create table Config (id integer not null primary key autoincrement, host text);
    delete from Config;
    delete from BackendServer;
	`
	if _, err = db.Exec(sqlStmt); err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}

	_, err = GetConfig("testHost")
	fmt.Println(err)

	backendServers := []BackendServer{
		NewBackendServer("app_server1"),
		NewBackendServer("app_server2"),
		NewBackendServer("app_server3"),
	}

	config, err := NewConfig("testHost", backendServers)
	assert.Equal(nil, err)

	err = config.Save()
	assert.Equal(nil, err)

	retrievedConfig, err := GetConfig("testHost")
	assert.Equal(nil, err)

	assert.Equal("testHost", retrievedConfig.Host)
	assert.Equal("app_server1", retrievedConfig.BackendServers[0].Host)
	assert.Equal("app_server2", retrievedConfig.BackendServers[1].Host)
	assert.Equal("app_server3", retrievedConfig.BackendServers[2].Host)

	config, err = NewConfig("testHost", backendServers)
	assert.Equal("Config already exists", err.Error())

	backendServers[2] = NewBackendServer("app_server4")
	retrievedConfig.BackendServers = backendServers
	err = retrievedConfig.Save()
	assert.Equal(nil, err)

	check, err := GetConfig("testHost")
	assert.Equal(nil, err)
	assert.Equal("testHost", check.Host)
	assert.Equal("app_server1", check.BackendServers[0].Host)
	assert.Equal("app_server2", check.BackendServers[1].Host)
	assert.Equal("app_server4", check.BackendServers[2].Host)
}
