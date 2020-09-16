package handler

import "net/http"

/**
 * @Author: yirufeng
 * @Email: yirufeng@foxmail.com
 * @Date: 2020/9/16 1:21 下午
 * @Desc:
 */

//HTTP请求拦截器：在执行目标函数之前执行
func HTTPInterceptor(h http.HandlerFunc) http.HandlerFunc{
	return http.HandlerFunc(
		func(w http.ResponseWriter, r *http.Request) {
			//执行安检的代码
			r.ParseForm()

			username := r.Form.Get("username")
			token := r.Form.Get("token")

			//校验username以及token是否有效
			if len(username) < 3 || !IsTokenValid(token) {
				w.WriteHeader(http.StatusForbidden)
				return
			}
			//如果有效就执行传入进来的形参的方法
			h(w, r)
		})
}