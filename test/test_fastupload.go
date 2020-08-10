package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"path/filepath"

	fsUtil "filestore-server/util"

	jsonit "github.com/json-iterator/go"
)

var (
	// 测试用户名
	username string
	// 测试用户密码
	password string
	// 当前用于测试上传的本地文件路径
	uploadFilePath string
	// 秒传接口
	targetURL = "http://localhost:28080/file/fastupload"
	// 登录接口
	apiUserSignin = "http://localhost:8080/user/signin"
)

// 登录
func signin(username, password string) (token string, err error) {
	resp, err := http.PostForm(
		apiUserSignin,
		url.Values{
			"username": {username},
			"password": {password},
		})
	if err != nil {
		return "", err
	}
	log.Println(username + " " + password)

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil || string(body) == "FAILED" {
		return "", err
	}

	// log.Println(string(body))
	token = jsonit.Get(body, "data").Get("Token").ToString()
	log.Println("signin token: " + token)

	return token, nil
}

func testUpload(username, token, uploadFilePath, filehash string) {
	resp, err := http.PostForm(targetURL, url.Values{
		"username": {username},
		"token":    {token},
		"filehash": {filehash},
		"filename": {filepath.Base(uploadFilePath)},
	})
	if err != nil {
		log.Printf("request upload error: %+v\n", err)
	}

	if resp != nil {
		body, err := ioutil.ReadAll(resp.Body)
		log.Printf("uploadRespErr: %+v\n", err)
		if err == nil {
			log.Printf("uploadRespBody: %+v\n", string(body))
		}
	}
}

func main() {
	flag.StringVar(&uploadFilePath, "fpath", "", "指定上传的文件路径")
	flag.StringVar(&username, "user", "", "测试用户名")
	flag.StringVar(&password, "pwd", "", "测试用户密码")
	flag.Parse()

	// 0: 登录，获取token
	token, err := signin(username, password)
	if err != nil {
		log.Println(err)
		return
	} else if token == "" {
		log.Println("登录失败，请检查用户名、密码")
		return
	}
	log.Println("token: " + token)

	fileHash, err := fsUtil.ComputeSha1ByShell(uploadFilePath)
	if err != nil {
		fmt.Println(err)
		return
	}

	testUpload(username, token, uploadFilePath, fileHash)
}
