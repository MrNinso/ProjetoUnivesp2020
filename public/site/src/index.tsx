import React from 'react'
import ReactDOM from 'react-dom'
import { Route, Switch, BrowserRouter} from 'react-router-dom'
import './index.css'
import * as serviceWorker from './serviceWorker'

import { LoginPage } from './Pages/Login/Login.page'
import { RoomsPage } from "./Pages/Rooms/Rooms.page"

var Sw: Switch | null

ReactDOM.render(
    <React.StrictMode>
        <BrowserRouter>
            <div>
                <Switch ref={(s) => Sw = s }>
                    <Route path="/app/login" component={ () => LoginPage(Sw) }  />
                    <Route path="/app/rooms" component={ RoomsPage }  />
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
