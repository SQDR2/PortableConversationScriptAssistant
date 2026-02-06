package services

import (
	"context"
	"fmt"
	"sidekick/backend/utils"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/errors"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	tmt "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/tmt/v20180321"
	"github.com/wailsapp/wails/v2/pkg/runtime"
)

type TranslationService struct {
	ctx    context.Context
	client *tmt.Client
	config *utils.AppConfig
}

func NewTranslationService() *TranslationService {
	return &TranslationService{}
}

func (s *TranslationService) Startup(ctx context.Context) {
	s.ctx = ctx
	// Load config
	// Assuming config.json is in the same directory as the executable or CWD
	// For dev, it's usually CWD.
	configPath := "config.json"
	config, err := utils.LoadConfig(configPath)
	if err != nil {
		// Try to find it relative to executable if not found in CWD
		// But for now, just log warning
		runtime.LogErrorf(s.ctx, "Failed to load config.json: %v. Translation feature will be disabled.", err)
		return
	}
	s.config = config

	if s.config.TencentCloud.SecretId == "" || s.config.TencentCloud.SecretId == "REPLACE_WITH_YOUR_SECRET_ID" {
		runtime.LogWarningf(s.ctx, "Tencent Cloud credentials not configured in config.json")
		return
	}

	credential := common.NewCredential(
		s.config.TencentCloud.SecretId,
		s.config.TencentCloud.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tmt.tencentcloudapi.com"
	client, err := tmt.NewClient(credential, s.config.TencentCloud.Region, cpf)
	if err != nil {
		runtime.LogErrorf(s.ctx, "Failed to create TMT client: %v", err)
		return
	}
	s.client = client
	runtime.LogInfof(s.ctx, "TranslationService initialized using region: %s", s.config.TencentCloud.Region)
}

func (s *TranslationService) GetConfigStatus() bool {
	if s.config == nil {
		return false
	}
	if s.config.TencentCloud.SecretId == "" || s.config.TencentCloud.SecretId == "REPLACE_WITH_YOUR_SECRET_ID" {
		return false
	}
	if s.config.TencentCloud.SecretKey == "" || s.config.TencentCloud.SecretKey == "REPLACE_WITH_YOUR_SECRET_KEY" {
		return false
	}
	return true
}

func (s *TranslationService) UpdateCredentials(secretId string, secretKey string) error {
	if s.config == nil {
		// Try to reload or init default if nil, but Startup should have handled it.
		// If Startup failed completely (file read error), maybe we need to create a new struct.
		// But assuming config.json exists (per requirement) but is empty.
		// If s.config is nil, it means LoadConfig failed entirely.
		// Let's create a default AppConfig if needed, but preserve existing one if possible.
		s.config = &utils.AppConfig{
			TencentCloud: utils.TencentCloudConfig{
				Region: "ap-guangzhou", // Default fallback
			},
		}
	}

	s.config.TencentCloud.SecretId = secretId
	s.config.TencentCloud.SecretKey = secretKey

	// Save to file
	configPath := "config.json"
	err := utils.SaveConfig(configPath, s.config)
	if err != nil {
		runtime.LogErrorf(s.ctx, "Failed to save config: %v", err)
		return fmt.Errorf("failed to save config: %w", err)
	}

	// Re-initialize client
	credential := common.NewCredential(
		s.config.TencentCloud.SecretId,
		s.config.TencentCloud.SecretKey,
	)
	cpf := profile.NewClientProfile()
	cpf.HttpProfile.Endpoint = "tmt.tencentcloudapi.com"
	client, err := tmt.NewClient(credential, s.config.TencentCloud.Region, cpf)
	if err != nil {
		runtime.LogErrorf(s.ctx, "Failed to create TMT client: %v", err)
		return fmt.Errorf("failed to create TMT client: %w", err)
	}
	s.client = client
	runtime.LogInfof(s.ctx, "TranslationService updated and re-initialized.")
	return nil
}

func (s *TranslationService) Translate(sourceText string, sourceLang string, targetLang string) (string, error) {
	if s.client == nil {
		return "", fmt.Errorf("Translation service not configured. Please check config.json")
	}

	request := tmt.NewTextTranslateRequest()
	request.SourceText = common.StringPtr(sourceText)
	request.Source = common.StringPtr(sourceLang)
	request.Target = common.StringPtr(targetLang)
	request.ProjectId = common.Int64Ptr(s.config.TencentCloud.ProjectId)

	response, err := s.client.TextTranslate(request)
	if err != nil {
		if _, ok := err.(*errors.TencentCloudSDKError); ok {
			return "", fmt.Errorf("Tencent Cloud API Error: %v", err)
		}
		return "", err
	}

	if response.Response.TargetText == nil {
		return "", fmt.Errorf("Received empty translation response")
	}

	return *response.Response.TargetText, nil
}
