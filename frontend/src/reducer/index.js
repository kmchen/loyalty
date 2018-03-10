import {combineReducers } from 'redux';
import * as actions from '../action/actionTypes';

export const message = (state = {}, action) => {
    switch(action.type) {
        case actions.riderRecordResponse:
            const {riderData} = action;
            return {...state, riderData};
        default:
            return state
    }
}

export const reducers = combineReducers({message});
