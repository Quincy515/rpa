// 五种级别：debug info warn error fatal
import Taro from "@tarojs/taro";
import fundebug from "fundebug-wxjs";
import { formatTime } from ".";
import { stringify } from "./serialize";

/**
 * 打印一些感兴趣的重要信息，或者bi上报的日志
 * @param {string} name
 * @param {Object} option
 */
const info = (name, option) => {
  report(name, option, "info");
};

/**
 * 有一些警告信息，但是也要上报
 * @param {string} name
 * @param {Object} option
 */
const warn = (name, option) => {
  report(name, option, "warn");
};

/**
 * 用于调试的一些详细信息
 * @param {string} name
 * @param {Object} option
 */
const debug = (name, option) => {
  report(name, option, "debug");
};

/**
 * 普通的异常
 * @param {string} name
 * @param {Object} option
 */
const error = (name, option) => {
  report(name, option, "error");
};

/**
 * 致命的错误 项目无法进行的那种
 * @param {string} name
 * @param {Object} option
 */
const fatal = (name, option) => {
  report(name, option, "fatal");
};

const report = (name, option, type = "info") => {
  // 设备信息
  try {
    var deviceInfo = Taro.getSystemInfoSync();
    var device = JSON.stringify(deviceInfo);
  } catch (e) {
    console.log("not support sysinfo", e.message);
  }

  option = stringify(option); //[object,object]

  // 记录时间和用户信息
  let time = formatTime(new Date());
  const user = Taro.getStorageInfoSync("userInfo");
  // TODO: 判断是否有用户信息
  console.log("user===>", user);
  if (type == "info" || type == "debug") {
    console.log(time, type, user, name, option, device);
  } else {
    fundebug.notify(type, option, {
      metaData: {
        device: device,
        user: user,
        option: option,
        name: name,
        type: type,
        time: time,
      },
    });
    // console.error(time, type, user, name, option, device);
  }
};

export default { debug, info, warn, error, fatal };
