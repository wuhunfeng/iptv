package main

import (
	"fmt"
	"io"
	"net/http"
	"os"
)

// 简单账号密码，自行修改
const (
	adminUser = "admin"
	adminPwd  = "admin123"
	dataPath  = "/data/source.txt"
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

	// 登录页面 GET / 提交 POST
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
<h4>拉取直播源</h4>
<form action="/fetch" method="post">
源地址：<input type="text" name="url" style="width:400px;" placeholder="https://xxx/xxx.m3u"><br><br>
<button type="submit">开始拉取并保存 source.txt</button>
</form>
<p><a href="/source.txt">查看/下载 source.txt</a></p>
<p><a href="/health">健康检查 /health</a></p>
<p><a href="/m3u">获取M3U订阅</a></p>
<p><a href="/">返回首页</a></p>
</body>
</html>`
		_, _ = w.Write([]byte(html))
	})

	// 拉取源并保存
	http.HandleFunc("/fetch", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.Write([]byte("仅支持POST提交"))
			return
		}
		url := r.PostFormValue("url")
		if url == "" {
			w.Write([]byte("URL不能为空 <br><a href='/dashboard'>返回面板</a>"))
			return
		}

		resp, err := http.Get(url)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("拉取失败：%v <br><a href='/dashboard'>返回面板</a>", err)))
			return
		}
		defer resp.Body.Close()

		// 读取内容
		content, err := io.ReadAll(resp.Body)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("读取内容失败：%v <br><a href='/dashboard'>返回面板</a>", err)))
			return
		}

		// 写入文件
		err = os.WriteFile(dataPath, content, 0644)
		if err != nil {
			w.Write([]byte(fmt.Sprintf("写入文件失败：%v <br><a href='/dashboard'>返回面板</a>", err)))
			return
		}

		w.Write([]byte("✅ 拉取成功，已保存到 source.txt <br><a href='/source.txt'>查看文件</a> | <a href='/dashboard'>返回面板</a>"))
	})

	// 对外提供source.txt访问
	http.HandleFunc("/source.txt", func(w http.ResponseWriter, r *http.Request) {
		content, err := os.ReadFile(dataPath)
		if err != nil {
			w.Write([]byte("文件不存在，请先拉取源"))
			return
		}
		w.Header().Set("Content-Type", "text/plain;charset=utf-8")
		_, _ = w.Write(content)
	})

	// m3u订阅输出接口（可以读取source.txt输出m3u，这里先保留原有示例）
	http.HandleFunc("/m3u", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/vnd.apple.mpegurl")
		m3uContent := `#EXTM3U
#EXTINF:-1 tvg-name="测试频道" tvg-logo="",测试频道
http://example.com/live.m3u8
`
		_, _ = w.Write([]byte(m3uContent))
	})

	fmt.Printf("✅ Server listen on %s\n", addr)
	err := http.ListenAndServe(addr, nil)
	if err != nil {
		fmt.Printf("❌ Start failed: %v\n", err)
		os.Exit(1)
	}
}
