# Gencoder

A code generator for any language or framework, based on [Handlebars](https://handlebarsjs.com/).

## Background

When generating code, you cannot add custom code, because it will be overwritten the next time you generate code. In
real world, code generators usually help to create some boilerplate code (like CRUD operations), and then you may
modify the generated code. However, when you regenerate the code, your changes will be overwritten, and you have to
merge the code manually (which is really annoying).

Gencoder is designed to solve this problem. It allows the code generator to recognize which parts of the file are
automatically generated and which parts are added manually. When regenerating the code, only the generated parts will be
overwritten, and the manual changes will remain untouched. No need for manual merging anymore.

## Install

```bash
go install github.com/DanielLiu1123/gencoder/cmd/gencoder@latest
```

Build from source:

```bash
make && CGO_ENABLED=0 go build -o gencoder cmd/gencoder/main.go
```

## Quick Start

Run a MySQL server:

```bash
docker rm -f test_mysql && docker run --name test_mysql -e MYSQL_ROOT_PASSWORD=root -e MYSQL_DATABASE=testdb -p 3306:3306 -p 33060:33060 -id mysql:latest && sleep 10 && docker exec -i test_mysql mysql -uroot -proot -e "\
    CREATE TABLE testdb.user ( \
        id INT AUTO_INCREMENT PRIMARY KEY, \
        username VARCHAR(64) NOT NULL COMMENT 'Username, required', \
        created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP COMMENT 'Record creation timestamp' \
    ) COMMENT='User account information';"
```

Init basic configuration:

```bash
gencoder init
```

Generate code:

```bash
gencoder generate
```

## Contributing

The [issue tracker](https://github.com/DanielLiu1123/gencoder/issues) is the preferred channel for bug reports,
feature requests and submitting pull requests.

If you would like to contribute to the project, please refer to [Contributing](./CONTRIBUTING.md).

## License

The MIT License.