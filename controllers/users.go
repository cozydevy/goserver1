package controllers

import (
	"course-go/config"
	"course-go/models"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/copier"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
)

type Users struct {
	DB *gorm.DB
}

type createUserForm struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=8"`
	Name     string `json:"name" binding:"required"`
}

type updateUserForm struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Role     string `json:"role"`
}

type getUserForm struct {
	Email       string `json:"email"`
	Password    string `json:"password"`
	Newpassword string `json:"newpassword"`
	Name        string `json:"name"`
	Role        string `json:"role"`
}
type userResponse struct {
	ID    uint   `json:"id"`
	Email string `json:"email"`
	Name  string `json:"name"`
	Role  string `json:"role"`
}
type userPass struct {
	Password string `json:"password"`
}

type usersPaging struct {
	Items  []userResponse `json:"items"`
	Paging *pagingResult  `json:"paging"`
}

func (u *Users) FindAll(ctx *gin.Context) {

	var users []models.User
	query := u.DB.Order("id desc").Find(&users)

	pagination := pagination{ctx: ctx, query: query, records: &users}
	paging := pagination.paginate()

	var serializedUsers []userResponse
	copier.Copy(&serializedUsers, &users)
	ctx.JSON(http.StatusOK, gin.H{
		"users": usersPaging{Items: serializedUsers, Paging: paging},
	})
}

func (u *Users) FindOne(ctx *gin.Context) {

	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

func (u *Users) Create(ctx *gin.Context) {

	var form createUserForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var user models.User

	copier.Copy(&user, &form)
	user.Password = user.GenerateEncryptedPassword()

	if err := u.DB.Create(&user).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusCreated, gin.H{"user": serializedUser})
}
func (u *Users) UpdatebyAdmin(ctx *gin.Context) {
	var form updateUserForm
	if err := ctx.ShouldBindJSON(&form); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	if form.Password != "" {
		user.Password = user.GenerateEncryptedPassword()
	}

	if err := u.DB.Model(&user).Update(&form).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}
func (u *Users) Update(ctx *gin.Context) {
	var users models.User
	var getdata getUserForm

	var form updateUserForm
	id := ctx.Param("id")

	if err := ctx.ShouldBindJSON(&getdata); err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	// u.DB.Find(&users, "id = ?", id)

	db := config.GetDB()
	if db.Where("id = ?", id).First(&users).RecordNotFound() {
		fmt.Printf("Found error\n")

		return
	}
	if getdata.Email != users.Email {
		form.Email = getdata.Email
	}

	form.Name = getdata.Name
	form.Role = getdata.Role

	var user models.User
	if getdata.Password != "" || getdata.Newpassword != "" {

		if err := bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(getdata.Password)); err != nil {
			copier.Copy(&user, &form)

			fmt.Printf("not match\n")
			fmt.Println("Match:   ", users.Password+"\n")
			fmt.Println("Match:   ", getdata.Password)
			return
		}
		match := bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(getdata.Password))
		fmt.Println("Match:   ", match)

		fmt.Println("\nMatch:   ", users.Password+"\n")
		fmt.Println("Match:   ", getdata.Password)
		form.Password = getdata.Newpassword
		copier.Copy(&user, &form)

		form.Password = user.GenerateEncryptedPassword()
	} else {
		copier.Copy(&user, &form)

	}

	if err := u.DB.Model(&user).Where("id = ?", id).Update(&form).Error; err != nil {
		ctx.JSON(http.StatusUnprocessableEntity, gin.H{"error": err.Error()})
		return
	}

	var serializedUser userResponse
	copier.Copy(&serializedUser, &users)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}
func (u *Users) Delete(ctx *gin.Context) {

	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	u.DB.Unscoped().Delete(&user)

	ctx.Status(http.StatusNoContent)
}

func (u *Users) Promote(ctx *gin.Context) {

	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user.Promote()
	u.DB.Save(user)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

func (u *Users) Demote(ctx *gin.Context) {

	user, err := u.findUserByID(ctx)
	if err != nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	user.Demote()
	u.DB.Save(user)

	var serializedUser userResponse
	copier.Copy(&serializedUser, &user)
	ctx.JSON(http.StatusOK, gin.H{"user": serializedUser})
}

func (u *Users) findUserByID(ctx *gin.Context) (*models.User, error) {
	id := ctx.Param("id")
	var user models.User

	if err := u.DB.First(&user, id).Error; err != nil {
		return nil, err
	}

	return &user, nil
}

func setUserImage(ctx *gin.Context, user *models.User) error {
	file, _ := ctx.FormFile("avatar")
	if file == nil {
		return nil
	}

	if user.Avatar != "" {
		user.Avatar = strings.Replace(user.Avatar, os.Getenv("HOST"), "", 1)
		pwd, _ := os.Getwd()
		os.Remove(pwd + user.Avatar)
	}

	path := "uploads/users/" + strconv.Itoa(int(user.ID))
	os.MkdirAll(path, os.ModePerm)
	filename := path + "/" + file.Filename
	if err := ctx.SaveUploadedFile(file, filename); err != nil {
		return err
	}

	db := config.GetDB()
	user.Avatar = os.Getenv("HOST") + "/" + filename
	db.Save(user)

	return nil
}
