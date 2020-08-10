// 上传佐证组件代码
import Taro, { Component } from "@tarojs/taro";
import { connect } from "@tarojs/redux";

import { View, Input, Image, Text } from "@tarojs/components";
import { AtIcon, AtProgress } from "taro-ui";
import _ from "lodash";

import PropTypes from "prop-types";
import { API } from "@/constants/api";

import "./uploadfile.scss";

@connect(({}) => ({}))
export default class GtUploadFile extends Component {
  static propTypes = {
    existingFiles: PropTypes.array, //已经存在的文件
    uploadFileType: PropTypes.string, //上传文件类型控制,通过uploadFileType控制多个类型以逗号隔开，示例:uploadFileType:"image/gif,image/jp2"
    maxCount: PropTypes.number, //最大上传图片数
    multiple: PropTypes.bool, //是否支持多选
    onUploadFile: PropTypes.func, // 点击上传按钮时触发的事件
    onConfirom: PropTypes.func, // 点击删除按钮时触发的事件
    onRemoveImage: PropTypes.func, // 点击删除按钮时触发的事件
    tempFile: PropTypes.bool, //是文件为true 图片为false
    recTypeId: PropTypes.number,
    recId: PropTypes.number,
    subRecTypeId: PropTypes.number,
    subRecId: PropTypes.number
  };
  static defaultProps = {
    existingFiles: [],
    uploadFileType: "image/*",
    maxCount: 1, //最大上传图片数
    multiple: true, //是否支持多选
    onUploadFile: null, //获取上传的
    onConfirom: null, //获取删除的
    onRemoveImage: null, //删除上传的
    tempFile: false,
    recTypeId: 0,
    recId: 0,
    subRecTypeId: 0,
    subRecId: 0
  };
  constructor(props) {
    super(props);
    this.state = {
      uploadTypeIsImage: true, //上传文件类型是否是上传图片
      viewFilesData: [], //图片文件，在view上面显示
      showFileProgress: false, //是否显示文件上传进度条
      fileProgress: 0, //文件上传进度条进度
      hiddenUploadBtn: false, //隐藏上传按钮
      isRemovFile: false
    };
  }
  componentWillMount() {
    const { uploadFileType } = this.props;
    let uploadType = _.startsWith("image/*", uploadFileType, 0);
    if (uploadType) {
      this.setState({
        uploadTypeIsImage: true
      });
    } else {
      this.setState({
        uploadTypeIsImage: false
      });
    }
  }
  componentWillReceiveProps(nextProps) {
    //如果值是异步获取
    const { existingFiles } = nextProps;
    const { viewFilesData, isRemovFile } = this.state;
    if (!viewFilesData.length && !isRemovFile) {
      this.setState({
        viewFilesData: existingFiles,
        hiddenUploadBtn: true //隐藏上传按钮
      });
    } else {
      return;
    }
  }

  //确定
  handConfirm() {
    this.props.onConfirom();
  }

  //获取服务器和本地用户的信息
  getIncidentalInfo() {
    const rootUrl = API.UPLOADURL; // 服务器地址
    const token = Taro.getStorageSync("token"); // 图片上传需要单独写入token
    const locale = Taro.getStorageSync("locale"); // 图片上传需要单独写入token
    const userLoginName = Taro.getStorageSync("userLoginName"); // 图片上传需要单独写入token
    return {
      rootUrl,
      token,
      locale,
      userLoginName
    };
  }

