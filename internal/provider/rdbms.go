package provider

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"stoo-kv/config"
	"strings"
)

type Rdbms struct {
	db  *gorm.DB
	cfg *config.Config
}
type kv struct {
	Namespace string `gorm:"column:namespace"`
	Profile   string `gorm:"column:profile"`
	Key       string `gorm:"column:key"`
	Value     string `gorm:"column:value"`
}

func NewMySql(config *config.Config) (*Rdbms, error) {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		config.Providers.Mysql.Username,
		config.Providers.Mysql.Password,
		config.Providers.Mysql.Host,
		config.Providers.Mysql.Port,
		config.Providers.Mysql.DatabaseName,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Rdbms{db: db, cfg: config}, nil
}

func NewPostgres(config *config.Config) (*Rdbms, error) {
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		config.Providers.Postgres.Host,
		config.Providers.Postgres.Username,
		config.Providers.Postgres.Password,
		config.Providers.Postgres.DatabaseName,
		config.Providers.Postgres.Port,
		config.Providers.Postgres.SslMode,
		config.Providers.Postgres.TimeZone,
	)
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &Rdbms{db: db, cfg: config}, nil
}

func (r *Rdbms) Set(key string, value any) error {
	namespace, profile, keyName := splitKey(key)
	return r.db.Table(r.cfg.RdbmsDefaultTable).Create(map[string]any{
		"namespace": namespace,
		"profile":   profile,
		"key":       keyName,
		"value":     value,
	}).Error
}

func (r *Rdbms) Get(key string) (string, error) {
	namespace, profile, keyName := splitKey(key)
	keyValue := &kv{}
	err := r.db.Limit(1).Debug().Table(r.cfg.RdbmsDefaultTable).Where("`namespace` = ? AND `profile` = ? AND `key` = ?", namespace, profile, keyName).Find(keyValue).Error
	return keyValue.Value, err
}

func (r *Rdbms) Delete(key string) error {
	namespace, profile, keyName := splitKey(key)
	return r.db.Table(r.cfg.RdbmsDefaultTable).Where("`namespace` = ? AND `profile` = ? AND `key` = ?", namespace, profile, keyName).Delete(&kv{}).Error
}

func (r *Rdbms) GetAll() (map[string]string, error) {
	kvMap := make(map[string]string)
	var keyValues []kv
	if err := r.db.Table(r.cfg.RdbmsDefaultTable).Select("`key`", "value").Find(&keyValues).Error; err != nil {
		return nil, err
	}
	for _, entry := range keyValues {
		kvMap[entry.Key] = entry.Value
	}
	return kvMap, nil
}

func (r *Rdbms) GetByNameSpaceAndProfile(namespace, profile string) (map[string]string, error) {
	kvMap := make(map[string]string)
	var keyValues []kv
	if err := r.db.Table(r.cfg.RdbmsDefaultTable).Where("namespace = ? AND profile = ?", namespace, profile).Select("`key`", "value").Find(&keyValues).Error; err != nil {
		return nil, err
	}
	for _, entry := range keyValues {
		kvMap[entry.Key] = entry.Value
	}
	return kvMap, nil
}

func splitKey(key string) (string, string, string) {
	keys := strings.Split(key, "::")
	return keys[0], keys[1], keys[2]
}
