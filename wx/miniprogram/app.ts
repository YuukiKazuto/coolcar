import { Coolcar } from "./service/request"

App<IAppOption>({
  globalData: {

  },
  onLaunch() {
    // 登录
    Coolcar.login()
  },
})