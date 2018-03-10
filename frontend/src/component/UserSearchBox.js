import React from "react";
import '../asset/UserSearchBox.css';

export default ({ text, username, handleTextChange }) => (
  <div>
      <div className="row">
          <div className="col-xs-12">
              <div className="chat">
                  <div className="col-xs-5 col-xs-offset-3">
                      <input
                        type="text"
                        value={text}
                        placeholder="Please enter user id here ...."
                        className="form-control"
                        onChange={handleTextChange}
                        onKeyDown={handleTextChange}
                      />
                  </div>
                  <div className="clearfix"></div>
              </div>
          </div>
      </div>
    </div>
);
