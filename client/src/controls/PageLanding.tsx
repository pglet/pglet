import React from 'react';
import { useParams } from "react-router-dom";
import { useSelector } from 'react-redux';
import { Page } from './Page'
import { Login } from './Login'
import { WebSocketContext } from '../WebSocket';
import { MessageBar, MessageBarType, Spinner } from '@fluentui/react'
import { getWindowHash } from './Utils'

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

    const ws = React.useContext(WebSocketContext);

    React.useEffect(() => {

        ws.registerWebClient(fullPageName, getWindowHash());

    }, [fullPageName, ws])

    const err = useSelector((state: any) => state.page.error);
    const page = useSelector((state: any) => state.page.controls['page']);

    if (err === "login_required") {
        return <Login />
    }
    else if (err) {
        return <MessageBar messageBarType={MessageBarType.error} isMultiline={true}>{err}</MessageBar>
    } else if (!page) {
        return <Spinner label="Loading page, please wait..." labelPosition="right" styles={{ root: { height: "35px" }}} />
    } else {
        return <Page control={page} pageName={fullPageName} />
    }
}