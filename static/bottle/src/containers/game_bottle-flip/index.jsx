import React, { Component } from 'react';

import styled from 'styled-components';

import { Redirect } from 'react-router-dom';

import { Toast } from 'antd-mobile'

import BottleFlip from './game';

import Player from './player';
import ScoreBoard from './score-board';
import Card from './card';

import client from  '../../client';
import { AuthContext } from '../../context';

import { vw } from '../../utils';

import { getCards, useCard } from '../../api';

const gameId = 'bottle-flip';


const Mine = styled(Player).attrs({mine: true})`
  position: absolute;
  left: 0;
  top: 0;
`

const Opponent = styled(Player).attrs({mine: false})`
  position: absolute;
  right: 0;
  top: 0;
`;


const UI = styled.div`
  position: absolute;
  left: 0;
  right: 0;
  top: 0;
  bottom: 0;
`

const TopScoreBoard = styled(ScoreBoard)`
  position: absolute;
  top: 0;
  left: 0;
  right: 0;
`;

const RBCard = styled(Card)`
  position: absolute;
  right: ${ vw(20) };
  bottom: ${ vw(20) };
`;

const CARD_ID = 1;

function* generateScore() {
  let score = 0;
  yield score;
  while(true) {
    yield  new Promise(resolve => {
      setTimeout(() => {
        resolve(score += Math.ceil(Math.random() * 2))
      }, 1000 + Math.random() * 3000);
    }) 
  }
}

class BottleFlipGame extends Component {

  TIME_LIMIT = 40

  state = {
    ending: false,
    card_num: 0,
    fall: false,
    mine: 0,
    opponent: 0,
    countdown: this.TIME_LIMIT
  }

  wrapper = React.createRef();
  game = new BottleFlip();

  timer = null;
  start = new Date();

  handleRoomMessage = ({data}) => {
    const { profile } = this.props;

    if (data.uid !== profile.uid) {
      this.setState({
          opponent: data.score
      })
    }
  }

  handleGameScore = (event) => {
    const { profile, params } = this.props;

    const score = event.score;
    client.call('room_message',  {room: params.room, game_id: gameId, data:{score, uid: profile.uid} })
    this.setState({
      mine: score
    })
  }

  showResult(data) {
    this.setState({
      ending: true
    })

    if (this.props.onResult) {
      this.props.onResult(data);
    }
  }

  handleGameResult = (result) => {
    const {params } = this.props;
    const data = {
      game_id: gameId,
      room: 'BOTTLE_FLIP:X',
    }
    this.showResult({...data, ...result, players: params.info});
  }

  handleGameFall = () => {
    this.setState({
      fall: true
    })
  }

  handleTick = async() => {
    const { profile, params } = this.props;

    const countdown = this.TIME_LIMIT - Math.ceil((new Date() - this.start) / 1000);
    if (countdown >= 0) {
      this.setState({
        countdown: countdown
      })

      if (countdown == 0) {
        clearInterval(this.timer);
        const data = {
          game_id: gameId,
          uid: profile.uid,
          value: this.state.mine,
          text: this.state.mine.toString(), 
          extra: {},
          room: params.room
        }

        const {success, result, message} = await client.call('game_result', data)
        
        if (success) {
          this.showResult({...data, ...result, players: params.info})
        }
      }
    }
  }


  handleCardClick = async() => {
    if (!this.state.fall || this.state.card_num === 0) return;

    const { profile } = this.props;

    Toast.loading('复活中...', 0);
    const resp = await useCard({uid: profile.uid, card_id: CARD_ID});
    Toast.hide();

    if (!resp.success) {
      Toast.info(resp.message);
      return;
    }

    this.setState({
      card_num: resp.payload.surplus
    });

    if (this.game.revive()) {
      this.setState({
        fall: false
      })
    }
  }


  async componentDidMount() {
    const { profile } = this.props;

    this.game.start(Math.random() * 10000);
    this.wrapper.current.appendChild(this.game.renderer.domElement);
    this.timer = setInterval(this.handleTick, 1000);

    this.game.addEventListener('score', this.handleGameScore);
    this.game.addEventListener('fall', this.handleGameFall);

    client.on('notify.room_message', this.handleRoomMessage);
    client.on('notify.game_result', this.handleGameResult);

    let response = await getCards({uid: profile.uid});
    if (response.success) {
      response.payload.forEach(card => {
        if (card.card_id === CARD_ID ) {
          this.setState({
            card_num: card.card_num
          })
        }
      })
    }
    
  }

  surrender() {
    const {profile, params} = this.props;
    client.push('surrender', {game_id: gameId, uid: profile.uid, room: params.room});
  }
  

  componentWillUnmount() {
    this.game.stop();
    clearInterval(this.timer);
    this.game.removeEventListener('score', this.handleGameScore);
    this.game.removeEventListener('fall', this.handleGameFall);
    
    client.off('notify.room_message', this.handleRoomMessage);
    client.off('notify.game_result', this.handleGameResult);

    this.wrapper.current.removeChild(this.game.renderer.domElement);

    if (!this.state.ending) {
      this.surrender()
    }
  }

