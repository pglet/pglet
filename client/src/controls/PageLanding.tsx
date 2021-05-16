import React from 'react';
import { useParams } from "react-router-dom";
import { useSelector } from 'react-redux';
import { Page } from './Page'
import { Signin } from './Signin'
import { WebSocketContext } from '../WebSocket';
import { mergeStyles, MessageBar, MessageBarType, PartialTheme, Spinner, ThemeProvider } from '@fluentui/react'
import { getWindowHash } from './Utils'
import { buildTheme, darkThemeColor } from './Theming';

interface ParamTypes {
    accountName: string,
    pageName: string
}

export const PageLanding = () => {

    let { accountName, pageName } = useParams<ParamTypes>();
    const [theme, setTheme] = React.useState<PartialTheme | undefined>();
    const ws = React.useContext(WebSocketContext);

    if (!accountName) {
        accountName = "p";
    }

    if (!pageName) {
        pageName = "index";
    }

    let fullPageName = `${accountName}/${pageName}`;

    const updateTheme = (themePrimaryColor:any, themeTextColor:any, themeBackgroundColor:any) => {
        var theme = buildTheme(themePrimaryColor, themeTextColor, themeBackgroundColor)
        setTheme(theme);
        document.documentElement.style.background = themeBackgroundColor;
    }

    React.useEffect(() => {
        updateTheme(darkThemeColor.primary, darkThemeColor.text, darkThemeColor.background);
        ws.registerWebClient(fullPageName, getWindowHash());

    }, [fullPageName, ws])

    const err = useSelector((state: any) => state.page.error);
    const signinOptions = useSelector((state: any) => state.page.signinOptions);
    const page = useSelector((state: any) => state.page.controls['page']);

    const className = mergeStyles({
        height: '100vh'
      });

    const renderContent = () => {
        if (err === "signin_required") {
            return <Signin signinOptions={signinOptions} />
        }
        else if (err) {
            return <MessageBar messageBarType={MessageBarType.error} isMultiline={true}>{err}</MessageBar>
        } else if (!page) {
            return <Spinner label="Loading page, please wait..." labelPosition="right" styles={{ root: { height: "35px" }}} />
        } else {
            return <Page control={page} pageName={fullPageName} updateTheme={updateTheme} />
        }
    }

    return <ThemeProvider theme={theme} className={className}>
        {theme && renderContent()}
        </ThemeProvider>
}