package controller

import (
	"database/sql"
	"errors"
	"github.com/gin-gonic/gin"
	"github.com/yuk-vg/vg-1day-2017-08-09/doppio/httputil"
	"github.com/yuk-vg/vg-1day-2017-08-09/doppio/model"
	"net/http"
	"strconv"
)

// Message is controller for requests to messages
type Message struct {
	DB     *sql.DB
	Stream chan *model.Message
}

// All は全てのメッセージを取得してJSONで返します
func (m *Message) All(c *gin.Context) {
	msgs, err := model.MessagesAll(m.DB)
	if err != nil {
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	if len(msgs) == 0 {
		c.JSON(http.StatusOK, make([]*model.Message, 0))
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": msgs,
		"error":  nil,
	})
}

// GetByID はパラメーターで受け取ったidのメッセージを取得してJSONで返します
func (m *Message) GetByID(c *gin.Context) {
	msg, err := model.MessageByID(m.DB, c.Param("id"))

	switch {
	case err == sql.ErrNoRows:
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusNotFound, resp)
		return
	case err != nil:
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"result": msg,
		"error":  nil,
	})
}

// Create は新しいメッセージ保存し、作成したメッセージをJSONで返します
func (m *Message) Create(c *gin.Context) {
	var msg model.Message

	if c.Request.ContentLength == 0 {
		resp := httputil.NewErrorResponse(errors.New("body is missing"))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := c.BindJSON(&msg); err != nil {
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// 1-2. ユーザー名を追加しよう
	// できる人は、ユーザー名が空だったら`anonymous`等適当なユーザー名で投稿するようにしてみよう

	inserted, err := msg.Insert(m.DB)
	if err != nil {
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	// bot対応
	m.Stream <- inserted

	c.JSON(http.StatusCreated, gin.H{
		"result": inserted,
		"error":  nil,
	})
}

// UpdateByID は...
func (m *Message) UpdateByID(c *gin.Context) {
	// 1-3. メッセージを編集しよう
	// ...
	var msg model.Message

	if c.Request.ContentLength == 0 {
		resp := httputil.NewErrorResponse(errors.New("body is missing"))
		c.JSON(http.StatusBadRequest, resp)
		return
	}

	if err := c.BindJSON(&msg); err != nil {
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	_, err := msg.Update(m.DB)
	if err != nil {
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}
	c.JSON(http.StatusCreated, gin.H{
		"error":  nil,
	})
}

// DeleteByID は...
func (m *Message) DeleteByID(c *gin.Context) {
	// 1-4. メッセージを削除しよう
	// ...
	var msg model.Message
	msg.ID, _ = strconv.ParseInt(c.Param("id"), 10, 64)

	err := msg.Delete(m.DB)
	if err != nil {
		resp := httputil.NewErrorResponse(err)
		c.JSON(http.StatusInternalServerError, resp)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"error": nil,
	})
}
