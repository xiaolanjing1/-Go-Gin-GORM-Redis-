package controller

import (
	"net/http"
	"os"
	"strconv"
	"time"
	"user-system/config"
	"user-system/middleware"
	"user-system/model"

	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

func UserList(c *gin.Context) {
	var users []model.User
	page := c.Query("page")
	size := c.Query("size")
	pageint, err := strconv.Atoi(page)
	if err != nil {
		pageint = 1
	}
	sizeint, err := strconv.Atoi(size)
	if err != nil {
		sizeint = 10
	}
	outpage := (pageint - 1) * sizeint
	err = config.Db.Limit(sizeint).Offset(outpage).Find(&users).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "请求失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"page":  pageint,
		"size":  sizeint,
		"users": users,
	})
}

func Delavatar(c *gin.Context) {
	username, _ := c.Get("username")
	var user model.User
	err := config.Db.Where("name=?", username).First(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "没有该用户",
		})
		return
	}
	if user.Avatar == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "该用户无上传头像",
		})
		return
	}
	path := "." + user.Avatar
	err = os.Remove(path)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "删除失败",
		})
		return
	}
	err = config.Db.Model(&user).Update("avatar", "").Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "更新失败",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"msg": "删除头像成功",
	})
}

func Updateemial(c *gin.Context) {
	username, _ := c.Get("username")
	var req struct {
		New_email string `json:"new_email"`
		Old_email string `json:"old_email"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}
	var user model.User
	err = config.Db.Where("name=?", username).First(&user).Error
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "该用户不存在",
		})
		return
	}
	if user.Email != req.Old_email {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "邮箱错误",
		})
		return
	}
	user.Email = req.New_email
	config.Db.Model(&user).Update("email", req.New_email)
	c.JSON(http.StatusOK, gin.H{
		"msg": "修改邮箱成功",
	})
}

func Logout(c *gin.Context) {
	token := c.GetHeader("Authorization")
	config.RDB.Del(
		config.Ctx,
		token,
	)
	c.JSON(http.StatusOK, gin.H{
		"msg": "退出登录成功",
	})

}

func Register(c *gin.Context) {
	var user model.User
	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "密码加密失败",
		})
		return
	}
	user.Password = string(hashPassword)
	config.Db.Create(&user)
	c.JSON(http.StatusOK, gin.H{
		"msg": "用户注册成功",
	})
}

func Login(c *gin.Context) {
	var res model.User
	err := c.ShouldBindJSON(&res)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}
	var user model.User
	config.Db.Where("name=?", res.Name).First(&user)
	if user.Id == 0 {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "没有该用户",
		})
		return
	}
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(res.Password),
	)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "密码错误",
		})
		return
	}
	if user.Email != res.Email {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "邮箱有误",
		})
		return
	}
	token, err := middleware.GenerateToken(user.Name)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "token生成失败",
		})
	}
	config.RDB.Set(
		config.Ctx,
		token,
		user.Name,
		time.Hour*24,
	)
	c.JSON(http.StatusOK, gin.H{
		"msg":   "成功登录",
		"token": token,
	})
}

func UserInfo(c *gin.Context) {
	username, _ := c.Get("username")
	c.JSON(http.StatusOK, gin.H{
		"username": username,
	})
}

func UserProfile(c *gin.Context) {
	username, _ := c.Get("username")
	var user model.User
	config.Db.Where("name=?", username).First(&user)
	c.JSON(http.StatusOK, gin.H{
		"name":     user.Name,
		"password": user.Password,
		"email":    user.Email,
	})
}

func UpdateUser(c *gin.Context) {
	username, _ := c.Get("username")

	var req struct {
		Name  string `json:"name"`
		Email string `json:"email"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}
	var user model.User
	config.Db.Where("name=?", username).First(&user)
	user.Name = req.Name
	user.Email = req.Email
	config.Db.Save(&user)
	c.JSON(http.StatusOK, gin.H{
		"msg": "修改成功",
	})
}

func UpdatePassword(c *gin.Context) {
	username, _ := c.Get("username")

	var req struct {
		NewPassword string `json:"new_password"`
		OldPassword string `json:"old_password"`
	}
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "参数错误",
		})
		return
	}
	var user model.User
	config.Db.Where("name=?", username).First(&user)
	err = bcrypt.CompareHashAndPassword(
		[]byte(user.Password),
		[]byte(req.OldPassword),
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "旧密码错误",
		})
		return
	}
	hashPassword, err := bcrypt.GenerateFromPassword([]byte(req.NewPassword), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "加密失败",
		})
		return
	}
	user.Password = string(hashPassword)
	config.Db.Save(&user)
	c.JSON(http.StatusOK, gin.H{
		"msg": "密码修改成功",
	})
}

func UploadAvatar(c *gin.Context) {
	username, _ := c.Get("username")
	file, err := c.FormFile("avatar")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "上传失败",
		})
		return
	}
	path := "./uploads/" + file.Filename
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}
	var user model.User
	config.Db.Where("name=?", username).First(&user)
	user.Avatar = "/uploads/" + file.Filename
	config.Db.Save(&user)
	c.JSON(http.StatusOK, gin.H{
		"msg":    "上传成功",
		"avatar": user.Avatar,
	})
}
