const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = app => {
  app.use(createProxyMiddleware("/ws", {target: "http://localhost:5000", ws: true, changeOrigin: true}))
  app.use(
    '/api',
    createProxyMiddleware({
      target: 'http://localhost:5000',
      changeOrigin: true,
    })
  ); 
}