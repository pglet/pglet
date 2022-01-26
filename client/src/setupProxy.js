const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = app => {
  app.use(createProxyMiddleware("/ws", {target: "http://localhost:8550", ws: true, changeOrigin: true}))
  app.use(
    '/api',
    createProxyMiddleware({
      target: 'http://localhost:8550',
      changeOrigin: true,
    })
  ); 
}