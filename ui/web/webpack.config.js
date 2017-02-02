var path = require('path');

module.exports = {
  entry: './index.js',
  output: {
    path: path.resolve(__dirname, 'static/assets/js'),
    filename: 'bundle.js'
  }
};
