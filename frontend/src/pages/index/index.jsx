import Taro, { Component } from "@tarojs/taro";
import { View, Button, Text } from "@tarojs/components";
import { connect } from "@tarojs/redux";
import { AtButton } from "taro-ui";

import { add, minus, asyncAdd, getList } from "@/actions/counter";
import { updateUser } from "@/actions/user.js";
import { miniLogin } from "@/utils/request";
import { ActivateCard } from "@/components/card";
import { GtUploadFile } from "@/components/uploadfile";

import "./index.scss";

@connect(
  ({ counter, user }) => ({
    counter,
    user
  }),
  dispatch => ({
    add() {
      dispatch(add());
    },
    dec() {
      dispatch(minus());
    },
    asyncAdd() {
      dispatch(asyncAdd());
    },
    getList() {
      dispatch(getList());
    },
    updateUser(params) {
      dispatch(updateUser(params));
    }
  })
)
class Index extends Component {
  config = {
    navigationBarTitleText: "首页"
  };

  componentWillReceiveProps(nextProps) {
    console.log(this.props, nextProps);
  }

  componentWillMount() {
    const token = Taro.getStorageSync("token");
    if (!token) {
      miniLogin().then(res => {
        // 查询用户信息 存储用户信息
        console.log("查询用户信息 存储用户信息：", res);
      });
    }
  }

  componentWillUnmount() {}

  componentDidShow() {}

  componentDidHide() {}

  setUserInfo = data => {
    const userInfo = data.detail.userInfo;
    console.log("data", userInfo);
    let params = {
      nick_name: userInfo.nickName,
      avatar_src: userInfo.avatarUrl,
      is_authorized: 1
    };
    console.log("params", params);
    this.props.updateUser(params);
  };

  navigateToPlan = () => {
    // 跳转到目的页面，打开新页面
    Taro.navigateTo({
      url: "/pages/plan/index"
    });
  };

  navigateToImg = () => {
    // 跳转到目的页面，打开新页面
    Taro.navigateTo({
      url: "/pages/image/index"
    });
  };
  render() {
    const { list } = this.props.counter;
    const test = list
      ? list.data.question_list.map((item, index) => {
          return <View>{item.question_caption}</View>;
        })
      : null;
    return (
      <View className="index">
        <Button className="add_btn" onClick={this.props.add}>
          +
        </Button>
        <Button className="dec_btn" onClick={this.props.dec}>
          -
        </Button>
        <Button className="dec_btn" onClick={this.props.asyncAdd}>
          async
        </Button>
        <View>
          <Text>{this.props.counter.num}</Text>
        </View>
        <View>
          <Text>Hello, World</Text>
        </View>
        <AtButton type="primary" onClick={this.props.getList}>
          list
        </AtButton>
        <Button onGetUserInfo={this.setUserInfo} open-type="getUserInfo">
          用户授权
        </Button>
        <ImgPicker />
        <ActivateCard />
        <AtButton type="primary" onClick={this.navigateToPlan}>
          计划页面
        </AtButton>
        <AtButton type="primary" onClick={this.navigateToImg}>
          上传图片页面
        </AtButton>
        <GtUploadFile />
        {test}
      </View>
    );
  }
}

export default Index;
