# Gencoder

A code generator for any languages/frameworks.

Applicable scenarios for gencoder:

- You need to modify the generated code and the modified code will not be overwritten
- High customization requirements for generated code

## Install

```bash
go install github.com/DanielLiu1123/gencoder/cmd/gencoder
```

## Quick Start

```shell
docker rm -f test_mysql && docker run --name test_mysql -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=testdb -p 3306:3306 -p 33060:33060 -id mysql:latest && sleep 10 && docker exec -i test_mysql mysql -uroot -proot -e "\
    CREATE TABLE testdb.user ( \
        id INT AUTO_INCREMENT PRIMARY KEY, \
        username VARCHAR(64) NOT NULL COMMENT 'Username, required', \
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Record creation timestamp' \
    ) COMMENT='User account information';"
```

```bash
gencoder init
```

```bash
gencoder generate
```