  //触发选择文件的事件
  handChooesFile() {
    if (process.env.TARO_ENV === "weapp") {
      //微信环境下触发微信小程序原生API的事件
      const { uploadTypeIsImage } = this.state;
      if (uploadTypeIsImage) {
        this.weappHandUpFiles();
      } else {
        this.uploadMessageFile();
      }
    } else if (process.env.TARO_ENV === "h5") {
      //h5环境下直接触发input框的点击事件
      let uploadInput = document.getElementById("uploadInput");
      uploadInput.click();
      //点击以后会触发input框上面的onChange事件
    }
  }
  /** -----h5端的处理----- */
  //input框的onChange事件
  handSelectFiles(e) {
    let files = e.target.files;
    this.handUploadFilesFun(files);
  }
  //文件循环遍历处理
  handUploadFilesFun(files) {
    let that = this;
    const { uploadTypeIsImage } = this.state;
    let tipsText = uploadTypeIsImage ? "张图片" : "个文件";
    let fileNamesData = [];
    let imagesSrcData = []; //转化为blob格式在浏览器上显示缓存的图片
    let uploadfilemaxsize = 10 * 1024 * 1024; //大小的上限
    let uploadData = []; //确定按钮时获取的值
    if (!files.length) {
      return;
    } else {
      for (let index = 0; index < files.length; index++) {
        let filesItem = files[index];
        if (filesItem.size > uploadfilemaxsize) {
          let uploadfilemsg =
            "上传文件大小超过系统规定上限(10M),请重新选择图片";
          Taro.showToast({
            title: uploadfilemsg,
            icon: "success",
            duration: 2000
          });
          return;
        } else {
          this.uploadFileItem(filesItem, index, data => {
            if (!data) {
              return;
            } else {
              if (data.success) {
                if (uploadTypeIsImage) {
                  //上传图片就添加URL
                  let filesSrcItem = URL.createObjectURL(filesItem);
                  imagesSrcData.push(filesSrcItem);
                } else {
                  //上传文件就添加filename
                  let fileName = filesItem.name;
                  fileNamesData.push(fileName);
                }
                let showViewFilesData = uploadTypeIsImage
                  ? imagesSrcData
                  : fileNamesData;
                that.setState({
                  viewFilesData: showViewFilesData
                });
                Taro.showToast({
                  title: `第${index + 1}${tipsText}上传成功`,
                  icon: "success",
                  duration: 2000
                });
                uploadData.push(data.data[0]);
                if (index === files.length - 1) {
                  that.props.onUploadFile(uploadData);
                  return;
                }
              } else {
                Taro.showToast({
                  title: `第${index + 1}${tipsText}上传失败`,
                  icon: "error",
                  duration: 2000
                });
                that.setState({
                  hiddenUploadBtn: false
                });
                console.log(data, "error");
                return;
              }
            }
          });
        }
      }
    }
  }
  //单个文件(图片)上传服务器
  uploadFileItem(filesItem, index, callback) {
    let that = this;
    const { rootUrl, token, locale, userLoginName } = this.getIncidentalInfo();
    const {
      maxCount,
      tempFile,
      recTypeId,
      recId,
      subRecTypeId,
      subRecId
    } = this.props;
    if (index > maxCount) {
      return;
    } else {
      let formData = new FormData();
      let xhr = new XMLHttpRequest();
      xhr.open("POST", rootUrl, true); // 第三个参数为async?，异步/同步
      formData.append(filesItem.name, filesItem);
      //把请求相关参数放入formData中
      formData.append("tempFile", tempFile);
      formData.append("recTypeId", recTypeId);
      formData.append("recId", recId);
      formData.append("subRecTypeId", subRecTypeId);
      formData.append("subRecId", subRecId);

      xhr.setRequestHeader("Authorization", `Bearer ${token}`);
      xhr.setRequestHeader("userToken", token);
      xhr.setRequestHeader("userLoginName", userLoginName);
      xhr.setRequestHeader("userLanguage", locale);

      that.setState({
        showFileProgress: true
      });
      //监听请求的进度并在回调中传入进度参数
      xhr.upload.addEventListener(
        "progress",
        e => {
          if (e.lengthComputable) {
            let progress = Math.round((e.loaded / e.total) * 100);
            that.setState({
              fileProgress: progress
            });
            if (progress === 100) {
              setTimeout(function() {
                that.setState({
                  hiddenUploadBtn: true,
                  showFileProgress: false
                });
              }, 30);
            }
          }
        },
        false
      ); // 第三个参数为useCapture?，是否使用事件捕获/冒泡
      //监听readyState的变化，完成时回调后端返回的response
      xhr.addEventListener(
        "readystatechange",
        e => {
          let response = e.currentTarget.response
            ? JSON.parse(e.currentTarget.response)
            : null;
          if (e.currentTarget.readyState === 4) {
            callback(response);
            xhr.upload.removeEventListener(
              "progress",
              event => {
                if (event.lengthComputable) {
                  let progress = Math.round((event.loaded / event.total) * 100);
                  that.setState({
                    fileProgress: progress
                  });
                  if (progress === 100) {
                    that.setState({
                      showFileProgress: false
                    });
                  }
                }
              },
              false
            );
          } else {
            console.log("upload loading ... ");
          }
        },
        false
      );
      xhr.send(formData);
    }
  }

