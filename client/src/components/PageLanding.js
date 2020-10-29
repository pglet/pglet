import React, { useEffect, useContext } from 'react';
import { useParams } from "react-router-dom";
import { useSelector } from 'react-redux';
import Page from './Page'
import { WebSocketContext } from '../WebSocket';

const PageLanding = () => {

    let { accountName, pageName } = useParams();

    let fullPageName = `${accountName}/${pageName}`;

    const ws = useContext(WebSocketContext);

    useEffect(() => {

        ws.registerWebClient(fullPageName);

    }, [fullPageName, ws])

    const err = useSelector(state => state.page.error);
    const root = useSelector(state => state.page.controls['_0']);

    if (err) {
        return <div>{err}</div>;
    } else {
        return <Page control={root} />
    }
}

export default PageLanding