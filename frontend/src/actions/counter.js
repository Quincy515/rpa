import Taro, { login } from "@tarojs/taro";

import { getApi, postApi, postFormApi } from "@/utils/request";
import { ADD, MINUS, LIST } from "@/constants/counter";
import { API } from "@/constants/api";

export const add = () => {
  return {
    type: ADD
  };
};
export const minus = () => {
  return {
    type: MINUS
  };
};

// 异步的action
export function asyncAdd() {
  return dispatch => {
    setTimeout(() => {
      dispatch(add());
    }, 2000);
  };
}

const testList = data => {
  return {
    type: LIST,
    payload: data
  };
};

export function getList() {
  return dispatch => {
    //   Taro.request({
    //     url: "http://localhost:2020/api/questions/",
    //     header: {
    //       "content-type": "application/json",
    //     },
    //   }).then((res) => {
    //     dispatch(testList(res.data));
    //   });

    getApi("http://localhost:2020/api/questions/", {
      offset: 1,
      limit: 3
    }).then(res => {
      dispatch(testList(res));
    });

    // postFormApi("http://localhost:2020/api/signup", {
    //   email: "boaa@163.com",
    //   password: "Aa123456",
    //   confirm_password: "Aa123456",
    // }).then((res) => {
    //   dispatch(testList(res));
    // });

    // postApi("http://localhost:2020/api/signup", {
    //   email: "boabc@163.com",
    //   password: "Aa123456",
    //   confirm_password: "Aa123456",
    // }).then((res) => {
    //   console.log("res===>", res);
    //   dispatch(testList(res));
    // });
  };
}