  /** -----weapp端的处理---- */
  //图片
  weappHandUpFiles() {
    let that = this;
    const {
      maxCount,
      tempFile,
      recTypeId,
      recId,
      subRecTypeId,
      subRecId
    } = this.props;
    const { rootUrl, token, locale, userLoginName } = this.getIncidentalInfo();
    let imagesSrcData = []; //在浏览器上显示缓存的图片
    let imagesUploadData = []; //确定按钮时获取的值
    Taro.chooseImage({
      count: maxCount,
      sizeType: ["original", "compressed"], // 可以指定是原图还是压缩图，默认二者都有
      sourceType: ["album", "camera"], // 可以指定来源是相册还是相机，默认二者都有
      success: res => {
        // 返回选定照片的本地文件路径列表，tempFilePath可以作为img标签的src属性显示图片
        let tempFilePaths = res.tempFilePaths;
        for (let i = 0; i < tempFilePaths.length; i++) {
          if (i > maxCount) {
            return;
          } else {
            let tempFileItem = tempFilePaths[i];
            imagesSrcData.push(tempFileItem);
            const uploadTask = Taro.uploadFile({
              url: rootUrl, //里面填写你的上传图片服务器API接口的路径
              filePath: tempFileItem, //要上传文件资源的路径 String类型
              name: "file", //按个人情况修改，文件对应的 key,开发者在服务器端通过这个 key 可以获取到文件二进制内容，(后台接口规定的关于图片的请求参数)
              header: {
                "Content-Type": "multipart/form-data", //记得设置
                // userToken: token,
                // userLoginName: userLoginName,
                // userLanguage: locale,
                Authorization: `Bearer ${token}`
              },
              formData: {
                //和服务器约定的token, 一般也可以放在header中
                tempFile: tempFile,
                recTypeId: recTypeId,
                recId: recId,
                subRecTypeId: subRecTypeId,
                subRecId: subRecId
              },
              success: data => {
                let result = JSON.parse(data.data);
                if (result.success) {
                  Taro.showToast({
                    title: `第${i + 1}张图片上传成功`,
                    icon: "success",
                    duration: 2000
                  });
                  imagesUploadData.push(result.data);
                  if (i === tempFilePaths.length - 1) {
                    this.props.onUploadFile(imagesUploadData);
                  }
                } else {
                  Taro.showToast({
                    title: `第${i + 1}张图片上传失败`,
                    icon: "error",
                    duration: 2000
                  });
                  return;
                }
              },
              fail: err => {
                console.log(err);
              }
            });
            uploadTask.progress(progress => {
              this.setState({
                fileProgress: progress.progress
              });
              if (progress.progress === 100) {
                that.setState({
                  hiddenUploadBtn: true
                });
                setTimeout(function() {
                  that.setState({
                    showFileProgress: false
                  });
                }, 600);
              }
            });
          }
        }
        this.setState({
          viewFilesData: imagesSrcData
        });
      },
      fail: err => {
        console.log(err);
      }
    });
  }
  //文件
  uploadMessageFile() {
    let that = this;
    const { tempFile, recTypeId, recId, subRecTypeId, subRecId } = this.props;
    const { rootUrl, token, locale, userLoginName } = this.getIncidentalInfo();
    let fileNamesData = [];
    let fileUploadData = [];
    Taro.chooseMessageFile({
      count: 1,
      type: "file",
      success: res => {
        // 返回选定的本地文件路径列表
        let tempFiles = res.tempFiles;
        for (let i = 0; i < tempFiles.length; i++) {
          if (!tempFiles.length) {
            return;
          } else {
            let tempFileItem = tempFiles[i];
            let fileName = tempFileItem.name;
            // let fileItemName = fileName.substring(0,fileName.lastIndexOf("."))
            fileNamesData.push(fileName);
            this.setState({
              showFileProgress: true
            });
            const uploadTask = Taro.uploadFile({
              url: rootUrl, //里面填写你的上传图片服务器API接口的路径
              filePath: tempFileItem.path, //要上传文件资源的路径 String类型
              name: "gantIbom", //按个人情况修改，文件对应的 key,开发者在服务器端通过这个 key 可以获取到文件二进制内容，(后台接口规定的关于图片的请求参数)
              header: {
                "Content-Type": "multipart/form-data", //记得设置
                userToken: token,
                userLoginName: userLoginName,
                userLanguage: locale
                // Authorization: `Bearer ${token}`
              },
              formData: {
                //和服务器约定的token, 一般也可以放在header中
                tempFile: tempFile,
                recTypeId: recTypeId,
                recId: recId,
                subRecTypeId: subRecTypeId,
                subRecId: subRecId
              },
              success: data => {
                let result = JSON.parse(data.data);
                if (result.success) {
                  Taro.showToast({
                    title: `第${i + 1}个文件上传成功`,
                    icon: "success",
                    duration: 2000
                  });
                  fileUploadData.push(result.data);
                  if (i === tempFiles.length) {
                    this.props.onUploadFile(fileUploadData);
                  }
                } else {
                  Taro.showToast({
                    title: `第${i + 1}个文件上传失败`,
                    icon: "error",
                    duration: 2000
                  });
                  return;
                }
              },
              fail: err => {
                console.log(err);
              }
            });
            uploadTask.progress(progress => {
              this.setState({
                fileProgress: progress.progress
              });
              if (progress.progress === 100) {
                setTimeout(function() {
                  that.setState({
                    showFileProgress: false,
                    hiddenUploadBtn: true
                  });
                }, 30);
              }
            });
          }
        }
        this.setState({
          viewFilesData: fileNamesData
        });
      },
      fail: err => {
        console.log(`文件上传失败，查看错误信息:${err}`);
      }
    });
  }

