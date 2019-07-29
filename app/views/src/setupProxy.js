// setupProxy.js is a development tool
// react-scripts let's you configure the proxy yourself using this file
const proxy = require('http-proxy-middleware');

module.exports = function(app) {
  app.use(proxy('/auth', { target: 'http://localhost:4000/' }));
};
