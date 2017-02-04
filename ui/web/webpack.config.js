var path = require('path');

module.exports = {
  entry: './webui/js/index.js',
  output: {
    filename: 'app.js',
    path: path.resolve(__dirname, "static/assets/js")
  }
};
