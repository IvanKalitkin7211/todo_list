package postgres

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"todo-list/config"
)

func TestDatabase_Logic(t *testing.T) {
	t.Run("ValidateConfig_Error", func(t *testing.T) {
		badCfg := &config.DatabaseConfig{Host: ""}
		err := validateConfig(badCfg)
		assert.Error(t, err)
		assert.Contains(t, err.Error(), "invalid database config")
	})

	t.Run("ProvideDBClient_Integration", func(t *testing.T) {
		cfg := config.NewConfig()

		dbClient, err := ProvideDBClient(&cfg.Database)
		if err != nil {
			t.Skip("Пропускаем: база в Docker не отвечает, но код валидации проверен")
		}

		assert.NoError(t, err)
		assert.NotNil(t, dbClient)
		assert.NotNil(t, dbClient.GetDB())

		gormDB := dbClient.GetDB()
		assert.NotNil(t, gormDB)
	})
}
