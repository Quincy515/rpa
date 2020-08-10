import Taro from "@tarojs/taro";
import { View } from "@tarojs/components";
import { AtImagePicker, AtButton } from "taro-ui";

import { COUNT } from "../../config";
import { API } from "@/constants/api";

export default class ChooseImage extends Taro.Component {
  constructor() {
    super(...arguments);
    this.state = {
      files: [],
      showUploadBtn: false,
      upLoadImg: []
    };
  }

  componentWillMount() {
    console.log("this.props.chooseImg", this.props.chooseImg);
    this.setState({
      files: this.props.chooseImg.files,
      showUploadBtn: this.props.chooseImgshowUploadBtn,
      upLoadImg: this.props.chooseImg.upLoadImg
    });
  }

  componentDidShow() {}

  componentDidHide() {}

  // files 值发生变化触发的回调函数, operationType 操作类型有添加，移除，如果是移除操作，则第三个参数代表的是移除图片的索引
  onChange(values, operationType, index) {
    if (operationType === "remove") {
      this.setState(
        prevState => {
          let oldSendImg = prevState.upLoadImg;
          oldSendImg.splice(oldSendImg[index], 1); // 删除已上传的图片地址
          return {
            files: values,
            upLoadImg: oldSendImg
          };
        },
        () => {
          const { files } = this.state;
          // 设置删除数据图片地址
          if (files.length === COUNT) {
            // 最多图片个数，隐藏添加图片按钮
            this.setState({
              showUploadBtn: false
            });
          } else if (files.length === 0) {
            this.setState({
              upLoadImg: []
            });
          } else {
            this.setState({
              showUploadBtn: true
            });
          }
        }
      );
    } else {
      values.map((item, index) => {
        if (item.url.indexOf(".pdf") > -1 || item.url.indexOf("PDF") > -1) {
          values[index].url = require("@/asset/images/PDF.png");
        }
      });
      this.setState(
        () => {
          return { files: values };
        },
        () => {
          const { files } = this.state;
          if (files.length == COUNT) {
            // 最多图片个数，隐藏添加图片按钮
            this.setState({
              showUploadBtn: false
            });
          } else {
            this.setState({
              showUploadBtn: true
            });
          }
        }
      );
    }
  }

  // 选择失败回调
  onFail(mes) {
    console.log(mes);
  }

  // 点击图片回调
  onImageClick(index, file) {
    console.log(index, file);
    let imgs = [];
    this.state.files.map((item, index) => {
      imgs.push(item.file.path);
    });
    if (imgs[index].indexOf(".pdf") > -1 || imgs[index].indexOf(".PDF") > -1) {
      Taro.downloadFile({
        url: imgs[index],
        success: function(res) {
          let filePath = res.tempFilePath;
          Taro.openDocument({
            filePath: filePath,
            success: function(res) {
              console.log("打开文档成功");
            }
          });
        }
      });
    } else {
      Taro.previewImage({
        current: imgs[index], // 当前显示图片
        urls: imgs // 所有图片
      });
    }
  }

  toUpload = () => {
    const { files } = this.state;
    const token = Taro.getStorageSync("token"); // 图片上传需要用户token
    if (files.length > 0 && token) {
      this.props.onFilesValue(files);
      const url = API.UPLOADURL;
      this.uploadRequest({ url, token, path: files });
    } else if (!token) {
      Taro.showToast({
        title: "请先登录",
        icon: "none",
        duration: 2000
      });
    } else {
      Taro.showToast({
        title: "请先选择图片",
        icon: "none",
        duration: 2000
      });
    }
  };

  // 上传组件
  uploadRequest = data => {
    let that = this;
    let i = data.i ? data.i : 0; // 当前上传的图片位置
    let success = data.success ? data.success : 0; // 上传成功的个数
    let fail = data.fail ? data.fail : 0; // 上传失败的个数
    Taro.showLoading({
      title: `正在上传第${i + 1}张`
    });
    // 发起上传
    Taro.uploadFile({
      url: data.url,
      header: {
        "content-type": "multipart/form-data",
        Authorization: "bearer " + data.token // 上传需要单独处理token
      },
      name: "file",
      filePath: data.path[i].url,
      success: resp => {
        // 图片上传成功，图片上传成功的变量+1
        let resultData = JSON.parse(resp.data);
        if (resultData.code === "1000") {
          success++;
          this.setState(
            prevState => {
              let oldUpload = prevState.uploadImg;
              oldUpload.push(resultData.data);
              return {
                uploadImg: oldUpload
              };
            },
            () => {
              //setState 会合并所有的setState操作，所以在这里等待图片传完之后再调用设置url方法
              // this.setFatherUploadSrc() // 设置数据图片地址字段
            }
          );
        } else {
          fail++;
        }
      },
      fail: () => {
        fail++; // 图片上传失败，图片上传失败的数量+1
      },
      complete: () => {
        Taro.hideLoading();
        i++; // 这个图片执行完上传后，开始上传下一个
        if (i == data.path.length) {
          // 当图片传完后，停止调用
          Taro.showToast({
            title: "上传完成",
            icon: "success",
            duration: 2000
          });
          console.log("成功: ", success, "失败: ", fail);
        } else {
          // 如图片还没有传完，则继续调用函数
          data.i = i;
          data.success = success;
          data.fail = fail;
          that.uploadRequest(data);
        }
      }
    });
  };
  render() {
    const { showUploadBtn } = this.state;
    return (
      <View>
        <AtImagePicker
          fmultiple={false}
          length={3} //单行的图片数量
          files={this.state.files}
          onChange={this.onChange.bind(this)}
          onFail={this.onFail.bind(this)}
          onImageClick={this.onImageClick.bind(this)}
          showAddBtn={showUploadBtn} //是否显示添加图片按钮
        />
        <AtButton
          type="primary"
          className="poof_submit_btn"
          onClick={this.toUpload}
        >
          上传图片
        </AtButton>
      </View>
    );
  }
}
