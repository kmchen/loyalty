import React, { Component } from 'react';
import { connect } from 'react-redux';
import Dashboard from './Dashboard';
import UserSearchBox from './UserSearchBox';
import { riderRecordRequest } from '../action/actionTypes';
import logo from '../asset/logo.svg';
import '../asset/App.css';

class App extends Component {
    constructor(props) {
        super(props);
        this.state = { userId: '' };
        this.handleTextChange = this.handleTextChange.bind(this);
    }

    handleTextChange(e) {
        if (e.keyCode === 13) {
            this.props.getUserData(this.state.userId)
        } else {
            this.setState({ userId: e.target.value });
        }
    }

  render() {
      const showDashboard = this.props.riderData ? {
          id: this.props.riderData.id,
          name: this.props.riderData.name,
          loyalty: this.props.riderData.loyalty,
          grade: this.props.riderData.grade,
          numRides: this.props.riderData.numRides
      } : null;
      return (
          <div className="App">
              <header className="App-header">
                  <img src={logo} className="App-logo" alt="logo" />
                  <h1 className="App-title">Loyalty Dashboard</h1>
              </header>
              <section>
                  <UserSearchBox
                    text={this.state.text}
                    username={this.state.username}
                    handleTextChange={this.handleTextChange}
                  />
                  { Boolean(showDashboard) &&
                      <Dashboard {...showDashboard}/>
                  }
              </section>
          </div>
      );
  }
}

const mapStateToProps = state => {
    return { 
      riderData: state.message.riderData
    }
}

const mapDispatchToProps = dispatch => {
    return {
        getUserData: userId => {
            dispatch({type: riderRecordRequest, userId: userId})
        }
    }
}

const RootApp = connect(
    mapStateToProps,
    mapDispatchToProps
)(App);

export default RootApp;
