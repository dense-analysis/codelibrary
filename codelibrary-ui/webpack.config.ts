import * as os from 'os'
import * as path from 'path'
import {Configuration} from 'webpack'
import * as webpackDevServer from 'webpack-dev-server'

const config: Configuration = {
  entry: path.resolve(__dirname, 'src', 'index.tsx'),
  output: {
    path: path.resolve(__dirname, 'dist'),
    filename: 'codelibrary-ui.js',
    publicPath: '/static/js',
  },
  module: {
    rules: [
      {
        test: /\.[jt]sx?$/,
        use: {
          loader: 'babel-loader',
          options: {
            presets: [
              [
                '@babel/preset-env',
                {
                  targets: {chrome: 100},
                },
              ],
              ['@babel/preset-react', {runtime: 'automatic'}],
              '@babel/preset-typescript',
            ],
          },
        },
        include: path.resolve(__dirname, 'src'),
        exclude: /node_modules/,
      },
      {
        test: /\.css$/,
        use: ['style-loader', 'css-loader'],
      },
    ],
  },
  mode: 'production',
  performance: {
    hints: false,
    maxEntrypointSize: 5 * 1024 * 1024,
    maxAssetSize: 5 * 1024 * 1024,
  },
  devServer: {
    hot: true,
    client: {
      logging: 'none',
      overlay: false,
    },
    static: {
      directory: path.join(__dirname, 'public/'),
    },
    allowedHosts: 'all',
    historyApiFallback: true,
    port: process.env.DEV_SERVER_PORT || '8000',
  },
}

export default config
