package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/go-resty/resty/v2"
)

type Service struct {
	apiKey string
	client *resty.Client
}

func NewService(apiKey string) *Service {
	return &Service{
		apiKey: apiKey,
		client: resty.New(),
	}
}

type OptimizeRequest struct {
	Title       string `json:"title"`
	Description string `json:"description"`
}

type OptimizeResponse struct {
	OptimizedTitle       string `json:"optimizedTitle"`
	OptimizedDescription string `json:"optimizedDescription"`
}

type glmRequest struct {
	Model    string       `json:"model"`
	Messages []glmMessage `json:"messages"`
}

type glmMessage struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type glmResponse struct {
	Choices []struct {
		Message struct {
			Content string `json:"content"`
		} `json:"message"`
	} `json:"choices"`
}

func (s *Service) OptimizeProduct(ctx context.Context, req OptimizeRequest) (*OptimizeResponse, error) {
	title := strings.TrimSpace(req.Title)
	description := strings.TrimSpace(req.Description)

	if title == "" && description == "" {
		return nil, fmt.Errorf("title or description is required")
	}

	prompt := fmt.Sprintf(`
请优化以下二手商品信息，使其更适合校园二手交易平台展示。

要求：
1. 标题简洁、真实、有吸引力。
2. 描述自然，不要夸张营销。
3. 适合大学校园二手交易场景。
4. 不要编造不存在的信息。
5. 只返回 JSON，不要返回 markdown，不要使用代码块。
6. 返回格式必须严格如下：

{
  "optimizedTitle": "...",
  "optimizedDescription": "..."
}

原始商品标题：%s
原始商品描述：%s
`, title, description)

	body := glmRequest{
		Model: "glm-4-flash",
		Messages: []glmMessage{
			{
				Role:    "user",
				Content: prompt,
			},
		},
	}

	var resp glmResponse

	httpResp, err := s.client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+s.apiKey).
		SetHeader("Content-Type", "application/json").
		SetBody(body).
		SetResult(&resp).
		Post("https://open.bigmodel.cn/api/paas/v4/chat/completions")

	if err != nil {
		return nil, fmt.Errorf("request zhipu failed: %w", err)
	}

	if httpResp.StatusCode() < 200 || httpResp.StatusCode() >= 300 {
		return nil, fmt.Errorf("zhipu status=%d body=%s", httpResp.StatusCode(), httpResp.String())
	}

	if len(resp.Choices) == 0 {
		return nil, fmt.Errorf("AI response empty, body=%s", httpResp.String())
	}

	content := strings.TrimSpace(resp.Choices[0].Message.Content)

	content = strings.TrimPrefix(content, "```json")
	content = strings.TrimPrefix(content, "```")
	content = strings.TrimSuffix(content, "```")
	content = strings.TrimSpace(content)

	var result OptimizeResponse
	if err := json.Unmarshal([]byte(content), &result); err != nil {
		return nil, fmt.Errorf("parse AI json failed: %w, raw=%s", err, resp.Choices[0].Message.Content)
	}

	if strings.TrimSpace(result.OptimizedTitle) == "" && strings.TrimSpace(result.OptimizedDescription) == "" {
		return nil, fmt.Errorf("AI result empty, raw=%s", resp.Choices[0].Message.Content)
	}

	return &result, nil
}
