package configs

var DefaultConfig = []byte(`
version: "v1.0.0"
redis:
  host: "localhost"
  port: 6379
  password: ""
db:
  type: "postgreSQL"
  port: 5432
  host: "localhost"
  username: "username"
  password: "password"
  name: "film"
  is_ssl_enable: false
handlers:
  admins: "3333333,2222222,1111111"
  redisDB: 1
  publicChannelId:
  promote_channels: "-1000,-2000"
`)
