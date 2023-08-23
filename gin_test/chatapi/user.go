package chatapi

import (
	"fmt"
	"gin_chat/chatdb"
	"gin_chat/models"
	"gin_chat/utils"
	"github.com/asaskevich/govalidator"
	"github.com/gin-gonic/gin"
	"math/rand"
	"net/http"
	"strconv"
)

// GetUserList
// @Tags 用户模块
// @Summary 获取用户列表
// @Success 200 {string} json{"code","message"}
// @Router /user/getUserList [get]
func GetUserList(c *gin.Context) {
	data := chatdb.GetUserList()

	c.JSON(http.StatusOK, gin.H{
		"code":    models.UserStatus.Normal,
		"message": "登录成功!",
		"data":    data,
	})
}

// CreateUser
// @Tags 用户模块
// @Summary 创建用户
// @param name query string false "用户名"
// @param password query string false "密码"
// @param repassword query string false "确认密码"
// @param phone query string false "手机号"
// @Success 200 {string} json{"code","message"}
// @Router /user/createUser [get]
func CreateUser(c *gin.Context) {
	user := models.UserBasic{}
	user.Name = c.Query("name")
	password := c.Query("password")
	rePassword := c.Query("repassword")

	salt := fmt.Sprintf("%06d", rand.Int31())
	user.Salt = salt

	u := chatdb.FindUserByName(user.Name)
	if u.Name != "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    models.UserStatus.ParamsError,
			"message": "用户名已注册",
			"data":    nil,
		})
		return
	}

	if password != rePassword {
		c.JSON(http.StatusOK, gin.H{
			"code":    models.UserStatus.ParamsError,
			"message": "两次密码不一致!",
			"data":    nil,
		})
		return
	}
	//user.Password = password
	user.Password = utils.MakePassword(password, salt)
	_, err := chatdb.CreateUser(&user)
	if err != nil {
		fmt.Printf("create user error:%v\n", err)
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"code":    models.UserStatus.Normal,
		"message": "新用户创建成功！",
		"data":    user,
	})
}

// DeleteUser
// @Tags 用户模块
// @Summary 删除用户
// @param id query string false "用户id"
// @Success 200 {string} json{"code","message"}
// @Router /user/deleteUser [get]
func DeleteUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.Query("id"))
	user.ID = uint(id)
	chatdb.DeleteUser(&user)
	c.JSON(http.StatusOK, gin.H{
		"code":    models.UserStatus.Normal,
		"message": "用户删除成功！",
		"data": gin.H{
			"id": user.ID,
		},
	})
}

// UpdateUser
// @Tags 用户模块
// @Summary 修改用户
// @param id formData string false "用户id"
// @param name formData string false "用户名"
// @param password formData string false "密码"
// @param phone formData string false "电话"
// @param email formData string false "email"
// @Success 200 {string} json{"code","message"}
// @Router /user/updateUser [post]
func UpdateUser(c *gin.Context) {
	user := models.UserBasic{}
	id, _ := strconv.Atoi(c.PostForm("id"))
	user.ID = uint(id)
	user.Name = c.PostForm("name")
	user.Password = c.PostForm("password")
	user.Phone = c.PostForm("phone")
	user.Email = c.PostForm("email")

	_, err := govalidator.ValidateStruct(&user)
	if err != nil {
		fmt.Printf("validate struct error:%v\n", err)
		c.JSON(http.StatusOK, gin.H{
			"code":    models.UserStatus.ParamsError,
			"message": "修改参数不匹配!",
			"data":    nil,
		})
		return
	}
	chatdb.UpdateUser(&user)
	c.JSON(http.StatusOK, gin.H{
		"code":    models.UserStatus.Normal,
		"message": "修改成功!",
		"data": gin.H{
			"id": user.ID,
		},
	})
}

// FindByUserNameAndPwd
// @Tags 用户模块
// @Summary 根据用户名和密码获取用户信息
// @param name query string false "用户名"
// @param password query string false "密码"
// @Success 200 {string} json{"code","message"}
// @Router /user/findByNameAndPwd [post]
func FindByUserNameAndPwd(c *gin.Context) {
	data := models.UserBasic{}

	name := c.Query("name")
	password := c.Query("password")

	user := chatdb.FindUserByName(name)
	if user.Name == "" {
		c.JSON(http.StatusOK, gin.H{
			"code":    models.UserStatus.ParamsError,
			"message": "该用户不存在",
			"data":    nil,
		})
		return
	}

	if user.Password != utils.MakePassword(password, user.Salt) {
		c.JSON(http.StatusOK, gin.H{
			"code":    models.UserStatus.ParamsError,
			"message": "用户名或密码不正确",
			"data":    data,
		})
		return
	}
	data = chatdb.FindUserByNameAndPwd(name, user.Password)

	c.JSON(http.StatusOK, gin.H{
		"code":    models.UserStatus.Normal,
		"message": "登录成功!",
		"data":    data,
	})
}
