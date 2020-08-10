import Taro from "@tarojs/taro";
import { View } from "@tarojs/components";
import ChooseImage from "@/components/image";

export default class ImagePicker extends Taro.Component {
  constructor() {
    super(...arguments);
    this.state = {
      chooseImg: {
        files: [],
        showUploadBtn: true,
        upLoadImg: []
      },
      files: []
    };
  }

  config = {
    navigationBarTitleText: "上传图片页面"
  };

  componentWillMount() {}

  // 拿到子组件上传图片的路径数组
  getOnFilesValue = value => {
    console.log(value);
    this.setState(
      {
        files: value
      },
      () => {
        console.log(this.state.files);
      }
    );
  };
  render() {
    return (
      <View className="home">
        <ChooseImage
          chooseImg={this.state.chooseImg}
          onFilesValue={this.getOnFilesValue.bind(this)}
        />
      </View>
    );
  }
}
