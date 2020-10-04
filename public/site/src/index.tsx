import React from 'react'
import ReactDOM from 'react-dom'
import { Route, Switch, BrowserRouter, useHistory } from 'react-router-dom'
import './index.css'
import * as serviceWorker from './serviceWorker'


import { LoginPage } from './Pages/Login/Login.page'
import { RoomsPage } from "./Pages/Rooms/Rooms.page"

var http = new XMLHttpRequest()

ReactDOM.render(
    <React.StrictMode>
        <BrowserRouter>
            <div>
                <Switch>
                    <Route path="/app/login" component={  () => LoginPage(http, useHistory()) }  />
                    <Route path="/app/rooms" component={ () => RoomsPage(http) }  />
                </Switch>
            </div>
        </BrowserRouter>
    </React.StrictMode>,
    document.getElementById('root')
);
//TODO <Route component={Error} />

// If you want your app to work offline and load faster, you can change
// unregister() to register() below. Note this comes with some pitfalls.
// Learn more about service workers: https://bit.ly/CRA-PWA
serviceWorker.unregister();
