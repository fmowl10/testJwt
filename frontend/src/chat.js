import React from 'react'
import 'semantic-ui-css/semantic.min.css'
import {Button, Input, Form} from 'semantic-ui-react'
import {w3cwebsocket as W3CWebSocket} from "websocket";

class Chat extends React.Component {
    constructor(props) {
        super(props)
        this.state = {token : this.props.location.state.token, text : '', message : '', enable : false};
        this.socket = new W3CWebSocket('wss://test.fmowl.com/ws', this.state.token);
    }
    componentWillMount() {
        this.socket.onopen = () => {
            this.setState({text :'Websocket connected'});
            this.setState({enable: true})
        }
        this.socket.onmessage = (message) => {
            this.setState({text : this.state.text +" \n"+ message.data})
        }
        this.socket.onclose = () => {
            this.setState({enable: false, text:'disconneted'})
        }
    }
    handleSubmit = () =>{
        this.socket.send(this.state.message)
    }
    handleChange = (e, {name, value}) => {
        this.setState({[name]: value});
    }
    render() {
        return (
            <div style={{backgroundColor:"#6b7572", padding:10, position:"absolute", top:0, right:0, bottom:0, left:0}}>
                <div style={{backgroundColor:"#f7f7f7", padding:10, width:"100%", height:"100%"}}>
                    <p style={{wordBreak:"break-all", wordWrap:"break-word", whiteSpace:"pre", width:200, heigth:"100%"}}>hello world<br/>{this.state.text}</p>
                    <div style={{position:"relative",bottom:0, width:"100%"}}>
                        <Form onSubmit={this.handleSubmit} >
                            <div>
                                <Input action="Submit" fluid disabled={!this.state.enable} onChange={this.handleChange} name="message" placeholder="input message"/>
                            </div>
                        </Form>
                    </div>
                </div>
            </div>
        );
    }
}

//<!--<Button style={{width:"30%"}} isabled={!this.state.enable} content="Submit"/>-->
export default Chat;