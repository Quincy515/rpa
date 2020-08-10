import Taro from "@tarojs/taro";
import { AtCard } from "taro-ui";
export default class ActivateCard extends Taro.Component {
  constructor() {
    super(...arguments);
    this.state = {};
  }

  onCardClick(index, file) {
    console.log(index, file);
  }
  render() {
    return (
      <AtCard
        note="小Tips"
        extra="额外信息"
        title="这是个标题"
        thumb="http://www.logoquan.com/upload/list/20180421/logoquan15259400209.PNG"
        onClick={this.onCardClick}
      >
        这也是内容区 可以随意定义功能
      </AtCard>
    );
  }
}
