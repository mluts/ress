var path = require('path');

module.exports = {
  entry: './index.js',
  output: {
    path: path.resolve(__dirname, 'static/assets/js'),
    filename: 'bundle.js'
  },

  module:  {
    loaders: {
      test: /\.js/,
      exclude: /node_modules/,
      loader: 'babel-loader',
      query: {
        presets: ['es2015', 'react']
      }
    }
  }
};
