在/root下面新增/web_go
所有檔案直接放在/root/web_go下面
web是一個編譯好的檔案
下面這三個資料夾需要被新增
product_picture
customer_picture
ad_picture
然後output.sql是database的表格
需要有
使用者: explosion
密碼: explosion
然後把output.sql匯入database: web(web是資料庫的名稱)

(以下是手動建立後端,如果上面的方法不行)
back.go是原始碼,include裡面的info.go是相關資訊
在/root/web_go的資料夾內
```
go mod init db.explosion.tw/web
go mod tidy
go get -u github.com/gin-gonic/gin
go get -u github.com/go-sql-driver/mysql
go get -u github.com/gin-contrib/cors
go get -u github.com/gin-contrib/sessions
go get -u github.com/gin-contrib/sessions/cookie
```
