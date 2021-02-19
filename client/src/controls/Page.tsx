import React from 'react'
import { shallowEqual, useSelector, useDispatch } from 'react-redux'
import { ControlsList } from './ControlsList'
import useTitle from '../hooks/useTitle'
import { Stack, IStackProps, IStackTokens, createTheme, ThemeProvider, mergeStyles } from '@fluentui/react';
import {
  BaseSlots,
  ThemeGenerator,
  themeRulesStandardCreator,
} from '@fluentui/react/lib/ThemeGenerator';
import { isDark } from '@fluentui/react/lib/Color';
import { IPageProps } from './Control.types'
import { WebSocketContext } from '../WebSocket';
import { changeProps } from '../slices/pageSlice'
import { defaultPixels, getWindowHash } from './Utils'

export const Page = React.memo<IPageProps>(({ control, pageName }) => {

  const ws = React.useContext(WebSocketContext);
  const dispatch = useDispatch();

  // page title
  let title = `${pageName} - pglet`;
  if (control.title) {
    title = control.title
  }
  useTitle(title)

  // theme
  const themePrimaryColor = control.themeprimarycolor ? control.themeprimarycolor : '#8e16c9'
  const themeTextColor = control.themetextcolor ? control.themetextcolor : '#020203'
  const themeBackgroundColor = control.themebackgroundcolor ? control.themebackgroundcolor : '#ffffff'

  //console.log("themeBackgroundColor:", themeBackgroundColor);  

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

  const theme = createTheme({
    ...{ palette: themeAsJson },
    isInverted: isDark(themeRules[BaseSlots[BaseSlots.backgroundColor]].color!),
  });

  // strip out the unnecessary shade slots from the final output theme
  // const abridgedTheme: IThemeRules = {};
  // for (const ruleName in themeRules) {
  //   if (themeRules.hasOwnProperty(ruleName)) {
  //     if (
  //       ruleName.indexOf('ColorShade') === -1 &&
  //       ruleName !== 'primaryColor' &&
  //       ruleName !== 'backgroundColor' &&
  //       ruleName !== 'foregroundColor' &&
  //       ruleName.indexOf('body') === -1
  //     ) {
  //       abridgedTheme[ruleName] = themeRules[ruleName];
  //     }
  //   }
  // }

  // const jsonTheme = JSON.stringify(ThemeGenerator.getThemeAsJson(abridgedTheme), undefined, 2)

  function updateHash(hash: string) {

    const payload: any = {
      i: "page",
      hash: hash
    }

    dispatch(changeProps([payload]));
    ws.updateControlProps([payload]);
    ws.pageEventFromWeb("page", 'hash', hash);
  }

  React.useEffect(() => {

    const hash = getWindowHash();
    const pageHash = control.hash !== undefined ? control.hash : "";

    if (pageHash !== hash) {
      window.location.hash = pageHash ? "#" + pageHash : "";
    }

    // https://danburzo.github.io/react-recipes/recipes/use-effect.html
    // https://codedaily.io/tutorials/72/Creating-a-Reusable-Window-Event-Listener-Hook-with-useEffect-and-useCallback
    const handleWindowClose = (e: any) => {
      ws.pageEventFromWeb(control.i, 'close', control.data);
    }

    const handleHashChange = (e: any) => {
      updateHash(getWindowHash());
    }

    window.addEventListener("beforeunload", handleWindowClose);
    window.addEventListener("hashchange", handleHashChange);

    return () => {
      window.removeEventListener("beforeunload", handleWindowClose);
      window.removeEventListener("hashchange", handleHashChange);
    }
    // eslint-disable-next-line
  }, [control, ws]);

  const childControls = useSelector((state: any) => control.c.map((childId: string) => state.page.controls[childId]), shallowEqual);

  if (control.visible === "false") {
    return null;
  }

  let disabled = (control.disabled === "true")

  // stack props
  const stackProps: IStackProps = {
    verticalFill: control.verticalfill ? control.verticalfill === "true" : true,
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

  document.documentElement.style.background = themeBackgroundColor;

  const stackTokens: IStackTokens = {
    childrenGap: control.gap ? control.gap : 10
  }

  const className = mergeStyles({
    height: '100vh'
  });

  return <ThemeProvider theme={theme} className={className}>
    <Stack tokens={stackTokens} {...stackProps}>
      <ControlsList controls={childControls} parentDisabled={disabled} />
    </Stack>
  </ThemeProvider>
})