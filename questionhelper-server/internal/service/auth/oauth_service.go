package auth

import (
	"questionhelper-server/internal/dto"
	"questionhelper-server/internal/repository/user"
)

// GetOAuthStatus 获取第三方绑定状态
func GetOAuthStatus(userID uint) (*dto.OAuthStatusResponse, error) {
	oauths, err := user.FindOAuthUsersByUserID(userID)
	if err != nil {
		return nil, err
	}

	status := &dto.OAuthStatusResponse{}
	for _, oauth := range oauths {
		providerStatus := &dto.OAuthProviderStatus{
			Bound:    true,
			Nickname: oauth.ProviderUsername,
			BoundAt:  oauth.CreatedAt.Format("2006-01-02 15:04:05"),
		}

		switch oauth.Provider {
		case "wechat":
			status.Wechat = providerStatus
		case "github":
			status.Github = providerStatus
		case "google":
			status.Google = providerStatus
		}
	}

	return status, nil
}
