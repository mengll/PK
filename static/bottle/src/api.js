import { request } from './utils';


export const getCards = (params) => request('/v1/get_card', params);

export const useCard = (params) => request('/v1/use_card', params);