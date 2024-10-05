#!/usr/bin/env bash

# create MySQL
docker rm -f test_mysql && docker run --name test_mysql -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=testdb -p 3306:3306 -p 33060:33060 -id mysql:latest && sleep 10 && docker exec -i test_mysql mysql -uroot -proot -e "\
    CREATE TABLE testdb.user ( \
        id INT AUTO_INCREMENT PRIMARY KEY, \
        username VARCHAR(64) NOT NULL COMMENT 'Username, required', \
        password VARCHAR(128) NOT NULL, \
        email VARCHAR(128) NOT NULL DEFAULT '' COMMENT 'User email, required', \
        first_name VARCHAR(64) COMMENT 'First name of the user', \
        last_name VARCHAR(64) COMMENT 'Last name of the user', \
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Record creation timestamp', \
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT 'Record update timestamp', \
        status ENUM('active', 'inactive', 'suspended') DEFAULT 'active' COMMENT 'Account status', \
        INDEX idx_name (username), \
        UNIQUE INDEX idx_email (email), \
        INDEX idx_status_created (status, created_at), \
        INDEX idx_full_name (first_name, last_name) \
    ) COMMENT='User account information';"

# create PostgreSQL
docker rm -f test_postgres && docker run --name test_postgres -e POSTGRES_USER=root -e POSTGRES_PASSWORD=root -e POSTGRES_DB=testdb -p 5432:5432 -id postgres:latest && sleep 5 && docker exec -i test_postgres psql -U root -d testdb -c "\
    CREATE TABLE \"user\" ( \
        id SERIAL PRIMARY KEY, \
        username VARCHAR(64) NOT NULL, \
        password VARCHAR(128) NOT NULL, \
        email VARCHAR(128) NOT NULL DEFAULT '', \
        first_name VARCHAR(64), \
        last_name VARCHAR(64), \
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, \
        updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, \
        status VARCHAR(9) DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')), \
        CONSTRAINT unique_email UNIQUE (email) \
    ); \
    CREATE INDEX idx_name ON \"user\" (username); \
    CREATE INDEX idx_status_created ON \"user\" (status, created_at); \
    CREATE INDEX idx_full_name ON \"user\" (first_name, last_name); \
    COMMENT ON COLUMN \"user\".username IS 'Username, required'; \
    COMMENT ON COLUMN \"user\".email IS 'User email, required'; \
    COMMENT ON COLUMN \"user\".first_name IS 'First name of the user'; \
    COMMENT ON COLUMN \"user\".last_name IS 'Last name of the user'; \
    COMMENT ON COLUMN \"user\".created_at IS 'Record creation timestamp'; \
    COMMENT ON COLUMN \"user\".updated_at IS 'Record update timestamp'; \
    COMMENT ON COLUMN \"user\".status IS 'Account status'; \
    COMMENT ON TABLE \"user\" IS 'User account information';"

# create MSSQL
docker rm -f test_sqlserver && \
docker run --name test_sqlserver -e ACCEPT_EULA=Y -e MSSQL_SA_PASSWORD=Sa123456.. -p 1433:1433 -id mcr.microsoft.com/mssql/server:2022-latest && sleep 10 && \
docker exec -i test_sqlserver /opt/mssql-tools18/bin/sqlcmd -S localhost -U sa -P 'Sa123456..' -No -Q "CREATE DATABASE testdb;" && \
docker exec -i test_sqlserver /opt/mssql-tools18/bin/sqlcmd -S localhost -U sa -P 'Sa123456..' -No -d 'testdb' -Q "\
    CREATE TABLE [user] (
		id         INT IDENTITY (1,1) PRIMARY KEY,
		username   NVARCHAR(64)  NOT NULL,
		password   NVARCHAR(128) NOT NULL,
		email      NVARCHAR(128) NOT NULL DEFAULT '',
		first_name NVARCHAR(64),
		last_name  NVARCHAR(64),
		created_at DATETIME               DEFAULT CURRENT_TIMESTAMP,
		updated_at DATETIME               DEFAULT CURRENT_TIMESTAMP,
		status     NVARCHAR(9)            DEFAULT 'active' CHECK (status IN ('active', 'inactive', 'suspended')),
		deleted_at DATETIME      NULL,
		CONSTRAINT unique_email UNIQUE (email)
	); \
    CREATE INDEX idx_name ON [user] (username); \
    CREATE INDEX idx_status_created ON [user] (status, created_at); \
    CREATE INDEX idx_full_name ON [user] (first_name, last_name); \
    EXEC sp_addextendedproperty @name = N'MS_Description', @value = N'User account information', @level0type = N'SCHEMA', @level0name = 'dbo', @level1type = N'TABLE', @level1name = 'user'; \
    EXEC sp_addextendedproperty @name = N'MS_Description', @value = N'Username, required', @level0type = N'SCHEMA', @level0name = 'dbo', @level1type = N'TABLE', @level1name = 'user', @level2type = N'COLUMN', @level2name = 'username'; \
    EXEC sp_addextendedproperty @name = N'MS_Description', @value = N'User email, required', @level0type = N'SCHEMA', @level0name = 'dbo', @level1type = N'TABLE', @level1name = 'user', @level2type = N'COLUMN', @level2name = 'email'; \
    EXEC sp_addextendedproperty @name = N'MS_Description', @value = N'First name of the user', @level0type = N'SCHEMA', @level0name = 'dbo', @level1type = N'TABLE', @level1name = 'user', @level2type = N'COLUMN', @level2name = 'first_name'; \
    EXEC sp_addextendedproperty @name = N'MS_Description', @value = N'Last name of the user', @level0type = N'SCHEMA', @level0name = 'dbo', @level1type = N'TABLE', @level1name = 'user', @level2type = N'COLUMN', @level2name = 'last_name'; \
    EXEC sp_addextendedproperty @name = N'MS_Description', @value = N'Record creation timestamp', @level0type = N'SCHEMA', @level0name = 'dbo', @level1type = N'TABLE', @level1name = 'user', @level2type = N'COLUMN', @level2name = 'created_at'; \
    EXEC sp_addextendedproperty @name = N'MS_Description', @value = N'Record update timestamp', @level0type = N'SCHEMA', @level0name = 'dbo', @level1type = N'TABLE', @level1name = 'user', @level2type = N'COLUMN', @level2name = 'updated_at'; \
    EXEC sp_addextendedproperty @name = N'MS_Description', @value = N'Account status', @level0type = N'SCHEMA', @level0name = 'dbo', @level1type = N'TABLE', @level1name = 'user', @level2type = N'COLUMN', @level2name = 'status';"
