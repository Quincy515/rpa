import { postApi } from "@/utils/request";
import { LOGIN } from "@/constants/user";
import { API } from "@/constants/api";

// 用户授权获取用户信息
const userData = data => {
  return { type: LOGIN, payload: data };
};

export function updateUser(params) {
  console.log("updateUser", params);
  return dispatch => {
    postApi(API.UPDATEUSER, params).then(res => {
      dispatch(userData(res));
    });
  };
}
