#!/usr/bin/env bash

set -eux

go install github.com/cratonica/2goarray@latest

2goarray mockSQL1 db_importer < ./.data/mock/2022-08-01-001.sql > mock_sql1_test.go
2goarray mockSQL2 db_importer < ./.data/mock/2022-08-03-001.sql > mock_sql2_test.go
2goarray mockSQL3 db_importer < ./.data/mock/2022-08-03-002.sql > mock_sql3_test.go
2goarray mockDBImporterYAML db_importer < ./.data/mock/db-importer.yaml > mock_dbimporter_yaml_test.go
