import {history} from 'umi';
import {PageLoading} from '@ant-design/pro-layout/es';
import {routers} from './routers';
import {POST} from "@/global";


export const initialStateConfig = {
  loading: <PageLoading/>,
};

// export const getInitialState = async ()=> {}
var config = eval("window.config")
var userName = eval("window.userName")
export async function getInitialState() {
  return await new Promise((resolve, reject) => {
    resolve({
      loading: false,
      pure: false,
      name: userName,
      avatar: config.avatar,
    })
  })
}

export const layout = (state: any) => {
  const {loading, pure} = state.initialState;

  return {
    logo: config.logo == "" ? false : config.logo,
    title: config.title == "" ? "Scheduler" : config.title,
    siderWidth: 250,
    loading: loading,
    pure: pure,
    logout: () => {
      POST("/api/logout", {}, (res: any, status: any) => {
        console.log(res);
        if (status == 200) {
          history.push('/login');
        }
      })
    },
    menuDataRender: () => routers,
    onPageChange: () => {
      const {location} = history;
      if (location.pathname == '/login' || location.pathname == '/init') {
        state.setInitialState({
          ...state.initialState,
          pure: true,
        })
      }
    },
  }
};