  render() {
    const { profile, params } = this.props;
    const players = params.info,
          opponent = players.filter(p => p.uid != profile.uid)[0];
    
    const {fall, card_num} = this.state;

    let card_type = '';
    if (card_num === 0) {
      card_type = 'invite';
    } else if (fall) {
      card_type = 'use';
    } else {
      card_type = 'normal';
    }

    return <UI>
      <div style={{ width: '100vw', height: '100vh' }} ref={this.wrapper}></div>
      <TopScoreBoard mine={this.state.mine} opponent={this.state.opponent} time={this.state.countdown}/>
      <Mine name={profile.nick_name} avatar={profile.avatar}/>
      <Opponent name={opponent.nick_name} avatar={opponent.avatar}/>
      <RBCard type={card_type} num={card_num} onClick={this.handleCardClick}/>
    </UI>;
  }
}

//1.时间间隔 2.获得分数
class BottleFlipAI extends Component {
  TIME_LIMIT = 40;
  
  score = generateScore();
  scoreStoped = false;

  state = {
    card_num: 0,
    fall: false,
    mine: 0,
    opponent: 0,
    countdown: this.TIME_LIMIT
  }

  wrapper = React.createRef();
  game = new BottleFlip();

  timer = null;
  start = new Date();

  handleGameScore = (event) => {
    const { profile, params } = this.props;
    const score = event.score;
    this.setState({
      mine: score
    })
  }

  handleGameFall = () => {
    this.setState({
      fall: true
    })
  }

  async startAIScore() {
    while(true) {
      const score = await this.score.next().value;
      if (this.scoreStoped) return;
      this.setState({opponent: score})
    }
  }

  endAIScore() {
    this.scoreStoped = true;
  }

  handleTick = async() => {
    const { profile, params } = this.props;

    const countdown = this.TIME_LIMIT - Math.ceil((new Date() - this.start) / 1000);
    if (countdown >= 0) {
      this.setState({
        countdown: countdown
      })

      if (countdown == 0) {
        clearInterval(this.timer);

        let _result = '';
        if (this.state.mine > this.state.opponent) {
          _result = 'win';
        } else if (this.state.mine === this.state.opponent) {
          _result = 'draw';
        } else {
          _result = 'lose';
        }

        const data = {
          game_id: gameId,
          uid: profile.uid,
          value: this.state.mine,
          text: this.state.mine.toString(), 
          extra: {},
          result: _result,
        }

        const {success, result} = await client.call('game_result_ai', data)
        
        if (success) {
          if (this.props.onResult) {
            this.props.onResult({...data, ...result, players: params.info});
          }
        }
      }
    }
  }

  handleCardClick = async() => {
    if (!this.state.fall || this.state.card_num === 0) return;

    const { profile } = this.props;

    Toast.loading('复活中...', 0);
    const resp = await useCard({uid: profile.uid, card_id: CARD_ID});
    Toast.hide();

    if (!resp.success) {
      Toast.info(resp.message);
      return;
    }

    this.setState({
      card_num: resp.payload.surplus
    });

    if (this.game.revive()) {
      this.setState({
        fall: false
      })
    }
  }

  async componentDidMount() {
    const { profile } = this.props;

    this.game.start(Math.random() * 10000);
    this.startAIScore();
    this.wrapper.current.appendChild(this.game.renderer.domElement);
    this.timer = setInterval(this.handleTick, 1000);
    this.game.addEventListener('score', this.handleGameScore)
    this.game.addEventListener('fall', this.handleGameFall);
    
    let response = await getCards({uid: profile.uid});
    if (response.success) {
      response.payload.forEach(card => {
        if (card.card_id === CARD_ID ) {
          this.setState({
            card_num: card.card_num
          })
        }
      })
    }

  }

  componentWillUnmount() {
    this.game.stop();
    this.endAIScore();
    clearInterval(this.timer);
    this.game.removeEventListener('score', this.handleGameScore);
    this.game.removeEventListener('fall', this.handleGameFall)
    this.wrapper.current.removeChild(this.game.renderer.domElement);
  }

  render() {
    const { profile, params } = this.props;
    const players = params.info,
          opponent = players.filter(p => p.uid != profile.uid)[0];
    
    const {fall, card_num} = this.state;

    let card_type = '';
    if (card_num === 0) {
      card_type = 'invite';
    } else if (fall) {
      card_type = 'use';
    } else {
      card_type = 'normal';
    }

    return <UI>
      <div style={{ width: '100vw', height: '100vh' }} ref={this.wrapper}></div>
      <TopScoreBoard mine={this.state.mine} opponent={this.state.opponent} time={this.state.countdown}/>
      <Mine name={profile.nick_name} avatar={profile.avatar}/>
      <Opponent name={opponent.nick_name} avatar={opponent.avatar}/>
      <RBCard type={card_type} num={card_num} onClick={this.handleCardClick}/>
    </UI>;
  }
}


export default class Wrapper extends Component {

  handleResult = params => {
    this.props.history.push({
      pathname: '/ending',
      state: params
    })
  }

  render() {
      const params = this.props.location.state;
      if (!params) {
        return <Redirect to="/" />;
      }


      return <AuthContext.Consumer>
        {
            ({profile}) => (
              !params.ai ? <BottleFlipGame profile={profile} params={params} onResult={this.handleResult}/>
                : <BottleFlipAI profile={profile} params={params} onResult={this.handleResult}/>
            )
        }  
      </AuthContext.Consumer>;
  }
}