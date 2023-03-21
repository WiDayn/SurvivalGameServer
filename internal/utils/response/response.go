package response

import "github.com/gin-gonic/gin"

func Write(c *gin.Context, text string) {
	// START 到 END 之间的内容就是信息
	_, err := c.Writer.Write([]byte("==START==\n" + text + "\n==END==\n"))
	if err != nil {
		return
	}
}
