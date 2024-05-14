package main

import (
	"database/sql"
	"fmt"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"encoding/hex"
	"mime/multipart"
	// "net/http"
	"crypto/sha256"
	"db.explosion.tw/web/include"
	"strconv"
	"time"
	_ "github.com/go-sql-driver/mysql"
)

type ad_picture struct {
	Ad_id        string `json:"adid"`
	Pic_location string `json:"adSrc"`
}
type product_preview struct {
	Product_id   string `json:"productId"`
	Pic_location string `json:"src"`
	Seller_name  string `json:"SupplierName"`
	Product_name string `json:"ItemName"`
	Price        string `json:"ItemSalePrice"`
	Item_sold    string `json:"HasBeenSold"`
}
type product_detail struct {
	Product_id   string `json:"productId"`
	Pic_location string `json:"src"`
	Seller_name  string `json:"SupplierName"`
	Product_name string `json:"ItemName"`
	Price        string `json:"ItemSalePrice"`
	Item_sold    string `json:"HasBeenSold"`
	Number       string `json:"quantity"`
}
type get_orders struct {
	Product_id string `json:"productId"`
	Number     string `json:"quantity"`
}
type order_records struct {
	Id           string `json:"idNumber"`
	Name         string `json:"customerName"`
	Seller_id    string `json:"supplierId"`
	Seller_name  string `json:"supplierName"`
	Product_name string `json:"productName"`
	Number       string `json:"quantity"`
	Unit_price   string `json:"unitPrice"`
	Order_date   string `json:"orderDate"`
}
type signup struct {
	Id       string                `json:"idNumber" form:"idNumber"`
	Address  string                `json:"address" form:"address"`
	Age      string                `json:"age" form:"age"`
	Name     string                `json:"customerName" form:"customerName"`
	Image    *multipart.FileHeader `json:"imageSrc" form:"imageSrc"`
	Job      string                `json:"occupation" form:"occupation"`
	Phone    string                `json:"phoneNumber" form:"phoneNumber"`
	Password string                `json:"password" form:"password"`
}
type show_restock struct {
	Seller_id     string `json:"idNumber" form:"idNumber"`
	Seller_name   string `json:"supplierName" form:"supplierName"`
	Product_image string `json:"src" form:"src"`
	Location      string `json:"location" form:"location"`
	Product_id    string `json:"productId" form:"productId"`
	Product_name  string `json:"ItemName" form:"ItemName"`
	Price         string `json:"unitPrice" form:"unitPrice"`
	Number        string `json:"quantity" form:"quantity"`
}
type create_restock struct {
	Seller_id     string                `json:"idNumber" form:"idNumber"`
	Seller_name   string                `json:"supplierName" form:"supplierName"`
	Product_image *multipart.FileHeader `json:"src" form:"src"`
	Location      string                `json:"location" form:"location"`
	Product_id    string                `json:"productId" form:"productId"`
	Product_name  string                `json:"ItemName" form:"ItemName"`
	Price         string                `json:"unitPrice" form:"unitPrice"`
	Number        string                `json:"quantity" form:"quantity"`
}
type customer_receivealbe struct {
	Id                string `json:"idNumber"`
	Name              string `json:"customerName"`
	Remaining_balance string `json:"amount"`
}
type customer_info struct {
	Id         string `json:"idNumber"`
	Address    string `json:"address"`
	Age        string `json:"age"`
	Name       string `json:"customerName"`
	Image      string `json:"imageSrc"`
	Job        string `json:"occupation"`
	Phone      string `json:"phoneNumber"`
	Join_date  string `json:"registrationDate"`
	Status     string `json:"status"`
	Password   string `json:"password"`
	Permission string `json:"permission"`
}
type user_info struct {
	Id       string `json:"idNumber"`
	Img_path string `json:"imageSrc"`
}
type login struct {
}

var db *sql.DB

