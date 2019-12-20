package actions

func (as *ActionSuite) Test_Admin_Login() {
	res := as.HTML("/admin/login").Get()
	as.Equal(200, res.Code)
}

func (as *ActionSuite) Test_Admin_Logout() {
	res := as.HTML("/admin/logout").Get()
	as.Equal(200, res.Code)
}
