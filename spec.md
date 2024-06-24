
## Variables

- `{{table}}`: The name of the table.
  - `{{columns}}`: The columns of the table.
    - `{{column}}`: The name of the column.
  - `{{description}}`: The description of the table.


```yaml
tables:
  - name: user
    columns:
      - name: id
        type: int
        description: The id of the user.
      - name: name
        type: string
        description: The name of the user.
    index:
      - name: id
        type: primary key
    description: The users table.
```

## Usage

```shell
gencoder -f gencoder.yml
```