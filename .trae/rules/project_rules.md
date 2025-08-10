1. 请保持对话语言为中文
2. 请在生成代码时添加函数级注释
3. 执行数据库相关命令的时候直接带上本地数据库密码123456789a,不需要用户再输入,数据库账号为root,数据库端口为3306,数据库名称为hotgo
4. 最后更新时间取当前时间,创建时间取当前时间
5. 产品账号密码为admin,123456,curl前需要先获取token,获取token的url为/admin/site/accountLogin,请求方式为POST,请求体为{"username":"admin","password":"M/IZq/iWn6mAgH1SuJuM5g==","cid":"","code":""}，响应为{
    "code": 0,
    "message": "操作成功",
    "data": {
        "id": 1,
        "username": "admin",
        "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpZCI6MSwicGlkIjowLCJkZXB0SWQiOjEwMCwiZGVwdFR5cGUiOiJjb21wYW55Iiwicm9sZUlkIjoxLCJyb2xlS2V5Ijoic3VwZXIiLCJ1c2VybmFtZSI6ImFkbWluIiwicmVhbE5hbWUiOiLlrZ_luIUiLCJhdmF0YXIiOiJodHRwczovL2dteWNvcy5mYWNtcy5jbi9ob3Rnby9hdHRhY2htZW50LzIwMjMtMDItMDkvY3FkcThlcjluZmtjaGRvcGF2LnBuZyIsImVtYWlsIjoiMTMzODE0MjUwQHFxLmNvbSIsIm1vYmlsZSI6IjE1MzAzODMwNTcxIiwiYXBwIjoiYWRtaW4iLCJsb2dpbkF0IjoiMjAyNS0wOC0xMCAyMDo0MzoxMyJ9.n2nGsUTqDk7A3tZU1jv2XE-Mrlcj9EyI6k3khHVnx5U",
        "expires": 604800
    },
    "timestamp": 1754829793,
    "traceID": "c852443c66685a18bdb21e7a53f0dc19"
}，curl需要添加 Authorization 头，值为响应体中的token