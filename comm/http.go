package comm

import (
	"fmt"
	"io"
	"net/http"
	"strings"

	"github.com/botuniverse/go-libonebot/utils"
	log "github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
)

type httpComm struct {
}

func (comm *httpComm) handleActionRequest(w http.ResponseWriter, r *http.Request) {
	log.Debugf("HTTP Request: %v", r)

	if r.Method != "POST" {
		log.Warn("Action 只支持通过 POST 方式请求")
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	} else if !strings.HasPrefix(r.Header.Get("Content-Type"), "application/json") {
		log.Warn("Action 请求体 MIME 类型必须是 application/json")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Warnf("获取 Action 请求体失败, 错误: %v", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	body := utils.BytesToString(bodyBytes)
	log.Debugf("HTTP Request Body: %v", body)
	if !gjson.Valid(body) {
		log.Warn("Action 请求体不是合法的 JSON")
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	actionRequest := gjson.Parse(body)
	log.Debugf("Action Request: %v", actionRequest)

	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.WriteHeader(http.StatusOK)
	w.Write(bodyBytes)
}

func (comm *httpComm) handleStatusPage(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("<h1>Everything is OK!</h1>"))
}

// Start an HTTP communication task.
func StartHTTPTask(host string, port uint16) {
	addr := fmt.Sprintf("%s:%d", host, port)
	log.Infof("正在启动 HTTP 通信方式, 监听地址: %v", addr)

	comm := &httpComm{}
	mux := http.NewServeMux()
	mux.HandleFunc("/status", comm.handleStatusPage)
	mux.HandleFunc("/", comm.handleActionRequest)

	go func() {
		if err := http.ListenAndServe(addr, mux); err != nil && err != http.ErrServerClosed {
			log.Errorf("HTTP 通信方式启动失败, 错误: %v", err)
		}
		log.Info("HTTP 通信方式已退出")
	}()
}
