存放所有语言的从 数据库类型 到 语言类型 的映射关系，这个映射关系应该是做兜底的，
用户可以配置自己的映射关系，如果没有配置，则使用这个映射关系。

```yaml
type-mapping:
  java:
    - type: int
      target: java.lang.Integer
    - type: varchar.*
      target: java.lang.String
    - type: timestamp
      target: java.util.Date
  go:
    - type: int
      target: 'int'
      target-when-nullable: '*int'
      target-when-nonnull: 'int'
    - type: varchar.*
      target: 'string'
    - type: timestamp
      target: time.Time
  rust:
    - type: int
      target: i32
    - type: varchar.*
      target: String
    - type: timestamp
      target: chrono::NaiveDateTime
```