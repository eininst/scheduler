import {defineConfig} from 'umi';
import aliyunTheme from '@ant-design/aliyun-theme';

console.log(process.env.NODE_ENV)

var configs: any = {
  // dynamicImport: {}
}
// development
// production
if (process.env.NODE_ENV == "production") {
  configs = {
    // dynamicImport: {},
    // mfsu: {},
    publicPath:"/assets/"
  }
}
export default defineConfig({
  // title: '管理平台业务模版',
  // base: "/mp",
  ...configs,
  nodeModulesTransform: {
    type: 'none',
  },
  layout: {},
  theme: aliyunTheme,
  fastRefresh: {},
  // theme: {
  //   '@primary-color': '#1DA57A',
  // },
  antd: {
    dark: false,
    compact: false,
  },
  cssLoader: {
    localsConvention: 'camelCase',
  },
});
