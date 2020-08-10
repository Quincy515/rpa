import Taro from "@tarojs/taro";
import { View } from "@tarojs/components";
import { AtAccordion, AtList, AtListItem, AtFloatLayout } from "taro-ui";
import ActivateCard from "@/components/card";
import ChooseImage from "@/components/image";

export default class PlanIndex extends Taro.Component {
  constructor() {
    super(...arguments);
    this.state = {
      open: false,
      chooseImg: {
        files: [],
        showUploadBtn: true,
        upLoadImg: []
      },
      files: []
    };
  }
  handleClick(value) {
    this.setState({
      open: value
    });
  }
  handleClose(value) {}
  render() {
    return (
      <AtList hasBorder={false}>
        <AtAccordion
          open={this.state.open}
          onClick={this.handleClick.bind(this)}
          title="标题一"
        >
          <AtFloatLayout
            isOpened
            title="这是个标题"
            onClose={this.handleClose.bind(this)}
          >
            <ActivateCard />
          </AtFloatLayout>
        </AtAccordion>
        <AtAccordion
          title="标题三"
          icon={{ value: "chevron-down", color: "red", size: "15" }}
        />
      </AtList>
    );
  }
}
