package llm

import (
    "buycar/config"
    "context"
    "errors"
)

// Client 定义统一的 LLM 客户端接口
type Client interface {
    // Generate 传入完整 prompt，返回模型输出的主体文本、原始响应（JSON 字符串）以及实际使用的模型名
    Generate(ctx context.Context, prompt string) (content string, raw string, modelUsed string, err error)
}

var (
    ErrLLMNotEnabled = errors.New("LLM 未启用，请检查配置 llm.enabled")
)

// NewClientFromConfig 根据配置创建对应的 LLM 客户端
func NewClientFromConfig() (Client, error) {
    if config.LLM == nil || !config.LLM.Enabled {
        return nil, ErrLLMNotEnabled
    }

    switch config.LLM.Provider {
    case "tongyi":
        return NewTongyiClient(TongyiOptions{
            ApiKey:         config.LLM.Tongyi.ApiKey,
            Model:          config.LLM.Tongyi.Model,
            Endpoint:       config.LLM.Tongyi.Endpoint,
            TimeoutSeconds: config.LLM.Tongyi.TimeoutSeconds,
        })
    default:
        return nil, errors.New("未知的 LLM 提供商: " + config.LLM.Provider)
    }
}