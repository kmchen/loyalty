import React from "react";
import "../asset/Dashboard.css";
import avatar from "../asset/avatar.png";

export default ({ id, name, loyalty, grade, numRides }) => {
    return (
        <ul>
            <div className="row show-grid">
                <div className="col-xs-12">
                    <div className="chatMessage">
                      <div key={id} className="box">
                        <p><strong>{name}</strong></p>
                        <p> Grade : {grade}</p>
                        <p> Loyaylty : {loyalty}</p>
                        <p> Number of rides : {numRides}</p>
                      </div>
                      <div className="imageHolder">
                      <img src={avatar}
                          className="img-responsive avatar" alt="logo" />
                    </div>
                    </div>
                  </div>
            </div>
        </ul>
    );
};
