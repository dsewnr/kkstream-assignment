package actions

func (as *ActionSuite) Test_HomeHandler() {
	// with no user login
	res := as.HTML("/").Get()
	as.Equal(307, res.Code)
}
