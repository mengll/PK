import React, { Component, Fragment } from 'react';
import ReactDOM from 'react-dom';

import { Toast } from 'antd-mobile'
import styled from 'styled-components';

import { env, vw, request } from '../../utils';

const  { wx } = window;

const Overlay = styled.div`
    z-index: 1000;
    position: absolute;
    left: 0; right: 0; top: 0; bottom: 0;
    background: rgba(0, 0, 0, 0.6);
`;

const Tips = styled.div`
    position: absolute;
    top: 0;
    right: ${vw(100)};

    width: ${vw(511)};
    height: ${vw(496)};
    background: url(${require('./share.png')});
    background-size: cover;
`;

// class WeixinShare extends Component {
//     static jssdk = new Promise(async resolve => {
//         // const params = await request('/index.php?d=api&c=wechat&m=jssdk', {url: window.location.href });

//         //wx.ready(resolve);
//     })

//     state = {
//         visible: false,
//     }

//     container = document.createElement('div');
    
    
//     componentWillMount() {
//         const {wx} = window;
//         const {title, content, image, url, onSuccess = () => {}} = this.props;
//         WeixinShare.jssdk.then(() => {
//             wx.onMenuShareTimeline({ //分享到朋友圈
//                 title: title,
//                 desc: content, // 分享描述
//                 link: url,
//                 imgUrl: image,
//                 success: onSuccess
//             });

//             wx.onMenuShareAppMessage({ //分享到朋友
//                 title: title, // 分享标题
//                 desc: content, // 分享描述
//                 link: url, // 分享链接
//                 imgUrl: image, // 分享图标
//                 success: onSuccess
//             });
//         });
//         document.body.appendChild(this.container);
//     }

// }

class WeixinShare extends Component {
    static jssdk = new Promise(async resolve => {
        try {
            const result = await request('/v1/config', {url: window.location.href.split('#')[0] });
            if (!result.success) {
                throw result.message;
            }
            const params = result.payload.weixin;
            wx.config({
                debug: false,
                appId: params.appId, // 必填，公众号的唯一标识
                timestamp: params.timestamp, // 必填，生成签名的时间戳
                nonceStr: params.nonceStr, // 必填，生成签名的随机串
                signature: params.signature, // 必填，签名，见附录1
                jsApiList: [
                    'checkJsApi', //判断当前客户端版本是否支持指定JS接口
                    'onMenuShareTimeline', //分享给好友
                    'onMenuShareAppMessage', //分享到朋友圈
                    'onMenuShareQQ' //分享到QQ
                ]
            });
            wx.ready(() => {
                resolve(true);
            })
        } catch(e) {
            resolve(false)
        }
    })

    static share = async (props) => {
        await WeixinShare.jssdk;

        const {
            title,
            content,
            image,
            url,
            onSuccess = () => {},
            onCancel= () => {}
        } = props;

        wx.onMenuShareTimeline({ //分享到朋友圈
            title: title,
            desc: content, // 分享描述
            link: url,
            imgUrl: image,
            success: onSuccess
        });

        wx.onMenuShareAppMessage({ //分享到朋友
            title: title, // 分享标题
            desc: content, // 分享描述
            link: url, // 分享链接
            imgUrl: image, // 分享图标
            success: onSuccess
        });

        return <Overlay onClick={onCancel}>
            <Tips/>
        </Overlay>
    }
}

class BrowserShare extends Component {
    static share = async (props) => {
        Toast.info('浏览器无法分享');
        console.log(props)
        return null;
    }
}

class AnfengGameShare extends Component {
    static share = async (props) => {
        const {title, content, image, url} = props;
        const data = {
            title,
            content,
            img: image,
            url,
            ext: { hack: 'NOT_EMPTY' }
        }
        window.anfeng.share(JSON.stringify(data))
        return null;
    }
}


class Share extends Component {
    state = {
        children: null
    }

    async share(props) {
        console.log('share.start');
      const children = await this.shareFunc().share({...props, onCancel:() => {
          this.setState({
              children: null
          })
      }})
      console.log('share', children);
      this.setState({
          children
      })
    }

    shareFunc() {
        switch (true) {
            case env.ANFENG_GAME: return AnfengGameShare;
            case env.WEIXIN: return WeixinShare;
            default: return BrowserShare;
        }
    }

    render() {
        return this.state.children;
    }
}



const container = document.createElement("div");
document.body.appendChild(container);

const share = React.createRef();

ReactDOM.render(<Share ref={share}/>, container);


export default share.current;

