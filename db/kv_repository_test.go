package db

import (
	"testing"

	"github.com/karintomania/kaigai-go-scraper/common"
	"github.com/stretchr/testify/require"
)

func TestKvRepository(t *testing.T) {
	common.SetLogger()

	db, cleanup := getTestEmptyDbConnection()
	defer cleanup()

	kvr := NewKvRepository(db)
	kvr.CreateTable()

	t.Run("Insert and Find", func(t *testing.T) {
		kvr.Truncate()

		testKey := "test-key"
		testValue := "test value 12345"

		kv := Kv{testKey, testValue}

		kvr.Insert(&kv)

		result := kvr.FindByKey(testKey)

		require.Equal(t, testValue, result)
	})

	t.Run("Update", func(t *testing.T) {
		kvr.Truncate()

		testKey := "test-key"
		testValue := "test value first"
		testValueUpdated := "test value first"

		kv := Kv{testKey, testValue}

		kvr.Insert(&kv)

		kv.Value = testValueUpdated

		kvr.Update(&kv)

		result := kvr.FindByKey(testKey)

		require.Equal(t, testValueUpdated, result)
	})

}
