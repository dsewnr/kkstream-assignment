package actions

import (
	"log"
	"os"
	"path/filepath"

	"github.com/gobuffalo/httptest"
)

func (as *ActionSuite) Test_Api_Upload() {

	testFile := "video0.mp4"
	testFilePath := filepath.Join("..", "testfiles", testFile)
	uploadFilePath := filepath.Join(UPLOAD_DIR, testFile)
	os.RemoveAll(UPLOAD_DIR)

	r, err := os.Open(testFilePath)
	defer r.Close()
	as.NoError(err)

	// setup a new httptest.File to hold the file information
	f := httptest.File{
		// ParamName is the name of the form parameter
		ParamName: "video",
		// FileName is the name of the file being uploaded
		FileName: testFile,
		// Reader is the file that is to be uploaded, any io.Reader works
		Reader: r,
	}

	// Post the file to /api/upload
	res, err := as.HTML("/api/upload").MultiPartPost(struct{}{}, f)
	as.NoError(err)
	as.Equal(200, res.Code)

	// assert the file exists on disk
	_, err = os.Stat(uploadFilePath)
	as.NoError(err)
}

func (as *ActionSuite) Test_Api_UploadMulti() {
	testFiles := []string{"video0.mp4", "video1.mp4", "video2.mp4"}

	os.RemoveAll(UPLOAD_DIR)

	var fs []httptest.File

	for _, testFile := range testFiles {
		testFilePath := filepath.Join("..", "testfiles", testFile)
		r, err := os.Open(testFilePath)
		defer r.Close()
		as.NoError(err)

		f := httptest.File{
			// ParamName is the name of the form parameter
			ParamName: "video",
			// FileName is the name of the file being uploaded
			FileName: testFile,
			// Reader is the file that is to be uploaded, any io.Reader works
			Reader: r,
		}

		fs = append(fs, f)
	}

	// Post the file to /api/uploadMulti
	res, err := as.HTML("/api/uploadMulti").MultiPartPost(struct{}{}, fs...)
	as.NoError(err)
	as.Equal(200, res.Code)

	// assert the file exists on disk
	for _, testFile := range testFiles {
		uploadFilePath := filepath.Join(UPLOAD_DIR, testFile)
		_, err := os.Stat(uploadFilePath)
		as.NoError(err)
	}
}

func (as *ActionSuite) Test_Api_UploadGcs() {

	testFile := "video0.mp4"
	testFilePath := filepath.Join("..", "testfiles", testFile)

	r, err := os.Open(testFilePath)
	defer r.Close()
	as.NoError(err)

	// setup a new httptest.File to hold the file information
	f := httptest.File{
		// ParamName is the name of the form parameter
		ParamName: "video",
		// FileName is the name of the file being uploaded
		FileName: testFile,
		// Reader is the file that is to be uploaded, any io.Reader works
		Reader: r,
	}

	// Post the file to /api/uploadGcs
	res, err := as.HTML("/api/uploadGcs").MultiPartPost(struct{}{}, f)
	log.Println(res.Result())
	as.NoError(err)
	as.Equal(200, res.Code)
}
