package handler

import (
	dblayer "fileSystem/db"
	"fileSystem/util"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

/**
 * @Author: yirufeng
 * @Email: yirufeng@foxmail.com
 * @Date: 2020/9/16 9:36 上午
 * @Desc:
 */

const (
	pwdSalt = "*#890"
)

//处理用户注册的handler
func SignupHandler(w http.ResponseWriter, r *http.Request) {
	//如果为get方法
	if r.Method == http.MethodGet {
		data, err := ioutil.ReadFile("./static/view/signup.html")
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Write(data)
		return
	}

	//到达这里说明是post请求
	r.ParseForm()
	username := r.Form.Get("username")
	passwd := r.Form.Get("password")

	//做一下简单的校验，例如用户名长度小于3以及密码长度小于5
	if len(username) < 3 || len(passwd) < 5 {
		w.Write([]byte("Invalid parameter."))
		return
	}

	//在简单的校验之后，给用户密码进行sha1的加密
	//1.首先定义一个盐值
	//2.通过之前的util中的sha1来进行sha1的加密处理
	encPasswd := util.Sha1([]byte(passwd + pwd_salt))
	suc := dblayer.UserSignUp(username, encPasswd)
	//插入成功
	if suc {
		w.Write([]byte("SUCCESS"))
	} else {
		w.Write([]byte("FAILED"))
	}
}

//用户登录流程
func SignInHandler(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()

	//获取用户名和密码
	username := r.Form.Get("username")
	password := r.Form.Get("password")
	encPassword := util.Sha1([]byte(password + pwdSalt))

	//1. 校验用户名以及密码
	pwdChecked := dblayer.UserSignin(username, encPassword)
	if !pwdChecked {
		w.Write([]byte("FAILED"))
		return
	}

	//2. 校验密码通过之后生成一个访问的凭证（token）
	//token的规则是自定义的，我们自己定义生成一个40位的token
	token := GenToken(username)
	//下一步便是要将token写入到我们的数据库中
	upRet := dblayer.UpdateToken(username, token)
	if !upRet {
		w.Write([]byte("FAILED"))
		return
	}
	//3. 登录之后重定向到主页，交给浏览器去做
	//这是之前写的代码
	//w.Write([]byte("http://" + r.Host + "/static/view/home.html"))

	//现在我们使用了自己写的一个工具类封装了我们的json操作，并且由于返回的数据量比较大，所以我们推荐使用json作为返回的数据类型
	resp := util.RespMsg{
		Code: 0,
		Msg:  "OK",
		//因为我们登录成功要返回让用户重定向页面的url地址，并且还要给用户一个token用于进行其他api的凭证访问
		Data: struct {
			Location string
			Username string
			Token    string
		}{
			Location: "http://" + r.Host + "/static/view/home.html",
			Username: username,
			Token:    token,
		},
	}
	//最后我们将response转换成jsonBytes并返回给客户端
	w.Write(resp.JSONBytes())
}


//查询用户信息
func UserInfoHandler(w http.ResponseWriter, r *http.Request) {
	//1. 解析请求参数
	r.ParseForm()
	username := r.Form.Get("username")
	token := r.Form.Get("token")

	//2. 验证token是否有效：单独使用一个函数来完成
	//isValidToken := IsTokenValid(token)
	//if !isValidToken {  //如果token无效直接返回403
	//	w.WriteHeader(http.StatusInternalServerError)
	//	return
	//}

	//3. 查询用户信息，需要查询数据库，在db/user.go中添加一个方法GetUserInfo
	//同时为了方便方法的返回值，我们自己定义了一个User的结构体，与我们的user表一一对应
	user, err := dblayer.GetUserInfo(username)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	//4. 组装并且将用户数据作为响应发送回去
	resp := util.RespMsg{
		Code: 0,
		Msg: "OK",
		Data: user,
	}
	w.Write(resp.JSONBytes())
}

//生成token
func GenToken(username string) string {
	//生成40位字符：md5加密之后为32位，我们对这个字符串(用户名 拼接 时间戳 拼接 我们的token盐值)进行md5加密，此时有32位
	//后面的8位我们采用时间戳的前8位
	//md5(username + 时间戳 + token_salt) + timestamp[:8]
	ts := fmt.Sprintf("%X", time.Now().Unix())
	tokenPrefix := util.MD5([]byte(username + ts + "_tokensalt"))
	return tokenPrefix + ts[:8]
}


//token是否有效
func IsTokenValid(token string) bool {
	if len(token) != 40 {
		return false
	}
	//1. 判断token的时效性，也就是是否过期，取出后8位判断是否过了我们自己指定的时间(例如1天)，如果超过我们制定时间，说明就失效了

	//2. 从数据库表tbl_user_token查询username对应的token信息

	//3. 对比两个token是否一致，
}
