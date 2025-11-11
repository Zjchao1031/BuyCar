package config

type mySQL struct {
	Addr     string
	Database string
	Username string
	Password string
	Charset  string
}

type redis struct {
	Addr     string
	Password string
	DB       int
}

type oss struct {
	Endpoint        string
	AccessKeyID     string `mapstructure:"accessKey-id"`
	AccessKeySecret string `mapstructure:"accessKey-secret"`
	BucketName      string
	Zone            string `mapstructure:"zone"`
	UseCdnDomains   bool   `mapstructure:"use_cdn_domains"`
}

// SMTP 配置
type smtp struct {
	Host     string
	Port     int
	User     string
	Password string
	From     string
	FromName string `mapstructure:"from_name"`
}

// Verify（验证码）配置
type verify struct {
	CodeLength int `mapstructure:"code_length"`
	TTLSeconds int `mapstructure:"ttl_seconds"`
}

type server struct {
	Addr string
	Port int
}

type turnstile struct {
	SecretKey string `mapstructure:"secret_key"`
	Enabled   bool
}

// LLM 配置
type llm struct {
    // 当前使用的提供商，例如: "tongyi"
    Provider string
    // 是否启用 LLM 生成功能
    Enabled  bool
    // 通义千问配置
    Tongyi   tongyi `mapstructure:"tongyi"`
}

type tongyi struct {
    // 通义千问 API Key（DashScope 密钥）
    ApiKey         string `mapstructure:"api_key"`
    // 使用的模型名称，示例: qwen-plus
    Model          string `mapstructure:"model"`
    // Chat Completions 端点（可选，默认官方地址）
    Endpoint       string `mapstructure:"endpoint"`
    // 请求超时时间（秒）
    TimeoutSeconds int    `mapstructure:"timeout_seconds"`
}

type config struct {
    MySQL     mySQL
    Redis     redis
    OSS       oss
    Smtp      smtp   `mapstructure:"smtp"`
    Verify    verify `mapstructure:"verify"`
    Server    server
    Turnstile turnstile `mapstructure:"turnstile"`
    LLM       llm      `mapstructure:"llm"`
}
