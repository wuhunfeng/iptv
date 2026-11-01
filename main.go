package main

import (
	"fmt"
	"net/http"
	"os"
)

// 简单账号密码，你可以自行修改
const (
	adminUser = "admin"
	adminPwd  = "admin123"
)

func main() {
	// 读取botkeep自动注入的PORT环境变量
	port := os.Getenv("PORT")
	if port == "" {
		port = "57516"
	}
	addr := fmt.Sprintf("0.0.0.0:%s", port)

	// 健康检查路由 /health，用来保活防休眠
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// 首页
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write([]byte("IPTV Control Center Go Version | Running OK <br><a href='/login'>前往登录</a>"))
	})

	// 登录页面 GET
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "text/html;charset=utf-8")
			html := `
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>登录</title>
</head>
<body>
<h3>IPTV 后台登录</h3>
<form method="post">
账号：<input type="text" name="user"><br><br>
密码：<input type="password" name="pwd"><br><br>
<button type="submit">登录</button>
</form>
</body>
</html>`
			_, _ = w.Write([]byte(html))
			return
		}

		// 登录提交 POST
		user := r.PostFormValue("user")
		pwd := r.PostFormValue("pwd")
		if user == adminUser && pwd == adminPwd {
			http.Redirect(w, r, "/dashboard", http.StatusFound)
			return
		}
		w.Write([]byte("账号或密码错误！<br><a href='/login'>返回登录</a>"))
	})

	// 后台面板 dashboard
	http.HandleFunc("/dashboard", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html;charset=utf-8")
		html := `
<!DOCTYPE html>
<html>
<head>
<meta charset="utf-8">
<title>IPTV Dashboard</title>
</head>
<body>
<h2>✅ IPTV Control Center 后台面板</h2>
<p>当前服务运行正常</p>
<p><a href="/health">健康检查 /health</a></p>
<p><a href="/">返回首页</a></p>
</body>
</html>`
		_, _ = w.Write([]byte(html))
	})

	fmt.Printf("✅ Server listen on %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("❌ Start failed: %v\n", err)
		os.Exit(1)
	}
}
