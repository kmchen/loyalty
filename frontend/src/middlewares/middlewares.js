import { send, client } from '../websocket';
import { riderRecordResponse, riderRecordRequest } from '../action/actionTypes';

export const riderResponseMiddleware = store => next => action => {
    client.onmessage = function(e) {
        const {type, rider} = JSON.parse(e.data);
            switch(type) {
                case riderRecordResponse:
                const riderResp = {type: riderRecordResponse, riderData: rider};
                store.dispatch(riderResp);
                next(riderResp);
                break;
              default:
                break;
        }
    }
}

export const riderRequestMiddleware = store => next => action => {
    switch (action.type) {
        case riderRecordRequest:
            send({type:riderRecordRequest, userId: action.userId})
            next(action);
            break;
        default:
            break
    }
};
