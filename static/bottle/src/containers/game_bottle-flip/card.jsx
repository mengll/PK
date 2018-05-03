import React, { Component } from 'react';
import styled from 'styled-components';

import { vw } from '../../utils';

const BADGE_SIZE = 34;

const Badge = styled.div`
    position: absolute;

    padding: 0 ${ vw(4) };

    right: ${ vw(-BADGE_SIZE/2) };
    top: ${ vw(-BADGE_SIZE/2) };

    display: flex;
    justify-content: center;
    align-items: center;

    border-radius ${ vw(BADGE_SIZE/2) };
    border: ${ vw(2) } solid rgba(255,255,255,0.8);

    font-size: ${ vw(20) };

    color: #FFF;
    background-color: #FB4217;
    min-width: ${ vw(BADGE_SIZE) };
    height: ${ vw(BADGE_SIZE) };
    border-radius: ${ vw(BADGE_SIZE/2) };
`


const Image = styled.div`
    position: relative;
    
    margin: 0 auto;

    width: ${ vw(81) };
    height: ${ vw(97) };
    background-image: url(${ require('./card.png') });
    background-size: cover;
`

const InviteText = styled.div`
    width: ${ vw(191) };
    height: ${ vw(44) };
    background-image: url(${ require('./text_invite.png') });
    background-size: cover;
`;

const NormalText = styled.div`
    width: ${ vw(72) };
    height: ${ vw(25) };
    background-image: url(${ require('./text_normal.png') });
    background-size: cover;
`;

const UseText = styled.div`
    width: ${ vw(87) };
    height: ${ vw(45) };
    background-image: url(${ require('./text_use.png') });
    background-size: cover;
`;

const Wrapper = styled.div`
    width: ${ vw(120) };
`

const Text = styled.div`
    margin-top: ${ vw(10) };
    height: ${ vw(50) };
    display: flex;
    justify-content: center;
    align-items: center;
`;

const Texts = {
    'use': UseText,
    'normal': NormalText,
    'invite': InviteText,
}


export default class Card extends Component {
    render() {
        const { type = 'invite', num = 0, ...rest } = this.props;
        return (
            <div {...rest}>
                <Wrapper>
                    <Image>
                        <Badge>{num}</Badge>
                    </Image>
                    <Text>
                        {
                            React.createElement(Texts[type])
                        }
                    </Text>
                </Wrapper>
            </div>
        );
    }
}