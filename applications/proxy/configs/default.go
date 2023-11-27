package configs

var (
	DefaultConfig = []byte(`
version: "v1.0.0"
botToken: "this_is_bot_token"
admins: "3333333,2222222,1111111"
db:
  type: "postgreSQL"
  username: "postgres"
  password: "postgres"
  host: "localhost"
  port: 5432
  name: "proxy"
  is_ssl_enable: false
handlers:
  publishChannelId: -1111111111111111
`)
)
