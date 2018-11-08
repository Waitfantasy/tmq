package persistent

import (
	"database/sql"
	"github.com/Waitfantasy/tmq/message"
	_ "github.com/go-sql-driver/mysql"
)

type MysqlPersistent struct {
	conn *sql.DB
}

// dsn: [username[:password]@][protocol[(address)]]/dbname[?param1=value1&...&paramN=valueN]
func NewMysqlPersistent(dsn string) (*MysqlPersistent, error) {
	if db, err := sql.Open("mysql", dsn); err != nil {
		return nil, err
	} else {
		return &MysqlPersistent{
			conn: db,
		}, nil
	}
}

func (p *MysqlPersistent) Store(msg *message.Message) error {
	panic("implement me")
}

func (p *MysqlPersistent) Update(msg *message.Message) error {
	panic("implement me")
}

func (p *MysqlPersistent) Find(id uint64) (*message.Message, error) {
	panic("implement me")
}
