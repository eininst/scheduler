import {BarChartOutlined, FileTextOutlined, FundOutlined, TeamOutlined,} from "@ant-design/icons";

var r = [ //此功能可以实现动态路由，用来渲染访问路由
  {
    path: '/',
    name: 'Dashboard',
    icon: <BarChartOutlined/>,
  },

  {
    path: '/task',
    name: '任务管理',
    icon: <FundOutlined/>,
  },
  {
    path: '/log',
    name: '请求日志',
    icon: <FileTextOutlined/>,
  },
]

var role = eval("window.role")
if (role == "admin") {
  r.push({
    path: '/user',
    name: '用户管理',
    icon: <TeamOutlined/>,
  },)
}

export const routers = r
