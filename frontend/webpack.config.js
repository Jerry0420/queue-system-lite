const path = require('path')
const MiniCssExtractPlugin = require('mini-css-extract-plugin')
const HtmlWebpackPlugin = require('html-webpack-plugin')
const webpack = require("webpack")

require('dotenv').config({ path: '../../etc/config/.env.frontend' })

module.exports = (env, argv) => {
  return {
      entry: ['@babel/polyfill', './src/index.tsx'],
      output: {
        filename: '[contenthash]..bundle.js',
        path: path.resolve(__dirname, './dist/'),
      },
      resolve: { extensions: ['.js', '.jsx', '.ts', '.tsx'] },
      module: {
        rules: [
          {
            test: /.ts$/,
            use: {
              loader: 'babel-loader',
              options: {
                presets: ['@babel/typescript', '@babel/preset-env'],
              },
            },
          },
          {
            test: /.tsx$/,
            use: {
              loader: 'babel-loader',
              options: {
                presets: ['@babel/typescript', '@babel/preset-react', '@babel/preset-env'],
              },
            },
          },
          {
            test: /\.(scss)$/,
            use: [
              MiniCssExtractPlugin.loader,
              'css-loader',
              'sass-loader',
              'postcss-loader'
            ],
          },
        ],
      },
      plugins: [
        new MiniCssExtractPlugin({
          filename: './bundle.css',
        }),
        new HtmlWebpackPlugin({
          template: './public/index.html',
        }),
        new webpack.DefinePlugin({
          "process.env": JSON.stringify(process.env),
        }),
      ],
      devServer: {
        static: {
          directory: path.join(__dirname, 'dist'),
        },
        compress: true,
        port: 3000,
        allowedHosts: 'all'
      },
      mode: 'development',
    }
}