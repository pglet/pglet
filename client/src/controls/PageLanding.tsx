import React from 'react';
import { useParams } from "react-router-dom";
import { useSelector } from 'react-redux';
import { Page } from './Page'
import { Signin } from './Signin'
import { WebSocketContext } from '../WebSocket';
import { mergeStyles, MessageBar, MessageBarType, Spinner, Theme, ThemeProvider } from '@fluentui/react'
import { getWindowHash } from './Utils'
import { buildTheme, getThemeExtraStyles } from './Theming';

interface ParamTypes {
    accountName: string,
    pageName: string
}

export const PageLanding = () => {

    let { accountName, pageName } = useParams<ParamTypes>();
    const [theme, setTheme] = React.useState<Theme | undefined>();
    const ws = React.useContext(WebSocketContext);

    if (!accountName) {
        accountName = "p";
    }

    if (!pageName) {
        pageName = "index";
    }

    let fullPageName = `${accountName}/${pageName}`;

    const updateTheme = (standardTheme:any, themePrimaryColor?:any, themeTextColor?:any, themeBackgroundColor?:any) => {
        var theme = buildTheme(standardTheme, themePrimaryColor, themeTextColor, themeBackgroundColor)
        setTheme(theme);
        document.documentElement.style.background = theme.semanticColors.bodyBackground;
    }

    React.useEffect(() => {
        updateTheme('dark');
        ws.registerWebClient(fullPageName, getWindowHash());

    }, [fullPageName, ws])

    const err = useSelector((state: any) => state.page.error);
    const signinOptions = useSelector((state: any) => state.page.signinOptions);
    const page = useSelector((state: any) => state.page.controls['page']);

    const className = mergeStyles({
        height: '100vh',
      }, getThemeExtraStyles(theme));

    const renderContent = () => {
        if (err === "signin_required") {
            return <Signin signinOptions={signinOptions} />
        }
        else if (err) {
            return <MessageBar messageBarType={MessageBarType.error} isMultiline={true}>{err}</MessageBar>
        } else if (!page) {
            return <Spinner label="Loading, please wait..." labelPosition="right" styles={{ root: { height: "100vh" }}} />
        } else {
            return <Page control={page} pageName={fullPageName} updateTheme={updateTheme} />
        }
    }

    return <ThemeProvider theme={theme} className={className}>
        {theme && renderContent()}
        </ThemeProvider>
}