import React, { Component } from 'react'

import styled from 'styled-components';

import { Redirect } from 'react-router-dom';

import { vw } from '../../utils';

const Popup = styled.div`
    position: fixed;
    left: 0;
    right: 0;
    top: 0;
    bottom: 0;
    background-color: rgba(0,0,0,0.6);
`;

const Center = styled.div`
    position: absolute;
    left: 50%;
    top: 50%;
    transform:translateX(-50%) translateY(-50%);
    width: ${vw(600)};

`

const Wrapper = styled.div`
    height: ${vw(700)};
    border-radius: ${vw(20)};
    background-color: #fff;
`

const Title = styled.div`
      box-sizing: border-box;
      text-align: center;
      height: ${vw(110)};
      font-size: ${vw(38)};
      padding-top: ${vw(34)};
`

const Body = styled.div`
      box-sizing: border-box;
      height: ${vw(590)};
      padding-bottom: ${vw(30)};
      font-size: ${vw(28)};
`

const Content = styled.div`
      box-sizing: border-box;
      height: ${vw(560)};
      padding: ${vw(10)} ${vw(30)};
      color:#555;
      overflow-y: auto;
      line-height: 1.6;
`

const Close = styled.div`
    width: ${vw(80)};
    height: ${vw(80)};

    background-image: url(${require('./btn_close.png')});
    background-size: cover;
    
    margin: 0 auto;
    margin-top: ${vw(20)};
`

// .modal {
//     position: absolute;
//     left: 50%;
//     top: 50%;
//     transform:translateX(-50%) translateY(-50%);
//     width: ${vw(600)};
  
  
//     .modal__wrapper {
//       height: ${vw(700)};
//       border-radius: ${vw(20)};
//       background-color: #fff;
//     }
  
//     .modal__title {
//       box-sizing: border-box;
//       text-align: center;
//       height: ${vw(110)};
//       font-size: ${vw(38)};
//       padding-top: ${vw(34)};
//     }
  
//     .modal__body {
//       box-sizing: border-box;
//       height: ${vw(590)};
//       padding-bottom: ${vw(30)};
//       font-size: ${vw(28)};
//     }
  
//     .modal__content {
//       box-sizing: border-box;
//       height: ${vw(560)};
//       padding: ${vw(10)} ${vw(30)};
//       color:#555;
//       overflow-y: auto;
//       line-height: 1.6;
//     }
  
//     .modal__close {
//       .background("btn_close.png");
//       .center-block;
//       margin-top: ${vw(20)};
//     }
//}

class Modal extends Component {
    render() {
        const {title = null, body = null, ...rest} = this.props;

        return <Popup {...rest}>
            <Center>
                <Wrapper>
                    <Title>{title}</Title>
                    <Body>
                        <Content>
                            {body}
                        </Content>
                    </Body>
                </Wrapper>
                <Close onClick={this.props.onClose}/>
            </Center>
        </Popup>
    }
}

export default class Rule extends Component {
    state = {
        closed: false
    }

    render() {
        const content = `
        赛事奖励（排行榜胜点高到低排序）：
        第一名到第十名: 奖励价值50元优惠券
        
        -
        
        比赛时间：
        
        2018年4月29日 - 2018年5月1日（23点59分 排行榜排序为准）
        
        -
        
        比赛说明：
        
        1、跳一跳PK赛，可随机匹配玩家也可以通过微信和QQ分享好友在线实时一同进行游戏；
        2、可以在游戏死亡后，使用【复活卡】进行复活；每人起始有三张复活卡，通过微信和QQ分享好友一同游戏，每局游 戏双方都可以获得一张复活卡奖励；
        3、如发现任何违规、套取奖励行为将视情节严重程度进行判罚：不予发放奖励、冻结通过推荐有奖所获得的奖励、依 法追究其法律责任；
        4、活动奖励将在5月2日中午12点统一进行发放；发放奖励以优惠券组合礼包发放到排行榜相应安锋账号中，可以在游 戏充值进行使用；
        `;

        const body = (
            <div>
                {
                    content.split('\n').map(line => <p>{line.trim()}</p>)
                }
            </div>
        );

        const { gameId } = this.props.match.params;
        if (this.state.closed) {
            return <Redirect to={`/game/${gameId}`}/>
        } else {
            return <Modal title={'游戏规则'} body={body} onClose={() => {
                this.setState({
                    closed: true
                })
            }}/>
        }
    }
}