func main() {
	r := gin.Default()
	r.SetTrustedProxies([]string{include.Ip})
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{"https://dbf.explosion.tw,http://127.0.0.1,http://172.17.0.1"},
		AllowMethods: []string{"GET", "POST", "PUT", "PATCH", "OPTIONS"},
		AllowHeaders: []string{"Authorization", "X-Requested-With", "Content-Type", "Upgrade", "Origin",
			"Connection", "Accept-Encoding", "Accept-Language", "Host", "Access-Control-Request-Method", "Access-Control-Request-Headers",
		},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return origin == "https://github.com"
		},
		MaxAge: 12 * time.Hour,
	}))

	store := cookie.NewStore([]byte(include.Session_password))
	r.Use(sessions.Sessions("session", store))
	db, _ = sql.Open("mysql", include.Db_path)
	r.POST("/",func(c *gin.Context) {
		c.String(200,"yes")
	})
	r.POST("/login", func(c *gin.Context) {
		session := sessions.Default(c)
		id := c.PostForm("idNumber")
		password := c.PostForm("password")
		rows, err := db.Query("SELECT password FROM customer WHERE id=?", id)
		checkErr(err)
		var result_password string
		rows.Next()
		rows.Scan(&result_password)
		rows.Close()
		if sha(password) == result_password {
			session.Set("id", id)
			session.Save()
			c.JSON(200, gin.H{
				"result": 1,
			})
		} else {
			c.JSON(200, gin.H{
				"result": 0,
			})
		}
	})

	r.POST("/signup", func(c *gin.Context) {
		var data signup
		if err := c.ShouldBind(&data); err != nil {
			fmt.Println(err)
		}
		rows, err := db.Query("SELECT id FROM customer WHERE id=?", data.Id)
		checkErr(err)
		var result_id string
		rows.Next()
		rows.Scan(&result_id)
		rows.Close()

		if result_id == "" {
			db.Exec("INSERT INTO customer (id,password,name,phone,address,age,job,image) VALUES (?,?,?,?,?,?,?,?)",
				data.Id, sha(data.Password), data.Name, data.Phone, data.Address, data.Age, data.Job, data.Id)
			fmt.Println(data.Image)
			dst := "/home/explosion/web_go/customer_picture/" + data.Id + ".jpg"
			c.SaveUploadedFile(data.Image, dst)
			c.String(200, "success")
		} else {
			c.String(200, "account exist")
		}
	})
	r.GET("/is_login", func(c *gin.Context) {
		session := sessions.Default(c)
		session_id := session.Get("id")
		if session_id != nil {
			id := session_id.(string)
			if login_status(id) {
				c.JSON(200, gin.H{
					"result": 1,
				})
			} else {
				c.JSON(200, gin.H{
					"result": 0,
				})
			}
		}else{
			c.JSON(200, gin.H{
				"result": 0,
			})
		}
	})
	r.GET("/logout", func(c *gin.Context) {
		session := sessions.Default(c)
		session.Delete("id")
		session.Save()
		c.String(200, "logout")
	})

	r.GET("/preview_product", func(c *gin.Context) {
		rows, err := db.Query("SELECT a.product_id,a.image,b.name,a.product_name,a.price,a.item_sold FROM product_record AS a JOIN customer AS b ON a.seller_id=b.id")
		checkErr(err)
		var array []product_preview
		var tmp product_preview
		for rows.Next() {
			rows.Scan(&tmp.Product_id, &tmp.Pic_location, &tmp.Seller_name, &tmp.Product_name, &tmp.Price, &tmp.Item_sold)
			tmp.Pic_location = include.Path + "product_icon/" + tmp.Pic_location
			array = append(array, tmp)
		}
		rows.Close()
		c.JSON(200, array)
	})
	r.POST("/show_product", func(c *gin.Context) {
		id := c.PostForm("productId")
		fmt.Println(id)
		rows, err := db.Query("SELECT a.product_id,a.image,b.name,a.product_name,a.price,a.item_sold,a.number FROM product_record AS a JOIN customer AS b ON a.seller_id=b.id WHERE a.product_id=?", id)
		checkErr(err)
		var array []product_detail
		var tmp product_detail
		for rows.Next() {
			rows.Scan(&tmp.Product_id, &tmp.Pic_location, &tmp.Seller_name, &tmp.Product_name, &tmp.Price, &tmp.Item_sold, &tmp.Number)
			tmp.Pic_location = include.Path + "product_icon/" + tmp.Pic_location
			array = append(array, tmp)
		}
		rows.Close()
		c.JSON(200, array)
	})
	r.POST("/get_order", func(c *gin.Context) {
		session := sessions.Default(c)
		id := session.Get("id")
		var order_list []get_orders
		if err := c.ShouldBindJSON(&order_list); err != nil {
			fmt.Println(err)
		}
		for _, data := range order_list {
			rows, _ := db.Query("SELECT item_sold,seller_id,number FROM product_record WHERE product_id=?", data.Product_id)
			rows.Next()
			var s_id string
			var sold int
			var number int
			rows.Scan(&sold, &s_id, &number)
			rows.Close()
			if(number > 0){
				result, err := db.Exec("INSERT INTO order_record (id,product_id,number,seller_id) VALUES (?,?,?,?)", id, data.Product_id, data.Number, s_id)
			checkErr(err)
			fmt.Println(result)
			tmp_num, _ := strconv.Atoi(data.Number)
			db.Exec("UPDATE product_record SET item_sold=?, number=? WHERE product_id=?", sold+tmp_num, number-tmp_num, data.Product_id)
			}
		}
	})
	r.GET("/show_restock", func(c *gin.Context) {
		rows, err := db.Query("SELECT a.seller_id,b.name,a.image,a.location,a.product_id,a.product_name,a.price,a.number FROM product_record AS a JOIN customer AS b ON a.seller_id=b.id;")
		checkErr(err)
		var array []show_restock
		var tmp show_restock
		for rows.Next() {
			rows.Scan(&tmp.Seller_id, &tmp.Seller_name, &tmp.Product_image, &tmp.Location, &tmp.Product_id, &tmp.Product_name, &tmp.Price, &tmp.Number)
			tmp.Product_image = include.Path + "product_icon/" + tmp.Product_image
			array = append(array, tmp)
		}
		rows.Close()
		c.JSON(200, array)
	})
	r.GET("/show_order", func(c *gin.Context) {
		rows, err := db.Query("SELECT a.id,b.name,a.seller_id,c.name,d.product_name,a.number,d.price,a.order_date FROM order_record AS a JOIN customer AS b ON a.id=b.id JOIN customer AS c ON a.seller_id=c.id JOIN product_record AS d ON a.product_id=d.product_id")
		checkErr(err)
		var array []order_records
		var tmp order_records
		for rows.Next() {
			rows.Scan(&tmp.Id, &tmp.Name, &tmp.Seller_id, &tmp.Seller_name, &tmp.Product_name, &tmp.Number, &tmp.Unit_price, &tmp.Order_date)
			array = append(array, tmp)
		}
		rows.Close()
		c.JSON(200, array)
	})
	r.GET("/seller_receviable", func(c *gin.Context) {
		session := sessions.Default(c)
		id := session.Get("id")
		rows, err := db.Query("SELECT a.id,b.name,SUM(a.number*c.price) FROM order_record AS a JOIN customer AS b ON a.id=b.id JOIN product_record as c ON a.product_id=c.product_id WHERE a.seller_id=? GROUP BY a.id", id)
		checkErr(err)
		var array []customer_receivealbe
		var tmp customer_receivealbe
		for rows.Next() {
			rows.Scan(&tmp.Id, &tmp.Name, &tmp.Remaining_balance)
			array = append(array, tmp)
		}
		rows.Close()
		c.JSON(200, array)
	})
	r.GET("/show_customer_info", func(c *gin.Context) {
		rows, err := db.Query("SELECT address,age,name,id,image,job,phone,join_date,premission FROM customer")
		checkErr(err)
		var array []customer_info
		var tmp customer_info
		for rows.Next() {
			rows.Scan(&tmp.Address, &tmp.Age, &tmp.Name, &tmp.Id, &tmp.Image, &tmp.Job, &tmp.Phone, &tmp.Join_date, &tmp.Permission)
			if tmp.Status == "0" {
				tmp.Status = "停用"
			}else{
				tmp.Status = "正常"
			}
			tmp.Image = include.Path + "icon/" + tmp.Image + ".jpg"
			array = append(array, tmp)
		}
		rows.Close()
		c.JSON(200, array)
	})
	r.PATCH("/update_customer_info", func(c *gin.Context) {
		var data customer_info
		if err := c.ShouldBindJSON(&data); err != nil {
			fmt.Println(err)
		}
		fmt.Println("status: ",data.Status)
		result, err := db.Exec("UPDATE customer SET password=?, name=?, phone=?, address=?, age=?, job=?, premission=?, status=? WHERE id=?",
			sha(data.Password), data.Name, data.Phone, data.Address, data.Age, data.Job, data.Permission, data.Status, data.Id)
		fmt.Println(result)
		checkErr(err)
	})
	r.POST("/create_restock", func(c *gin.Context) {
		var data create_restock
		if err := c.ShouldBind(&data); err != nil {
			fmt.Println(err)
		}
		result, err := db.Exec("INSERT INTO product_record (seller_id, product_name, image, price, number, location) VALUES (?,?,?,?,?,?)",
			data.Seller_id, data.Product_name, data.Product_image.Filename, data.Price, data.Number, data.Location)
		fmt.Println(result)
		fmt.Println(err)
		dst := "/root/web_go/product_picture/" + data.Product_image.Filename
		c.SaveUploadedFile(data.Product_image, dst)
		c.String(200, "success")

	})
	r.GET("/user_info", func(c *gin.Context) {
		session := sessions.Default(c)
		id := session.Get("id")
		rows, err := db.Query("SELECT image FROM customer WHERE id=?", id)
		checkErr(err)
		rows.Next()
		var result string
		rows.Scan(&result)
		result = include.Path + "icon/" + result + ".jpg"
		c.JSON(200, gin.H{
			"idNumber": id,
			"imageSrc": result,
		})
	})
	r.GET("/ad/:id", func(c *gin.Context) {
		id := c.Param("id")
		imagePath := "/root/web_go/ad_picture/" + id
		c.File(imagePath)
	})
	r.GET("/icon/:id", func(c *gin.Context) {
		id := c.Param("id")
		imagePath := "/root/web_go/customer_picture/" + id
		c.File(imagePath)
	})
	r.GET("/is_admin", func(c *gin.Context) {
		session := sessions.Default(c)
		id := session.Get("id")
		rows, err := db.Query("SELECT premission from customer where id=?", id)
		checkErr(err)
		rows.Next()
		var result string
		rows.Scan(&result)
		if result == "1" {
			c.JSON(200, gin.H{
				"result": 1,
			})
		} else {
			c.JSON(200, gin.H{
				"result": 0,
			})
		}
	})
	r.GET("/product_icon/:id", func(c *gin.Context) {
		id := c.Param("id")
		imagePath := "/root/web_go/product_picture/" + id
		c.File(imagePath)
	})

	r.Run(include.Ip)
}
func checkErr(err error) {
	if err != nil {
		fmt.Println("Error:", err)
		panic(err.Error())
	}
}
func sha(input string) string {
	h := sha256.New()
	h.Write([]byte(input))
	sha1_hash := hex.EncodeToString(h.Sum(nil))
	return sha1_hash
}

func login_status(id string) bool {
	rows, err := db.Query("SELECT id FROM customer WHERE id=?", id)
	checkErr(err)
	var result_id string
	rows.Next()
	rows.Scan(&result_id)
	rows.Close()
	if result_id != "" {
		return true
	} else {
		return false
	}
}
