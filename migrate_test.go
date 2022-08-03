package db_importer

import (
	"fmt"
	"github.com/ismdeep/rand"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"os"
	"testing"
)

func TestImport(t *testing.T) {
	configRoot := fmt.Sprintf("/tmp/%v", rand.TimeBased())
	assert.NoError(t, os.MkdirAll(fmt.Sprintf("%v/sql", configRoot), 0777))
	assert.NoError(t, ioutil.WriteFile(fmt.Sprintf("%v/db-importer.yaml", configRoot), mockDBImporterYAML, 0777))
	assert.NoError(t, ioutil.WriteFile(fmt.Sprintf("%v/sql/2022-08-01-001.sql", configRoot), mockSQL1, 0777))
	assert.NoError(t, ioutil.WriteFile(fmt.Sprintf("%v/sql/2022-08-03-001.sql", configRoot), mockSQL2, 0777))
	assert.NoError(t, ioutil.WriteFile(fmt.Sprintf("%v/sql/2022-08-03-002.sql", configRoot), mockSQL3, 0777))
	assert.NoError(t,
		Migrate(configRoot))
}
