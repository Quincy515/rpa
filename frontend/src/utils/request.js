import Taro from "@tarojs/taro";
import queryString from "query-string";
import { HTTP_STATUS, COMMON_STATUS } from "@/constants/code";
import { API } from "@/constants/api";
import log from "@/utils/log";

const METHOD_GET = "GET";
const METHOD_POST = "POST";
const TRY_AGAIN_COUNT = 3;

const requestApi = (cfg, tryAgainCount = TRY_AGAIN_COUNT) => {
  // 添加 token
  const token = Taro.getStorageSync("token");
  let header = Object.assign({}, cfg.header, {
    Authorization: "Bearer " + token
  });
  let config = {
    ...cfg,
    header
  };

  return Taro.request(config)
    .then(res => {
      if (
        res.statusCode == HTTP_STATUS.SUCCESS &&
        res.data.code == COMMON_STATUS.SUCCESS
      ) {
        log.debug("request_api_success", res.data);
        return res.data;
      } else {
        if (tryAgainCount) {
          return requestApi(cfg, tryAgainCount - 1); // 失败后重试3次
        }
        log.error("request_api_err", res.data);
        // return res.data;
      }
    })
    .catch(e => {
      if (tryAgainCount) {
        return requestApi(cfg, tryAgainCount - 1); // 失败后重试3次
      }
      log.error("request_api_err", e);
    });
};

// 封装 request get 请求
export const getApi = (url, params) => {
  let thisParams = { ...params };
  if (!thisParams._timestamp) {
    thisParams._timestamp = Date.now();
  }

  let requestUrl = genUrl(url, thisParams);

  let config = {
    method: METHOD_GET,
    url: requestUrl
  };

  return requestApi(config, TRY_AGAIN_COUNT);
};

// 给 get 请求添加时间戳防止浏览器缓存
function genUrl(url, params) {
  let paramStr = queryString.stringify(params);
  let splitChar = url.indexOf("?") === -1 ? "?" : "&";
  return url + splitChar + paramStr;
}
// http://localhost:2020/api/questions/?_timestamp=1593669571542&limit=1&offset=1

// 封装 request post 请求
export const postApi = (url, params) => {
  let config = {
    url: url,
    method: METHOD_POST,
    header: {
      "Content-Type": "application/json"
    },
    data: JSON.stringify(params)
  };
  return requestApi(config, TRY_AGAIN_COUNT);
};

// 封装 request post 请求
export const postFormApi = (url, params) => {
  let config = {
    url: url,
    method: METHOD_POST,
    header: {
      "Content-Type": "application/x-wwww-from-urlencoded"
    },
    data: encodeParam(params)
  };
  return requestApi(config, TRY_AGAIN_COUNT);
};

const encodeParam = params => {
  let paramsArr = [];
  Object.keys(params).forEach(key => {
    if (typeof params[key] != "undefined") {
      paramsArr.push(key + "=" + encodeURIComponent(params[key]));
    }
  });
  return paramsArr.join("&");
};

// 微信小程序登录
export const miniLogin = () => {
  return new Promise((resolve, reject) => {
    Taro.login({
      success: function(res) {
        if (res.code) {
          postApi(API.LOGIN, { js_code: res.code })
            .then(res => {
              Taro.setStorageSync("token", res.data.access_token);
              // resolve(token); // 通过
            })
            .catch(err => {
              // reject(err); // 拒绝
            });
        }
      }
    });
  });
};
