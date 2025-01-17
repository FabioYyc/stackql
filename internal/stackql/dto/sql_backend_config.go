package dto

import (
	"fmt"
	"os"
	"strings"

	"github.com/stackql/stackql/internal/stackql/constants"
	"github.com/xo/dburl"
	"gopkg.in/yaml.v2"
)

type SQLBackendSchemata struct {
	TableSchema     string `json:"tableSchema" yaml:"tableSchema"`
	IntelViewSchema string `json:"intelViewSchema" yaml:"intelViewSchema"`
	OpsViewSchema   string `json:"opsViewSchema" yaml:"opsViewSchema"`
}

type SQLBackendCfg struct {
	DbEngine              string             `json:"dbEngine" yaml:"dbEngine"`
	DSN                   string             `json:"dsn" yaml:"dsn"`
	DSNEnvVar             string             `json:"dsnEnvVar" yaml:"dsnEnvVar"`
	Schemata              SQLBackendSchemata `json:"schemata" yaml:"schemata"`
	DbInitFilePath        string             `json:"dbInitFilepath" yaml:"dbInitFilepath"`
	SQLSystem             string             `json:"sqlDialect" yaml:"sqlDialect"`
	SchemaType            string             `json:"schemaType" yaml:"schemaType"`
	InitMaxRetries        int                `json:"initMaxRetries" yaml:"initMaxRetries"`
	InitRetryInitialDelay int                `json:"initRetryInitialDelay" yaml:"initRetryInitialDelay"`
}

func (sqlCfg SQLBackendCfg) GetDSN() string {
	dsn := sqlCfg.DSN
	if sqlCfg.DSNEnvVar != "" {
		dsn = os.Getenv(sqlCfg.DSNEnvVar)
	}
	return dsn
}

func (sqlCfg SQLBackendCfg) GetSchemaType() string {
	return sqlCfg.SchemaType
}

func (sqlCfg SQLBackendCfg) GetDatabaseName() (string, error) {
	dsn := sqlCfg.GetDSN()
	if dsn == "" {
		return "", fmt.Errorf("cannot ionfer db name from empty dsn")
	}
	dbUrl, err := dburl.Parse(dsn)
	if err != nil {
		return "", fmt.Errorf("error parsing postgres dsn: %s", err.Error())
	}
	if dbUrl == nil {
		return "", fmt.Errorf("error parsing postgres dsn, nil url generated")
	}
	return strings.TrimLeft(dbUrl.Path, "/"), nil
}

func (sqlCfg SQLBackendCfg) GetTableSchemaName() string {
	return sqlCfg.Schemata.TableSchema
}

func (sqlCfg SQLBackendCfg) GetOpsViewSchemaName() string {
	return sqlCfg.Schemata.OpsViewSchema
}

func (sqlCfg SQLBackendCfg) GetIntelViewSchemaName() string {
	return sqlCfg.Schemata.IntelViewSchema
}

func GetSQLBackendCfg(s string) (SQLBackendCfg, error) {
	rv := SQLBackendCfg{}
	err := yaml.Unmarshal([]byte(s), &rv)
	if rv.DbEngine == "" {
		rv.DbEngine = constants.DbEngineDefault
	}
	if rv.SQLSystem == "" {
		rv.SQLSystem = constants.SQLDialectDefault
	}
	if rv.InitMaxRetries < 1 {
		rv.InitMaxRetries = 10
	}
	if rv.InitRetryInitialDelay < 1 {
		rv.InitRetryInitialDelay = 10
	}
	return rv, err
}
