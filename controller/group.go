package controller

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/QuantumNous/new-api/common"
	"github.com/QuantumNous/new-api/model"
	"github.com/QuantumNous/new-api/setting"
	"github.com/QuantumNous/new-api/setting/ratio_setting"

	"github.com/gin-gonic/gin"
)

func getCurrentUserGroup(c *gin.Context) (string, error) {
	userId := c.GetInt("id")
	user, err := model.GetUserCache(userId)
	if err != nil {
		return "", err
	}
	userGroup := strings.TrimSpace(user.Group)
	if userGroup == "" {
		userGroup = "default"
	}
	return userGroup, nil
}

func GetGroups(c *gin.Context) {
	groupNames := make([]string, 0)
	for groupName := range ratio_setting.GetGroupRatioCopy() {
		groupNames = append(groupNames, groupName)
	}
	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    groupNames,
	})
}

func GetUserGroups(c *gin.Context) {
	usableGroups := make(map[string]map[string]interface{})
	userGroup, err := getCurrentUserGroup(c)
	if err != nil {
		common.ApiError(c, err)
		return
	}

	if ratio_setting.ContainsGroupRatio(userGroup) {
		usableGroups[userGroup] = map[string]interface{}{
			"ratio": ratio_setting.GetGroupRatio(userGroup),
			"desc":  setting.GetUsableGroupDescription(userGroup),
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "",
		"data":    usableGroups,
	})
}

func normalizeTokenGroupForCurrentUser(c *gin.Context, tokenGroup string) (string, error) {
	userGroup, err := getCurrentUserGroup(c)
	if err != nil {
		return "", err
	}
	tokenGroup = strings.TrimSpace(tokenGroup)
	if tokenGroup == "" {
		return userGroup, nil
	}
	if tokenGroup != userGroup {
		return "", fmt.Errorf("无权创建 %s 分组的令牌", tokenGroup)
	}
	return tokenGroup, nil
}
