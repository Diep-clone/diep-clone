const webpack = require('webpack');
const path = require('path');

module.exports = {
  mode: 'development' /* or none or production */,
  entry: {
    core: [
      '@babel/polyfill',
      './src/js/index.js', 
      './src/styles/index.css'
    ]
  },

  output: {
    path: path.resolve(__dirname, 'dist/js')
    // filename: '',
    // publicPath: ''
  },

  module: {
    rules: [
      {
        test : /\.js$/,
        include: [
          path.resolve(__dirname, 'src/js')
        ],
        exclude: /node_modules/,
        use: {
          loader: 'babel-loader',
          options: {
            presets: ['@babel/preset-env'],
            plugins: ['@babel/plugin-proposal-class-properties']
          },
        }
      },
      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader'],
      }
    ]
  },  

  plugins: [],
  optimization: {},

  resolve: {
    modules: ['node_modules'],
    extensions: ['.js', '.json', '.css'],
  }
};