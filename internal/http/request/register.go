package request

import (
	"github.com/gin-gonic/gin"
)

type RegisterInput struct {
	DiscordAccountName string `json:"discordAccountName"`
	Password           string `json:"password"`
	Login              string `json:"login"`
}

func GetValidatedRegisterInputPayload(c *gin.Context) (*RegisterInput, error) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		return nil, err
	}

	return &input, nil
}
