import React, { useEffect, useContext } from 'react';
import { useParams } from "react-router-dom";
import { useSelector } from 'react-redux';
import { Page } from './Page'
import { WebSocketContext } from '../WebSocket';
import { MessageBar, MessageBarType } from '@fluentui/react'

interface ParamTypes {
    accountName: string,
    pageName: string
}

export const PageLanding = () => {

    let { accountName, pageName } = useParams<ParamTypes>();

    if (!accountName) {
        accountName = "public";
    }

    if (!pageName) {
        pageName = "index";
    }

    let fullPageName = `${accountName}/${pageName}`;

    const ws = useContext(WebSocketContext);

    useEffect(() => {

        ws.registerWebClient(fullPageName);

    }, [fullPageName, ws])

    const err = useSelector((state: any) => state.page.error);
    const root = useSelector((state: any) => state.page.controls['page']);

    if (err) {
        return <MessageBar messageBarType={MessageBarType.error} isMultiline={false}>{err}</MessageBar>
    } else {
        return <Page control={root} />
    }
}