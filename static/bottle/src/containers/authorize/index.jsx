import React, { Component } from 'react';
import styled from 'styled-components';

import { Redirect } from 'react-router-dom';

import client from '../../client';

import {  AuthContext } from '../../context';
import { vw } from '../../utils';

const Background = styled.div`
    background: rgba(242,242,242,1);
    min-height: 100vh;
    display: flex;
    justify-content: center;
    align-items: center;
`;

const Loading = styled.div`
    width: ${vw(640)};
    height: ${vw(360)};
    background-image: url(${require('./loading.gif')});
    background-size: cover;
`;

class Authorize extends Component {

    async redirect() {
        const {success, result} = await client.call('authorize', {})
        if (success) {
            window.location.href = result.url;
        }
    }

    async componentDidMount() {
        const accessToken = this.props.auth.accessToken;

    
        if (accessToken) {
            localStorage.setItem('accessToken', accessToken);

            const {success, result} = await client.call('login', {access_token: accessToken})
            
            if (success) {
                this.props.auth.update(result);
                const authorizeFrom = localStorage.getItem('authorizeFrom') || '/';
                window.location.href = authorizeFrom;
            } else {
                localStorage.removeItem('accessToken');
                this.redirect();
            }

        } else {
            this.redirect();
        }
    }

    render() {
        return <Background>
            <Loading/>
        </Background>;
    }
}


export default class Wrapper extends Component {
    state= {
        auth: null
    }

    render() {
        console.log(this.props);
        const { accessToken = null} = this.props.match.params;
        const { from = null } = this.props.location.state || {};
        if (from) {
            localStorage.setItem('authorizeFrom', from);
        }
        return (
            <AuthContext.Consumer>
            {
                auth => <Authorize auth={{...auth, accessToken: accessToken || auth.accessToken}}/>
            }
            </AuthContext.Consumer>
        );
    }
}