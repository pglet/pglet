const { createProxyMiddleware } = require('http-proxy-middleware');

module.exports = app => {
  app.use(createProxyMiddleware("/ws", {target: "http://localhost:5000", ws: true}))
}