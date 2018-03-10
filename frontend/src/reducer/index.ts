import {combineReducers } from 'redux';
import * as actions from '../action/actionTypes';

export const message = (state = {}, action) => {
  switch(action.type) {
    case actions.RIDER_RECORD_RESPONSE:
      //const {userId} = action.payload;
      return state;
    //case actions.MESSAGE_USER_JOINED:
      //return {...state, newUser: action.payload.userId}
    //case actions.MESSAGE_STROKE:
      //return {...state, stroke: action}
    //default:
      //return state
  }
}

export const reducers = combineReducers({message});
