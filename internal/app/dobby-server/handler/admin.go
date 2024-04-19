package handler

import (
	"github.com/labstack/echo/v4"
	"localdev/dobby-server/internal/app/dobby-server/model"
	"localdev/dobby-server/internal/app/dobby-server/view/crud"
	"localdev/dobby-server/internal/pkg/util"
	"strconv"
)

type AdminHandler struct {
	h *BaseHandler
}

func (a AdminHandler) HandleUserList(c echo.Context) error {
	userTable, err := a.h.UserApi.GetAllUser()
	var userList []model.UserCrud
	for _, user := range userTable {
		userList = append(userList, model.UserCrud{
			User:      user,
			EditUrl:   "/admin/user/" + strconv.Itoa(user.Id) + "/edit",
			ViewUrl:   "/admin/user/" + strconv.Itoa(user.Id),
			UpdateUrl: "/admin/user/" + strconv.Itoa(user.Id),
			DeleteUrl: "/admin/user/" + strconv.Itoa(user.Id),
		})
	}
	util.Panic(err)

	return render(c, view.UserList(*a.h.UserSession, *a.h.Tool, userList))
}

func (a AdminHandler) HandleUserEdit(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := a.h.UserApi.GetUserById(id)
	util.Panic(err)
	userCrud := model.UserCrud{
		User:      *user,
		EditUrl:   "/admin/user/" + strconv.Itoa(user.Id) + "/edit",
		ViewUrl:   "/admin/user/" + strconv.Itoa(user.Id),
		UpdateUrl: "/admin/user/" + strconv.Itoa(user.Id),
		DeleteUrl: "/admin/user/" + strconv.Itoa(user.Id),
	}
	util.Panic(err)

	return render(c, view.UserEdit(userCrud))
}

func (a AdminHandler) HandleUserUpdate(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := a.h.UserApi.GetUserById(id)
	util.Panic(err)
	user.Active = c.FormValue("active") == "on"
	user.Title = c.FormValue("title")
	var userPermissions []string
	for _, permission := range model.GetAllPermissions() {
		if c.FormValue(string(permission)) == "on" {
			userPermissions = append(userPermissions, string(permission))
		}
	}
	user.Permissions = model.PermissionsToString(userPermissions)
	user.Username = c.FormValue("username")
	err = a.h.UserApi.UpdateUser(id, user)
	util.Panic(err)

	userCrud := model.UserCrud{
		User:      *user,
		EditUrl:   "/admin/user/" + strconv.Itoa(user.Id) + "/edit",
		ViewUrl:   "/admin/user/" + strconv.Itoa(user.Id),
		UpdateUrl: "/admin/user/" + strconv.Itoa(user.Id),
		DeleteUrl: "/admin/user/" + strconv.Itoa(user.Id),
	}
	util.Panic(err)

	return render(c, view.UserView(*a.h.UserSession, userCrud))
}

func (a AdminHandler) HandleUserView(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := a.h.UserApi.GetUserById(id)
	util.Panic(err)

	userCrud := model.UserCrud{
		User:      *user,
		EditUrl:   "/admin/user/" + strconv.Itoa(user.Id) + "/edit",
		ViewUrl:   "/admin/user/" + strconv.Itoa(user.Id),
		UpdateUrl: "/admin/user/" + strconv.Itoa(user.Id),
		DeleteUrl: "/admin/user/" + strconv.Itoa(user.Id),
	}
	util.Panic(err)

	return render(c, view.UserView(*a.h.UserSession, userCrud))
}

func (a AdminHandler) HandleUserNewForm(c echo.Context) error {
	return render(c, view.UserNew())
}

func (a AdminHandler) HandleUserNew(c echo.Context) error {
	//insert new user
	var userPermissions []string
	for _, permission := range model.GetAllPermissions() {
		if c.FormValue(string(permission)) == "on" {
			userPermissions = append(userPermissions, string(permission))
		}
	}
	user := model.User{
		Username:    c.FormValue("username"),
		Active:      c.FormValue("active") == "on",
		Title:       c.FormValue("title"),
		Permissions: model.PermissionsToString(userPermissions),
	}
	err := a.h.UserApi.InsertUser(user)
	util.Panic(err)

	//Reload user list with the new user added
	userTable, err := a.h.UserApi.GetAllUser()
	var userList []model.UserCrud
	for _, user := range userTable {
		userList = append(userList, model.UserCrud{
			User:      user,
			EditUrl:   "/admin/user/" + strconv.Itoa(user.Id) + "/edit",
			ViewUrl:   "/admin/user/" + strconv.Itoa(user.Id),
			UpdateUrl: "/admin/user/" + strconv.Itoa(user.Id),
			DeleteUrl: "/admin/user/" + strconv.Itoa(user.Id),
		})
	}
	util.Panic(err)

	return render(c, view.UserList(*a.h.UserSession, *a.h.Tool, userList))
}

func (a AdminHandler) HandleUserDelete(c echo.Context) error {
	id, _ := strconv.Atoi(c.Param("id"))
	err := a.h.UserApi.DeleteUserById(id)
	util.Panic(err)

	//Reload user list with the user deleted
	userTable, err := a.h.UserApi.GetAllUser()
	var userList []model.UserCrud
	for _, user := range userTable {
		userList = append(userList, model.UserCrud{
			User:      user,
			EditUrl:   "/admin/user/" + strconv.Itoa(user.Id) + "/edit",
			ViewUrl:   "/admin/user/" + strconv.Itoa(user.Id),
			UpdateUrl: "/admin/user/" + strconv.Itoa(user.Id),
			DeleteUrl: "/admin/user/" + strconv.Itoa(user.Id),
		})
	}

	return render(c, view.UserList(*a.h.UserSession, *a.h.Tool, userList))
}
