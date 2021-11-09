const path = require('path')
const webpack = require('webpack')

module.exports = env => {
  return {
    plugins: [new webpack.DefinePlugin({})],
    entry: path.resolve(__dirname, 'src') + '/index.js',
    output: {
      path: path.resolve(__dirname, 'dist'),
      filename: 'bundle.js'
    },
    module: {
      rules: [
        {
          test: /\.js$/,
          exclude: /node_modules/,
          use: {
            loader: 'babel-loader'
          }
        }
      ]
    }
  }
}
