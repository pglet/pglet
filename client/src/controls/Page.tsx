import React from 'react'
import { shallowEqual, useSelector, useDispatch } from 'react-redux'
import { ControlsList } from './ControlsList'
import { Stack, IStackProps, IStackTokens, Theme, MessageBar, Spinner, MessageBarType, ThemeProvider, mergeStyles } from '@fluentui/react';
import { Signin } from './Signin'
import { ISigninProps } from './Control.types'
import { WebSocketContext } from '../WebSocket';
import { changeProps } from '../slices/pageSlice'
import { defaultPixels, getThemeColor, getWindowHash, isFalse, isTrue } from './Utils'
import { useParams } from 'react-router-dom';
import { buildTheme, getThemeExtraStyles } from './Theming';

interface ParamTypes {
  accountName: string,
  pageName: string
}

export const Page = () => {

  const ws = React.useContext(WebSocketContext);
  const dispatch = useDispatch();
  const defaultTheme = buildTheme('light', null, null, null);

  let { accountName, pageName } = useParams<ParamTypes>();
  const [prevStandardTheme, setPrevStandardTheme] = React.useState<any | undefined>();
  const [prevThemePrimaryColor, setPrevThemePrimaryColor] = React.useState<any | undefined>();
  const [prevThemeTextColor, setPrevThemeTextColor] = React.useState<any | undefined>();
  const [prevThemeBackgroundColor, setPrevThemeBackgroundColor] = React.useState<any | undefined>();
  const [theme, setTheme] = React.useState<Theme | undefined>();

  const err = useSelector((state: any) => state.page.error);
  const signinOptions = useSelector((state: any) => state.page.signinOptions);
  const control = useSelector((state: any) => state.page.controls['page']);
  const childControls = useSelector((state: any) => {
    const page = state.page.controls['page'];
    return page ? page.c.map((childId: string) => state.page.controls[childId]) : []
  }, shallowEqual);

  if (!accountName) {
    accountName = "p";
  }

  if (!pageName) {
    pageName = "index";
  }

  let fullPageName = `${accountName}/${pageName}`;

  const updateTheme = (standardTheme?: any, themePrimaryColor?: any, themeTextColor?: any, themeBackgroundColor?: any) => {

    standardTheme = standardTheme ?? "light";

    if (standardTheme !== prevStandardTheme ||
      themePrimaryColor !== prevThemePrimaryColor ||
      themeTextColor !== prevThemeTextColor ||
      themeBackgroundColor !== prevThemeBackgroundColor) {

      // build theme
      var theme = buildTheme(standardTheme, themePrimaryColor, themeTextColor, themeBackgroundColor)
      setTheme(theme);
      document.documentElement.style.background = theme.semanticColors.bodyBackground;
    }

    setPrevStandardTheme(standardTheme);
    setPrevThemePrimaryColor(themePrimaryColor);
    setPrevThemeTextColor(themeTextColor);
    setPrevThemeBackgroundColor(themeBackgroundColor);
  }

  const data = {
    fireUpdateHashEvent: true
  }

  function updateHash(hash: string) {
    if (data.fireUpdateHashEvent) {
      const payload: any = {
        i: "page",
        hash: hash
      }

      dispatch(changeProps([payload]));
      ws.updateControlProps([payload]);
      ws.pageEventFromWeb("page", 'hashChange', hash);
    }

    data.fireUpdateHashEvent = true;
  }

  React.useEffect(() => {

    if (!control) {
      ws.registerWebClient(fullPageName, getWindowHash());
      return;
    }

    // page title
    let title = `${fullPageName} - Pglet`;
    if (control.title) {
      title = control.title
    }
    document.title = title;

    // handle resize
    let resizeTimeout: any = null;
    function handleResize() {
      clearTimeout(resizeTimeout);
      resizeTimeout = setTimeout(() => {
        //console.log("window size:", window.innerHeight, window.innerWidth);
        const payload: any = [{
          i: "page",
          win_width: String(window.innerWidth)
        },
        {
          i: "page",
          win_height: String(window.innerHeight)
        }];

        dispatch(changeProps(payload));
        ws.updateControlProps(payload);
        ws.pageEventFromWeb("page", 'resize', `${window.innerWidth} ${window.innerHeight}`);
      }, 250)
    }

    // theme
    updateTheme(control.theme, control.themeprimarycolor, control.themetextcolor, control.themebackgroundcolor);

    const hash = getWindowHash();
    const pageHash = control.hash !== undefined ? control.hash : "";

    if (pageHash !== hash) {
      window.location.hash = pageHash ? "#" + pageHash : "";
      data.fireUpdateHashEvent = false;
    }

    const handleHashChange = (e: any) => {
      updateHash(getWindowHash());
    }

    window.addEventListener("hashchange", handleHashChange);
    window.addEventListener('resize', handleResize);

    return () => {
      window.removeEventListener("hashchange", handleHashChange);
      window.removeEventListener("resize", handleResize);
    }
    // eslint-disable-next-line
  }, [control, ws]);

  const className = mergeStyles({
    height: '100vh',
  }, getThemeExtraStyles(theme));

  const renderPage = () => {

    if (!theme || isFalse(control.visible)) {
      return null;
    }

    let disabled = isTrue(control.disabled)

    // stack props
    const stackProps: IStackProps = {
      verticalFill: control.verticalfill ? isTrue(control.verticalfill) : undefined,
      horizontalAlign: control.horizontalalign === '' ? undefined : (control.horizontalalign ? control.horizontalalign : "start"),
      verticalAlign: control.verticalalign === '' ? undefined : (control.verticalalign ? control.verticalalign : "start"),
      styles: {
        root: {
          width: control.width ? defaultPixels(control.width) : "100%",
          height: control.height !== undefined ? defaultPixels(control.height) : undefined,
          padding: control.padding ? defaultPixels(control.padding) : "10px",
          margin: control.margin !== undefined ? defaultPixels(control.margin) : undefined,
          backgroundColor: control.bgcolor ? getThemeColor(theme, control.bgcolor) : undefined,
        }
      },
    };

    const stackTokens: IStackTokens = {
      childrenGap: control.gap ? control.gap : 10
    }

    const authProviders = control.signin ? control.signin.split(",").map((s: string) => s.trim().toLowerCase()) : [];
    const signinGroups = isTrue(control.signingroups)

    const handleDismiss = () => {
      const payload: any = {
        i: "page",
        signin: ''
      }

      dispatch(changeProps([payload]));
      ws.updateControlProps([payload]);
      ws.pageEventFromWeb("page", 'dismissSignin', "");
    }

    let signinProps: ISigninProps = {
      signinOptions: {
        gitHubEnabled: authProviders.includes("github") || authProviders.includes("*"),
        gitHubGroupScope: signinGroups,
        azureEnabled: authProviders.includes("azure") || authProviders.includes("*"),
        azureGroupScope: signinGroups,
        googleEnabled: authProviders.includes("google") || authProviders.includes("*"),
        googleGroupScope: signinGroups
      },
      onDismiss: isTrue(control.signinallowdismiss) ? handleDismiss : undefined
    }

    return <>
      <Stack tokens={stackTokens} {...stackProps}>
        <ControlsList controls={childControls} parentDisabled={disabled} />
      </Stack>
      {authProviders.length > 0 &&
        <Signin {...signinProps} />
      }
    </>
  }

  const renderContent = () => {
    if (err === "signin_required") {
      return <Signin signinOptions={signinOptions} />
    }
    else if (err) {
      return <MessageBar messageBarType={MessageBarType.error} isMultiline={true}>{err}</MessageBar>
    } else if (!control) {
      return <Spinner label="Loading, please wait..." labelPosition="right" styles={{ root: { height: "100vh" } }} />
    } else {
      return renderPage()
    }
  }

  return <ThemeProvider theme={theme ?? defaultTheme} className={className}>
    {renderContent()}
  </ThemeProvider>
};