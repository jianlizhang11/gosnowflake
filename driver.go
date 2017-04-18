// Go Snowflake Driver - Snowflake driver for Go's database/sql package
//
// Copyright (c) 2017 Snowflake Computing Inc. All right reserved.
//

package gosnowflake

import (
	"database/sql"
	"database/sql/driver"
	"log"
	"net"
	"net/http"
)

type SnowflakeDriver struct{}

type DialFunc func(addr string) (net.Conn, error)

func (d SnowflakeDriver) Open(dsn string) (driver.Conn, error) {
	log.Println("Open")
	var err error
	sc := &snowflakeConn{
		SequeceCounter: 0,
	}
	sc.cfg, err = ParseDSN(dsn)
	if err != nil {
		return nil, err
	}

	// Authenticate
	sc.Rest = &snowflakeRestful{
		Host:     sc.cfg.Host,
		Port:     sc.cfg.Port,
		Protocol: sc.cfg.Protocol,
		Client:   &http.Client{}, // create a new client
	}
	sessionParameters := make(map[string]string)
	sessionInfo, err := Authenticate(
		sc.Rest,
		sc.cfg.User,
		sc.cfg.Password,
		sc.cfg.Account,
		sc.cfg.Database,
		sc.cfg.Schema,
		sc.cfg.Warehouse,
		sc.cfg.Role,
		sc.cfg.Passcode,
		sc.cfg.PasscodeInPassword,
		"", // TODO: OKTA support
		"",
		"",
		sessionParameters)
	if err != nil {
		// TODO: error handling
		return nil, nil
	}

	sc.cfg.Database = sessionInfo.DatabaseName
	sc.cfg.Schema = sessionInfo.SchemaName
	sc.cfg.Role = sessionInfo.RoleName
	sc.cfg.Warehouse = sessionInfo.WarehouseName
	return sc, nil
}

func init() {
	sql.Register("snowflake", &SnowflakeDriver{})
}