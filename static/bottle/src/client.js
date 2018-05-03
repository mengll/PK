import EventEmitter from "eventemitter2";
import { request } from './utils';
import { Toast } from 'antd-mobile';

// START          = "af01"
// LOGIN          = "af02"
// LOGOUT         = "af03"
// CREATE_ROOM    = "af04"
// SEARCH_MATCH   = "af05"
// GAME_HEART     = "af06"
// JOIN_CANCEL    = "af07"
// ROOM_MESSAGE   = "af08"
// OUT_ROOM       = "af09"
// RECONNECT      = "af10"
// NOW_ONLINE_NUM = "af11"
// JOIN_ROOM      = "af12"
// GAME_RESULT    = "af13"
// AUTHORIZE      = "af14"
// TIME_OUT       = "af15"
// DISCONNECT     = "af16"
// ONLINE         = "af17"
// USER_MESSAGE   = "af18"
// ENTER_GAME	   = "af19"

const routes = {
    start: 'af01',
    login: 'af02',
    logout: 'af03',
    create_room: 'af04',
    search_match: 'af05',
    game_heart: 'af06',
    join_cancel: 'af07',
    room_message: 'af08',
    out_room: 'af09',
    reconnect: 'af10',
    now_online_num: 'af11',
    join_room: 'af12',
    game_result: 'af13',
    authorize: 'af14',
    timeout: 'af15',
    online: 'af17',
    user_message: 'af18',
    enter_game: 'af19',
    surrender: 'af20',
    game_result_ai: 'af21',
}

const cmds = {
    'af01': 'start',
    'af02': 'login',
    'af03': 'logout',
    'af04': 'create_room',
    'af05': 'search_match',
    'af06': 'game_heart',
    'af07': 'join_cancel',
    'af08': 'room_message',
    'af09': 'out_room',
    'af10': 'reconnect',
    'af11': 'now_online_num',
    'af12': 'join_room',
    'af13': 'game_result',
    'af14': 'authorize',
    'af15': 'timeout',
    'af16': 'disconnect',
    'af17': 'online',
    'af18': 'user_message',
    'af19': 'enter_game',
    'af20': 'surrender',
    'af21': 'game_result_ai'
}
class Client extends EventEmitter {

    socket = null;
    pending = {}
    seq = 0;

    connecting = false;

    constructor() {
        super();
        
    }

    handleMessage = (event) => {
        try {
            const pack = JSON.parse(event.data);

            const { error_code, data, msg, message_id: seq } = pack;

            const success = error_code === 0;

            const response = {
                success,
                result: data,
                message: msg
            }

            const callback = this.pending[seq];
            
            if (callback) {
                // RESPONSE
                delete this.pending[seq];
                callback(response);
            } else {
                // NOTIFY
                const method = cmds[msg];
                if (method === undefined) {
                    console.log('client.notify.unknow', pack);
                } else {
                    if (success) {
                        this.notify({
                            method,
                            params: data
                        })
                    }
                }
            }
        } catch (error) {
            console.log('client.error.parse', error)
        }
    }

    async createSocket() {
        const {success, payload, message} = await request('/v1/config');
        if (!success) {
            throw message;
        }
        let socket = null;
        socket = new WebSocket(payload.websocket);
        socket.onmessage = this.handleMessage;
        return socket;
    }

    notify(event) {
        console.log('notify', event);
        const {method, params} = event;
        this.emit(`notify.${method}`, params);
    }

    async connected() {
        if (!this.connecting) {
            this.connecting = true;
            const promise = new Promise(async resolve => {
                const retry = () => {
                    Toast.loading('正在连接..', 0);
                    setTimeout(() => {
                        this.connecting = false;
                        this.connected().then(() => {
                            window.location.href = '/';
                            resolve();
                        })
                    }, 3000)
                }

                if (this.socket && this.socket.readyState == WebSocket.OPEN) {
                    resolve()
    
                }else if (this.socket == null || this.socket.readyState != WebSocket.CONNECTING ) {
                    try {
                        this.socket = await this.createSocket();
                        this.socket.addEventListener('open', () => {
                            resolve();
                        })
                        this.socket.addEventListener('error', () => {
                            retry();
                        })
                    } catch (e) {
                        retry();
                    }
                }
            })
            promise.then(() => {
                this.connecting = false;
                this.emit('open')
            })
            return promise;
        } else {
            return new Promise(resolve => {
                this.on('open', () => {
                    resolve()
                })
            })
        }
    }






    async call(method, params) {
        console.log('client.call.' + method, params);
        await this.connected();
        
        const cmd = routes[method];
        
        if (cmd === undefined) {
            throw new Error("Unknow method!");
        }

        return new Promise(
            (resolve) => {
                ++this.seq;

                const action = {
                    cmd,
                    data: {...params},
                    message_key: "",
                    message_id: this.seq.toString()
                }

                this.pending[this.seq] = (response) => {
                    console.log('client.response', response);
                    resolve(response)
                }

                this.socket.send(JSON.stringify(action))
            }
        )
    }

    async push(method, params) {
        method != 'game_heart' && console.log('client.push', {method, params});
        await this.connected();
        
        const cmd = routes[method];
        
        if (cmd === undefined) {
            throw new Error("Unknow method!");
        }

        return new Promise(
            (resolve) => {
                ++this.seq;

                const action = {
                    cmd,
                    data: {...params},
                    message_key: "",
                    message_id: this.seq.toString()
                }

                this.socket.send(JSON.stringify(action))

                resolve();
            }
        )
    }
    
}


export {
    routes
};

export default new Client()

