# vue element admin零碎学习笔记

## login产生的跨域的问题
前端发送post请求的时候如果发现是跨域的时候会先发送OPTIONS请求，服务器返回的是支持的格式，和允许的域等信息，最开始修改了axios.create来切换mock和真是的接口发现会产生跨域的问题，正确的做法应该是在devServer中中添加
```
    after: require('./mock/mock-server.js'),
    proxy: {
      [process.env.VUE_APP_BASE_API]:{
        target:'http://127.0.0.2:8012',
        changeOrigin:true
      },
    },
```
## post get的差异
get 主要是为了获取数据参数一般放在url里，这样就可能造成如果穿用户名密码会被记录在浏览器记录里，post则不会所以登录时候用post方式


