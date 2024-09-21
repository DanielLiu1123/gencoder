# Gencoder

A code generator for any language or framework that preserves your custom changes during regeneration, powered by [Handlebars](https://handlebarsjs.com/).

## Background

You cannot add code to the generated code because it will be overwritten the next time you generate it.
In real world, code generators are often used to create some boilerplate code (like CRUD operations),
then you develop your own code based on it. However, when you regenerate the code, your changes will be overwritten,
and you have to merge them manually (which is very annoying).

Gencoder is designed to solve this problem. It can recognize which parts of the file are generated
and which parts are manually added. When regenerating, it only overwrites the automatically generated parts,
and the manual parts remain unchanged. No more manual merging.

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