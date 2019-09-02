import React from 'react'
import ReactDOM from 'react-dom'
import {BrowserRouter as Router, Route, Switch} from "react-router-dom"
import Input from './input.js'
import Chat from './chat.js'

class App extends React.Component {
    render() {
        return (
                <Router>
                    <Switch>
                    <Route exact path="/" component={Input} />
                    <Route path="/chat" component={Chat}/>
                    </Switch>
                </Router>
        );
    }
}
ReactDOM.render(<App />, document.getElementById('root'));