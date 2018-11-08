package http

import (
	"github.com/Waitfantasy/tmq/message/manager"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type api struct {
	gin     *gin.Engine
	manager *manager.Manager
}

func (a *api) register() {
	a.gin = gin.Default()
	g := a.gin.Group("/tmq")
	g.POST("/prepare", a.prepare())
	g.GET("/commit/:id", a.commit())
	g.GET("/rollback/:id", a.rollback())
	g.GET("/consumerCommit/:id")
}

func (a *api) prepare() gin.HandlerFunc {
	return func(c *gin.Context) {
		// topic
		topic := c.PostForm("topic")
		if topic == "" {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": "topic param is empty",
			})
			return
		}

		// retry second
		var retrySecond int
		if param := c.PostForm("retry_second"); param == "" {
			retrySecond = 0
		} else {
			if sec, err := strconv.Atoi(param); err != nil {
				c.JSON(http.StatusOK, map[string]interface{}{
					"success": false,
					"message": err.Error(),
				})
			} else {
				retrySecond = sec
			}
		}

		// body
		body := c.PostForm("body")

		// create prepare msg
		if msg, err := a.manager.Prepare(topic, retrySecond, body); err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
				"id":      msg.Id,
			})
		}
	}
}

func (a *api) commit() gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("id")
		if param == "" {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": "id param is empty",
			})
			return
		}

		id, err := strconv.ParseUint(param, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		if _, err = a.manager.Send(id); err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}
	}
}

func (a *api) rollback() gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("id")
		if param == "" {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": "id param is empty",
			})
			return
		}

		id, err := strconv.ParseUint(param, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		if _, err = a.manager.Cancel(id); err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}
	}
}

func (a *api) consumerCommit() gin.HandlerFunc {
	return func(c *gin.Context) {
		param := c.Param("id")
		if param == "" {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": "id param is empty",
			})
			return
		}

		id, err := strconv.ParseUint(param, 10, 64)
		if err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
			return
		}

		if _, err = a.manager.ConsumerCommit(id); err != nil {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": false,
				"message": err.Error(),
			})
		} else {
			c.JSON(http.StatusOK, map[string]interface{}{
				"success": true,
			})
		}
	}
}