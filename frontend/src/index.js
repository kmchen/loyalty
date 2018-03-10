import React from 'react';
import ReactDOM from 'react-dom';
import './asset/index.css';
import 'websocket';
import RootApp from './component/App';
import logger from 'redux-logger'
import { createStore, applyMiddleware } from 'redux';
import { reducers } from './reducer/index';
import { Provider } from 'react-redux';
import { riderRequestMiddleware, riderResponseMiddleware } from './/middlewares/middlewares';

const store = createStore(
    reducers,
    applyMiddleware(logger, riderRequestMiddleware, riderResponseMiddleware)
);

ReactDOM.render(
    <Provider store={store}>
        <RootApp />
    </Provider>
    , document.getElementById('root'));
