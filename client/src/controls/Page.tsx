import React, { useEffect, useContext } from 'react'
import { shallowEqual, useSelector } from 'react-redux'
import { ControlsList } from './ControlsList'
import useTitle from '../hooks/useTitle'
import { Stack, IStackProps, IStackTokens, createTheme, ThemeProvider } from '@fluentui/react';
import {
  BaseSlots,
  ThemeGenerator,
  themeRulesStandardCreator,
} from '@fluentui/react/lib/ThemeGenerator';
import { isDark } from '@fluentui/react/lib/Color';
import { IPageProps } from './Control.types'
import { WebSocketContext } from '../WebSocket';
import { defaultPixels } from './Utils'

export const Page = React.memo<IPageProps>(({ control, pageName }) => {

  const ws = useContext(WebSocketContext);

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

  useEffect(() => {
    // https://danburzo.github.io/react-recipes/recipes/use-effect.html
    // https://codedaily.io/tutorials/72/Creating-a-Reusable-Window-Event-Listener-Hook-with-useEffect-and-useCallback
    const handleWindowClose = (e: any) => {
      ws.pageEventFromWeb(control.i, 'close', control.data);
    }
    window.addEventListener("beforeunload", handleWindowClose);
    return () => window.removeEventListener("beforeunload", handleWindowClose);
  }, [control, ws]);

  const childControls = useSelector((state: any) => control.c.map((childId: string) => state.page.controls[childId]), shallowEqual);

  if (control.visible === "false") {
    return null;
  }

  let disabled = (control.disabled === "true")

  // stack props
  const stackProps: IStackProps = {
    verticalFill: control.verticalfill ? control.verticalfill : true,
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

  document.documentElement.style.background = themeBackgroundColor

  return <ThemeProvider theme={theme}>
    <Stack tokens={stackTokens} {...stackProps}>
      <ControlsList controls={childControls} parentDisabled={disabled} />
    </Stack>
  </ThemeProvider>
})