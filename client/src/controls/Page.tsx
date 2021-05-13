import React from 'react'
import { shallowEqual, useSelector, useDispatch } from 'react-redux'
import { ControlsList } from './ControlsList'
import useTitle from '../hooks/useTitle'
import { Stack, IStackProps, IStackTokens, createTheme, ThemeProvider, mergeStyles, PartialTheme } from '@fluentui/react';
import {
  BaseSlots,
  ThemeGenerator,
  themeRulesStandardCreator,
} from '@fluentui/react/lib/ThemeGenerator';
import { Login } from './Login'
import { isDark } from '@fluentui/react/lib/Color';
import { ILoginProps, IPageProps } from './Control.types'
import { WebSocketContext } from '../WebSocket';
import { changeProps } from '../slices/pageSlice'
import { defaultPixels, getWindowHash, isFalse, isTrue } from './Utils'

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

  function buildTheme() {
    // theme
    const themePrimaryColor = control.themeprimarycolor ? control.themeprimarycolor : '#8e16c9'
    const themeTextColor = control.themetextcolor ? control.themetextcolor : '#020203'
    const themeBackgroundColor = control.themebackgroundcolor ? control.themebackgroundcolor : '#ffffff'

    // theme
    let themeRules = themeRulesStandardCreator();
    function changeColor(baseSlot: BaseSlots, newColor: any) {
      const currentIsDark = isDark(themeRules[BaseSlots[BaseSlots.backgroundColor]].color!);
      ThemeGenerator.setSlot(themeRules[BaseSlots[baseSlot]], newColor, currentIsDark, true, true);
      if (currentIsDark !== isDark(themeRules[BaseSlots[BaseSlots.backgroundColor]].color!)) {
        // isInverted got swapped, so need to refresh slots with new shading rules
        ThemeGenerator.insureSlots(themeRules, currentIsDark);
      }
    }

    changeColor(BaseSlots.primaryColor, themePrimaryColor);
    changeColor(BaseSlots.backgroundColor, themeBackgroundColor);
    changeColor(BaseSlots.foregroundColor, themeTextColor);
    changeColor(BaseSlots.backgroundColor, themeBackgroundColor);

    const themeAsJson: {
      [key: string]: string;
    } = ThemeGenerator.getThemeAsJson(themeRules);

    setTheme(createTheme({
      ...{ palette: themeAsJson },
      isInverted: isDark(themeRules[BaseSlots[BaseSlots.backgroundColor]].color!),
    }));

    document.documentElement.style.background = themeBackgroundColor;
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

    buildTheme();

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

  const loginProviders = control.login ? control.login.split(",").map((s:string) => s.trim().toLowerCase()) : [];
  const loginGroups = isTrue(control.logingroups)

  const handleDismiss = () => {
    const payload: any = {
      i: "page",
      login: ''
    }

    dispatch(changeProps([payload]));
    ws.updateControlProps([payload]);
    ws.pageEventFromWeb("page", 'loginCancelled', "");
  }

  let loginProps: ILoginProps = {
    loginOptions: {
      gitHubEnabled: loginProviders.includes("github") || loginProviders.includes("*"),
      gitHubGroupScope: loginGroups,
      azureEnabled: loginProviders.includes("azure") || loginProviders.includes("*"),
      azureGroupScope: loginGroups,
      googleEnabled: loginProviders.includes("google") || loginProviders.includes("*"),
      googleGroupScope: loginGroups
    },
    onDismiss: isTrue(control.loginallowdismiss) ? handleDismiss : undefined
  }

  return <ThemeProvider theme={theme} className={className}>
      <Stack tokens={stackTokens} {...stackProps}>
        <ControlsList controls={childControls} parentDisabled={disabled} />
      </Stack>
      { loginProviders.length > 0 &&
        <Login {...loginProps} />
      }
    </ThemeProvider>
})