  //文件上传缓存以后确定选中的文件上传服务器
  handConfirmFile() {
    this.props.onConfirom();
  }
  //删除文件
  handDeleteFile() {
    this.setState(
      {
        hiddenUploadBtn: false,
        isRemovFile: true,
        viewFilesData: []
      },
      () => {
        const { viewFilesData } = this.state;
        this.props.onRemoveImage(viewFilesData);
      }
    );
  }
  render() {
    const {
      uploadTypeIsImage,
      fileProgress,
      viewFilesData,
      showFileProgress,
      hiddenUploadBtn
    } = this.state;
    console.log(hiddenUploadBtn, "hiddenUploadBtn");
    const { uploadFileType } = this.props;
    let showUploadFileView = null;
    let uploadIcon = uploadTypeIsImage ? "image" : "upload";
    let uplpadIconSize = uploadTypeIsImage ? "40" : "18";
    let fileProgressClass = showFileProgress
      ? "progress show-progress"
      : "progress hidden-progress";
    showUploadFileView = (
      <View className="uploadfile-content">
        {viewFilesData.map(fileVal => {
          return uploadTypeIsImage ? (
            <View className="img-item-warp">
              <View className="img-warp">
                <Image className="img-item" src={fileVal} />
              </View>
              <View
                className="conforim-img"
                onClick={this.handConfirmFile.bind(this)}
              >
                确认
              </View>
              <View className="icon-btn-warp">
                <AtIcon
                  className="icon-btn"
                  onClick={this.handDeleteFile.bind(this)}
                  value="close-circle"
                  size="24"
                  color="#6e6e6e"
                ></AtIcon>
              </View>
            </View>
          ) : (
            <View className="file-item-warp">
              <AtIcon
                className="icon-btn"
                value="file-generic"
                size="20"
                color="#336633"
              ></AtIcon>
              <Text className="uploadfile-name">{fileVal}</Text>
              <AtIcon
                className="close-icon-btn"
                onClick={this.handDeleteFile.bind(this)}
                value="close-circle"
                size="12"
                color="#6e6e6e"
              ></AtIcon>
            </View>
          );
        })}
      </View>
    );
    return (
      <View className="uploadfile-components">
        {showUploadFileView}
        <View
          className={
            hiddenUploadBtn
              ? "btn-warp btn-warp-hidden"
              : "btn-warp btn-warp-show"
          }
        >
          <View
            className={uploadTypeIsImage ? "upload-img-btn" : "upload-file-btn"}
            onClick={this.handChooesFile.bind(this)}
          >
            <AtIcon
              className="upload-icon"
              value={uploadIcon}
              size={uplpadIconSize}
              color="#fff"
            ></AtIcon>
          </View>
        </View>
        <View className={fileProgressClass}>
          <AtProgress percent={fileProgress} />
        </View>
        <Input
          className="uploadfile-ipt"
          type="file"
          name="file"
          id="uploadInput"
          accept={uploadFileType}
          onChange={this.handSelectFiles.bind(this)}
        />
      </View>
    );
  }
}
