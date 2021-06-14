package db

import (
	"Panda/conf"
	"fmt"
	"github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/agent"
	"net"
	"os"
	"sync"
	"time"
)

// NewMySQLByDSN connects to MySQL by DSN
func NewMySQLByDSN(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open("mysql", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mysql: %s", err)
	}

	db.DB().SetConnMaxLifetime(60 * time.Second)
	db.DB().SetMaxOpenConns(2000)
	db.DB().SetMaxIdleConns(100)

	return db, nil
}

// GenMySQLDSN generates DSN for MySQL
func GenMySQLDSN(cfg conf.MySQL) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&timeout=30s&parseTime=true",
		cfg.Username, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
}

// NewMySQL connects to MySQL with config struct
func NewMySQL(cfg conf.MySQL) (*gorm.DB, error) {
	return NewMySQLByDSN(GenMySQLDSN(cfg))
}

var mySQL = make(map[string]*gorm.DB)
var mySQLLock = sync.Mutex{}

// GetMySQL create MySQL connection instance
func GetMySQL(cfg conf.MySQL) (*gorm.DB, error) {
	dsn := GenMySQLDSN(cfg)
	if db, ok := mySQL[dsn]; ok {
		return db, nil
	}

	mySQLLock.Lock()
	defer mySQLLock.Unlock()
	if db, ok := mySQL[dsn]; ok {
		return db, nil
	}

	if db, err := NewMySQLByDSN(dsn); err == nil {
		mySQL[dsn] = db
		return mySQL[dsn], nil
	} else {
		return nil, err
	}
}

type SSHDialer struct {
	*ssh.Client
}

func (dialer *SSHDialer) Dial(addr string) (net.Conn, error) {
	return dialer.Client.Dial("tcp", addr)
}

func RegisterMySQLDail(cfg conf.SSH) (*ssh.Client, error) {
	var agentClient agent.Agent

	// Establish a connection to the local ssh-agent
	if conn, err := net.Dial("unix", os.Getenv("SSH_AUTH_SOCK")); err == nil {
		// Create a new instance of the ssh agent
		agentClient = agent.NewClient(conn)
	}

	// The client configuration with configuration option to use the ssh-agent
	sshConfig := &ssh.ClientConfig{
		User: cfg.User,
		Auth: []ssh.AuthMethod{},
	}

	// When the agentClient connection succeeded, add them as AuthMethod
	if agentClient != nil {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PublicKeysCallback(agentClient.Signers))
	}

	// When there's a non empty password add the password AuthMethod
	if cfg.Password != "" {
		sshConfig.Auth = append(sshConfig.Auth, ssh.PasswordCallback(func() (string, error) {
			return cfg.Password, nil
		}))
	}

	sshConfig.HostKeyCallback = ssh.InsecureIgnoreHostKey()

	// Connect to the SSH Server
	sshConn, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", cfg.Host, cfg.Port), sshConfig)
	if err != nil {
		return nil, err
	}

	// Now we register the SSHDialer with the ssh connection as a parameter
	mysql.RegisterDial("tcp", (&SSHDialer{sshConn}).Dial)
	return sshConn, nil
}
