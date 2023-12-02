package openapi

import (
	"log"
	"net/http"
	"path/filepath"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/74th/vscode-book-r2-golang/domain/entity"
	"github.com/74th/vscode-book-r2-golang/domain/repository"
	"github.com/74th/vscode-book-r2-golang/domain/usecase"
)

// サーバAPI
type Server struct {
	server         http.Server
	tokenValidator repository.TokenValidator
	interactor     *usecase.Interactor
}

// サーバAPIのインスタンスを作成する
func New(addr string, webroot string, interactor *usecase.Interactor, tokenValidator repository.TokenValidator) *Server {
	s := &Server{
		server: http.Server{
			Addr: addr,
		},
		interactor:     interactor,
		tokenValidator: tokenValidator,
	}

	s.setRouter(webroot)

	return s
}

// サーバを開始する
func (s *Server) Serve() error {
	errChan := make(chan error)
	timer := time.NewTimer(100 * time.Millisecond)
	go func() {
		if err := s.server.ListenAndServe(); err != http.ErrServerClosed {
			log.Fatalf("could not start server: %s", err.Error())
		}
	}()
	select {
	case err := <-errChan:
		return err
	case <-timer.C:
	}
	timer.Stop()
	return nil
}

func (s *Server) Close() error {
	return s.server.Close()
}

// ルータの設定
func (s *Server) setRouter(webroot string) {
	router := gin.Default()
	api := router.Group("/api")
	api.Use(s.tokenValidate)
	api.GET("/tasks", s.list)
	api.POST("/tasks", s.create)
	api.PATCH("/tasks/:id/done", s.done)

	router.StaticFile("/", filepath.Join(webroot, "index.html"))
	router.Static("/js", filepath.Join(webroot, "js"))
	router.Static("/css", filepath.Join(webroot, "css"))
	s.server.Handler = router
}

// GET /tasks
// タスク一覧
func (s *Server) list(c *gin.Context) {
	tasks, err := s.interactor.ShowTasks()
	if err != nil {
		log.Print("error", err)
		c.Status(500)
		return
	}

	c.JSON(http.StatusOK, tasks)
}

// POST /tasks
// タスクの追加
func (s *Server) create(c *gin.Context) {
	task := new(entity.Task)

	err := c.ShouldBindJSON(task)
	if err != nil {
		log.Print("deserialize error", err)
		c.Status(401)
		return
	}

	task, err = s.interactor.CreateTask(task)
	if err != nil {
		log.Print("error", err)
		c.Status(500)
		return
	}

	c.JSON(200, task)
}

// POST /tasks/:id/done
// タスク完了
func (s *Server) done(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.Status(400)
		return
	}

	task, err := s.interactor.DoneTask(id)
	if err != nil {
		c.Status(404)
		return
	}

	c.JSON(200, task)
}

func (s *Server) tokenValidate(c *gin.Context) {
	auth := c.Request.Header.Get("Authorization")
	if len(auth) == 0 {
		c.Status(401)
		c.Abort()
		return
	}

	ok, err := s.tokenValidator.Validate(c.Request.Context(), auth)
	if err != nil {
		c.Status(500)
		c.Abort()
		return
	}

	if !ok {
		c.Status(401)
		c.Abort()
		return
	}
}
