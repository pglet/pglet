import React from 'react'
import { shallowEqual, useSelector, useDispatch } from 'react-redux'
import { ControlsList } from './ControlsList'
import useTitle from '../hooks/useTitle'
import { Stack, IStackProps, IStackTokens, ThemeProvider, mergeStyles, PartialTheme } from '@fluentui/react';
import { Signin } from './Signin'
import { ISigninProps, IPageProps } from './Control.types'
import { WebSocketContext } from '../WebSocket';
import { changeProps } from '../slices/pageSlice'
import { defaultPixels, getWindowHash, isFalse, isTrue } from './Utils'
import { lightThemeColor, buildTheme } from './Theming'

export const Page = React.memo<IPageProps>(({ control, pageName }) => {

  const ws = React.useContext(WebSocketContext);
  const dispatch = useDispatch();
  const [theme, setTheme] = React.useState<PartialTheme | undefined>();

  // page title
  let title = `${pageName} - pglet`;
  if (control.title) {
    title = control.title
  }
  useTitle(title)

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

    // theme
    const themePrimaryColor = control.themeprimarycolor ? control.themeprimarycolor : lightThemeColor.primary
    const themeTextColor = control.themetextcolor ? control.themetextcolor : lightThemeColor.text
    const themeBackgroundColor = control.themebackgroundcolor ? control.themebackgroundcolor : lightThemeColor.background

    var theme = buildTheme(themePrimaryColor, themeTextColor, themeBackgroundColor)
    setTheme(theme);
    document.documentElement.style.background = themeBackgroundColor;

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

    return () => {
      window.removeEventListener("hashchange", handleHashChange);
    }
    // eslint-disable-next-line
  }, [control, ws]);

  const childControls = useSelector((state: any) => control.c.map((childId: string) => state.page.controls[childId]), shallowEqual);

  if (isFalse(control.visible)) {
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
        margin: control.margin !== undefined ? defaultPixels(control.margin) : undefined
      }
    },
  };

  const stackTokens: IStackTokens = {
    childrenGap: control.gap ? control.gap : 10
  }

  const className = mergeStyles({
    height: '100vh'
  });

  const authProviders = control.signin ? control.signin.split(",").map((s:string) => s.trim().toLowerCase()) : [];
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

  return <ThemeProvider theme={theme} className={className}>
      <Stack tokens={stackTokens} {...stackProps}>
        <ControlsList controls={childControls} parentDisabled={disabled} />
      </Stack>
      { authProviders.length > 0 &&
        <Signin {...signinProps} />
      }
    </ThemeProvider>
})