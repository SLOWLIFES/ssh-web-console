package routers

import (
	"embed"
	"io/fs"
	"net/http"
	"os"

	"github.com/genshen/ssh-web-console/src/controllers"
	"github.com/genshen/ssh-web-console/src/controllers/files"
	"github.com/genshen/ssh-web-console/src/utils"
)

const (
	RunModeDev  = "dev"
	RunModeProd = "prod"
)

//go:embed build
var statikFS embed.FS

func getFileSystem() http.FileSystem {
	fsys, err := fs.Sub(statikFS, "build")
	if err != nil {
		return nil
	}
	return http.FS(fsys)
}

func Register() {
	// serve static files
	// In dev mode, resource files (for example /ssh/static/*) and views(fro example /index.html) are served separately.
	// In production mode, resource files and views are served by statikFS (for example /*).

	http.Handle(utils.Config.Prod.StaticPrefix, http.StripPrefix(utils.Config.Prod.StaticPrefix, http.FileServer(getFileSystem())))

	bct := utils.Config.SSH.BufferCheckerCycleTime
	// api
	http.HandleFunc("/api/signin", controllers.SignIn)
	http.HandleFunc("/api/sftp/upload", controllers.AuthPreChecker(files.FileUpload{}))
	http.HandleFunc("/api/sftp/ls", controllers.AuthPreChecker(files.List{}))
	http.HandleFunc("/api/sftp/dl", controllers.AuthPreChecker(files.Download{}))
	http.HandleFunc("/ws/ssh", controllers.AuthPreChecker(controllers.NewSSHWSHandle(bct)))
	http.HandleFunc("/ws/sftp", controllers.AuthPreChecker(files.SftpEstablish{}))
}

/*
 * disable directory index, code from https://groups.google.com/forum/#!topic/golang-nuts/bStLPdIVM6w
 */
type justFilesFilesystem struct {
	fs http.FileSystem
}

func (fs justFilesFilesystem) Open(name string) (http.File, error) {
	f, err := fs.fs.Open(name)
	if err != nil {
		return nil, err
	}
	return neuteredReaddirFile{f}, nil
}

type neuteredReaddirFile struct {
	http.File
}

func (f neuteredReaddirFile) Readdir(count int) ([]os.FileInfo, error) {
	return nil, nil
}
