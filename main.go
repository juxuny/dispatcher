package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/juxuny/dispatcher/log"
	"net/http"
	"os"
	"os/exec"
	"path"
	"strings"
)

var (
	interpreter string
	workDir     string
	address     string
	tokenFile   string

	logger = log.NewLogger("[main]")
)

func parse() {
	flag.StringVar(&interpreter, "b", "/bin/bash", "script interpreter")
	flag.StringVar(&workDir, "w", "scripts", "working directory")
	flag.StringVar(&address, "l", ":8080", "http server listen address")
	flag.StringVar(&tokenFile, "tf", "token.list", "access token file")
	flag.Parse()
}

type ScriptHandler struct{}

func respData(w http.ResponseWriter, data map[string]interface{}, code ...int) {
	w.Header().Add("Content-Type", "application/json")
	httpStatus := http.StatusOK
	if len(code) > 0 {
		httpStatus = code[0]
	}
	w.WriteHeader(httpStatus)
	resp := map[string]interface{}{
		"code": httpStatus,
		"data": data,
	}
	respData, _ := json.Marshal(resp)
	_, _ = w.Write(respData)
}

func respMessage(w http.ResponseWriter, message string, code ...int) {
	w.Header().Add("Content-Type", "application/json")
	httpStatus := http.StatusOK
	if len(code) > 0 {
		httpStatus = code[0]
	}
	w.WriteHeader(httpStatus)
	resp := map[string]interface{}{
		"code":    httpStatus,
		"message": message,
	}
	respData, _ := json.Marshal(resp)
	_, _ = w.Write(respData)
}

func respError(w http.ResponseWriter, message string, code ...int) {
	w.Header().Add("Content-Type", "application/json")
	httpStatus := http.StatusBadRequest
	if len(code) > 0 {
		httpStatus = code[0]
	}
	w.WriteHeader(httpStatus)
	resp := map[string]interface{}{
		"code":  httpStatus,
		"error": message,
	}
	respData, _ := json.Marshal(resp)
	_, _ = w.Write(respData)
}

func (*ScriptHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		respError(w, fmt.Sprintf("unsupported http method: %s", r.Method))
		return
	}
	if err := r.ParseForm(); err != nil {
		logger.Error(err)
		respError(w, "parse form error")
		return
	}

	// get script file name
	fileName := r.URL.Path
	if strings.Index(fileName, "/scripts") != 0 {
		http.NotFound(w, r)
		return
	}
	fileName = strings.Replace(fileName, "/scripts", "", 1) // remove the prefix '/scripts'
	fileName = path.Join(workDir, fileName)
	logger.Debug("run scripts: " + fileName)
	if _, err := os.Stat(fileName); os.IsNotExist(err) {
		http.NotFound(w, r)
		return
	}
	env := make([]string, 0)
	for k, vs := range r.Form {
		if len(vs) > 0 {
			env = append(env, fmt.Sprintf("%s=%s", k, vs[0]))
		}
	}
	if tokenFile != "" {
		token := r.FormValue("token")
		if ok, err := checkToken(tokenFile, token); err != nil || !ok {
			respError(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
	}
	if len(env) > 0 {
		logger.Debug("environments: ", env)
	}

	// run script
	cmd := exec.Command(interpreter, fileName)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = env
	if err := cmd.Start(); err != nil {
		logger.Error()
		respError(w, "execute failed")
		return
	}

	respMessage(w, "OK")
}

func main() {
	parse()
	logger.Info("start http server " + address)
	if err := http.ListenAndServe(address, &ScriptHandler{}); err != nil {
		panic(err)
	}
}
