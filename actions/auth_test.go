package actions

func (as *ActionSuite) Test_Auth_Login() {
	res := as.HTML("/auth/login").Get()
	as.Equal(200, res.Code)
}

func (as *ActionSuite) Test_Auth_CallbackLine() {
	as.Fail("Not Implemented!")
}

func (as *ActionSuite) Test_Auth_DoLogin() {
	as.Fail("Not Implemented!")
}


func (as *ActionSuite) Test_Auth_Logout() {
	as.Fail("Not Implemented!")
}

