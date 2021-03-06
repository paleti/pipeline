package oracle

import (
	"fmt"
	"strings"

	"github.com/banzaicloud/pipeline/pkg/providers/oracle"
	"github.com/banzaicloud/pipeline/pkg/providers/oracle/model"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

// Migrate executes the table migrations for the provider.
func Migrate(db *gorm.DB, logger logrus.FieldLogger) error {
	tables := []interface{}{
		&ObjectStoreBucketModel{},
		&model.Cluster{},
		&model.NodePool{},
		&model.NodePoolSubnet{},
		&model.NodePoolLabel{},
		&model.Profile{},
		&model.ProfileNodePool{},
		&model.ProfileNodePoolLabel{},
	}

	var tableNames string
	for _, table := range tables {
		tableNames += fmt.Sprintf(" %s", db.NewScope(table).TableName())
	}

	logger.WithFields(logrus.Fields{
		"provider":    oracle.Provider,
		"table_names": strings.TrimLeft(tableNames, " "),
	}).Info("migrating provider tables")

	return db.AutoMigrate(tables...).Error
}
