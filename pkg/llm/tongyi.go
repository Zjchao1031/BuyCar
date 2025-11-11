package llm

import (
    "bytes"
    "context"
    "encoding/json"
    "errors"
    "io"
    "net/http"
    "time"
)

// TongyiOptions 用于配置通义千问客户端
type TongyiOptions struct {
    ApiKey         string
    Model          string
    Endpoint       string
    TimeoutSeconds int
}

type tongyiClient struct {
    apiKey   string
    model    string
    endpoint string
    http     *http.Client
}

// NewTongyiClient 创建通义千问客户端
func NewTongyiClient(opt TongyiOptions) (Client, error) {
    if opt.ApiKey == "" {
        return nil, errors.New("Tongyi API Key 不能为空")
    }
    if opt.Model == "" {
        opt.Model = "qwen-plus"
    }
    if opt.Endpoint == "" {
        // 采用文本生成通用端点，输入为纯 prompt 文本
        opt.Endpoint = "https://dashscope.aliyuncs.com/api/v1/services/aigc/text-generation/generation"
    }
    if opt.TimeoutSeconds <= 0 {
        opt.TimeoutSeconds = 20
    }
    return &tongyiClient{
        apiKey:   opt.ApiKey,
        model:    opt.Model,
        endpoint: opt.Endpoint,
        http: &http.Client{Timeout: time.Duration(opt.TimeoutSeconds) * time.Second},
    }, nil
}

// tongyiRequest 与 tongyiResponse 的结构兼容文本生成接口与部分历史返回结构
type tongyiRequest struct {
    Model      string                 `json:"model"`
    Input      string                 `json:"input"`
    Parameters map[string]interface{} `json:"parameters,omitempty"`
}

type tongyiResponse struct {
    RequestID string          `json:"request_id"`
    Output    json.RawMessage `json:"output"`
    // 某些版本会直接返回 output_text
    OutputText string `json:"output_text"`
}

// Generate 调用通义千问，返回主体文本与原始 JSON
func (c *tongyiClient) Generate(ctx context.Context, prompt string) (string, string, string, error) {
    reqBody := tongyiRequest{
        Model: c.model,
        Input: prompt,
        Parameters: map[string]interface{}{
            // 要求返回简洁文本以便直接展示或存储
            "result_format": "text",
        },
    }
    b, _ := json.Marshal(reqBody)

    httpReq, err := http.NewRequestWithContext(ctx, http.MethodPost, c.endpoint, bytes.NewReader(b))
    if err != nil {
        return "", "", c.model, err
    }
    httpReq.Header.Set("Content-Type", "application/json")
    httpReq.Header.Set("Authorization", "Bearer "+c.apiKey)

    resp, err := c.http.Do(httpReq)
    if err != nil {
        return "", "", c.model, err
    }
    defer resp.Body.Close()
    raw, _ := io.ReadAll(resp.Body)
    if resp.StatusCode < 200 || resp.StatusCode >= 300 {
        return "", string(raw), c.model, errors.New("Tongyi API 调用失败: " + resp.Status)
    }

    var tr tongyiResponse
    if err := json.Unmarshal(raw, &tr); err != nil {
        // 返回原始 JSON，便于排障
        return "", string(raw), c.model, err
    }

    // 兼容解析：优先 output_text，其次 output.message/content/text
    if tr.OutputText != "" {
        return tr.OutputText, string(raw), c.model, nil
    }

    // 尝试解析 Output 内的可能结构
    // 1) {"text":"..."}
    var generic map[string]interface{}
    if len(tr.Output) > 0 {
        if err := json.Unmarshal(tr.Output, &generic); err == nil {
            if txt, ok := generic["text"].(string); ok && txt != "" {
                return txt, string(raw), c.model, nil
            }
            // 2) choices[0].message.content[0].text 风格
            if choices, ok := generic["choices"].([]interface{}); ok && len(choices) > 0 {
                if choice, ok := choices[0].(map[string]interface{}); ok {
                    if msg, ok := choice["message"].(map[string]interface{}); ok {
                        if contentArr, ok := msg["content"].([]interface{}); ok && len(contentArr) > 0 {
                            if item0, ok := contentArr[0].(map[string]interface{}); ok {
                                if txt, ok := item0["text"].(string); ok && txt != "" {
                                    return txt, string(raw), c.model, nil
                                }
                            }
                        }
                    }
                }
            }
        }
    }

    // 若无法解析主体文本，仍返回原始响应
    return "", string(raw), c.model, nil
}