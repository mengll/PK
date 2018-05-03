import { request } from './utils';


let card_num = 3;

export const getCards = (params) => {
    return {
        success: true,
        payload: [
            {
                name: "复活卡",
                card_num: card_num,
                card_id: 1
            }
        ]
    }
}



export const useCard = (params) => {
    return {
        success: true,
        payload: {
            surplus: -- card_num
        }
    }
}