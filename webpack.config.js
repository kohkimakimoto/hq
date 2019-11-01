const webpack = require('webpack');
const path = require('path');
const MiniCssExtractPlugin = require('mini-css-extract-plugin');
const CopyPlugin = require('copy-webpack-plugin');

module.exports = {
  entry: {
    main: './ui/index.tsx'
  },
  output: {
    path: path.join(__dirname, '/res/ui/assets'),
    filename: 'bundle.js'
  },
  module: {
    rules: [
      {
        test: /\.tsx?$/,
        use: 'ts-loader'
      },
      {
        test: /\.(ttf|eot|woff|woff2|svg)$/,
        loader: 'file-loader',
        options: {
          name: '[name].[ext]'
        }
      },
      {
        test: /\.(png|jpe?g|gif)(\?.*)?$/,
        loader: 'url-loader',
        query: {
          limit: 10000,
          name: 'img/[name].[ext]'
        }
      },
      {
        test: /\.scss$/,
        use: [
          MiniCssExtractPlugin.loader,
          {
            loader: 'css-loader'
          },
          {
            loader: 'sass-loader',
            options: {
              sourceMap: true
            }
          }
        ]
      },
      {
        test: /\.css$/,
        use: [MiniCssExtractPlugin.loader, 'css-loader']
      }
    ]
  },
  optimization: {
    splitChunks: {
      name: 'vendor',
      chunks: 'initial'
    }
  },
  plugins: [
    new webpack.DefinePlugin({
      'process.env.NODE_ENV': JSON.stringify(process.env.NODE_ENV || 'development')
    }),
    new MiniCssExtractPlugin({
      filename: 'style.css'
    }),
    new CopyPlugin([
      {
        from: path.join(__dirname, '/ui/static/*'),
        to: path.join(__dirname, '/res/ui/'),
        context: path.join(__dirname, '/ui/static')
      }
    ])
  ],
  resolve: {
    extensions: ['.ts', '.tsx', '.js', '.scss', '.css']
  },
  performance: { hints: false },
  devtool: process.env.NODE_ENV === 'production' ? false : 'source-map',
  mode: process.env.NODE_ENV === 'production' ? 'production' : 'development'
};
