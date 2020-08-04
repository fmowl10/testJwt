import React from 'react'
import 'semantic-ui-css/semantic.min.css'
import {Grid, Form, Button, Message} from 'semantic-ui-react'

class Input extends React.Component {
    constructor(props) {
        super(props);
        this.state = {key:'', role:'', token:'', value : '', error:false};
        this.isActive = true;
        this.options = [{key:'local', value:'local', text:'local'}, 
          {key:'host', value:'host', text:'host'}];
    }

    postInput= (url='', data='') =>  {
      return fetch(url, {
        method: 'POST',
        headers: {
          'Content-Type' : 'application/json'
        },
        body : JSON.stringify(data)
      }).then(response =>  { 
            if (response.ok) {
                return response.json()
            } else {
                console.log(response.text())
                return "not working"
            }
      });
    }


    getData = (url= '', token='') => {
      return fetch(url, {
        method : 'Get',
        headers: {
          'Authorization' : "Bearer " + token
        }
      })
      .then(response => response.text())
      .catch(error =>{console.log(error); return "not working"});
    }

    handleSubmit = async () => {
        const value = {key: this.state.key, role : this.state.role}
        if (value.key.length <= 0 || value.role.length <= 0) {
            this.setState({error:true})
            return ;
        }
        const data = await this.postInput('/jwt', value)
        this.setState({token: data.token, error: false});
    }

    onChatButtonClick = () => {
      this.props.history.push({
        pathname: '/chat',
        state: {token: this.state.token, key: this.state.key}
      });
    }

    handleChange = (event, {name, value}) => {
        this.setState({[name]: value});
    }

    onApiButtonClicked = () => {
        this.setState({value: ''})
        this.getData("/api", this.state.token)
        .then(data => {
            this.setState({value: data})
        })
    }
    

    render() {
        return (
          <div style={
            {width:400, border:"5px solid black", borderRadius:20, padding:10}}>
            <Grid columns={2}>
              <Grid.Column>
              <Form onSubmit={this.handleSubmit} error={this.state.error}>
                {this.state.error && <Message error header="input both thing"/>}
                <Form.Field>
                  <label>Key</label>
                  <Form.Input 
                    name="key" 
                    onChange={this.handleChange} 
                    placeholder='input key'
                  />
                  <Form.Dropdown 
                    name="role"
                    onChange={this.handleChange} 
                    placeholder="pick a role" 
                    options={this.options}
                    selection
                    clearable
                    />
                  </Form.Field>
                  <Form.Button content='Submit'/>
              </Form>
              </Grid.Column>
              <Grid.Column>
                {this.state.token !== 'not working' && this.state.token && 
                    <div>
                        <div style={
                          {width:180, wordBreak:"break-all", wordWrap:"break-word", }}>
                            {this.state.token}</div>
                        <Button name="apiButton"onClick={this.onApiButtonClicked}>
                            Get Hello
                        </Button>
                        <Button name="chatButton"onClick={this.onChatButtonClick}>
                        	Go Chat</Button>
                        <h1 id='api'>{this.state.value}</h1>
                    </div>
                }
              </Grid.Column>
            </Grid>
            </div>
        ); 
    }
}

export default Input;