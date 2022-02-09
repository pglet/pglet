import { Theme, isDark, createTheme, IStyle } from '@fluentui/react';
import { BaseSlots, ThemeGenerator, themeRulesStandardCreator } from '@fluentui/react/lib/ThemeGenerator';

export const lightThemeColor = {
  primary: '#6529d6',
  text: '#020203',
  background: '#ffffff'
}

export const darkThemeColor = {
  primary: '#8f87ff',
  text: '#d0d4d9',
  background: '#24292e'
}

export function buildTheme(standardTheme: any, themePrimaryColor: any, themeTextColor: any, themeBackgroundColor: any): Theme {

  let primaryColor = lightThemeColor.primary;
  let textColor = lightThemeColor.text;
  let backgroundColor = lightThemeColor.background;

  if (standardTheme && standardTheme.toLowerCase() === 'dark') {
    primaryColor = darkThemeColor.primary;
    textColor = darkThemeColor.text;
    backgroundColor = darkThemeColor.background;
  }

  if (themePrimaryColor) {
    primaryColor = themePrimaryColor
  }
  if (themeTextColor) {
    textColor = themeTextColor
  }
  if (themeBackgroundColor) {
    backgroundColor = themeBackgroundColor
  }

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

  changeColor(BaseSlots.primaryColor, primaryColor);
  changeColor(BaseSlots.backgroundColor, backgroundColor);
  changeColor(BaseSlots.foregroundColor, textColor);
  changeColor(BaseSlots.backgroundColor, backgroundColor);

  const themeAsJson: {
    [key: string]: string;
  } = ThemeGenerator.getThemeAsJson(themeRules);

  let theme = createTheme({
    ...{ palette: themeAsJson },
    isInverted: isDark(themeRules[BaseSlots[BaseSlots.backgroundColor]].color!),
  });

  // 
  if (isDefaultDark(theme)) {
    customizeDarkTheme(theme);
  }

  return theme;
}

export function getThemeExtraStyles(theme: Theme | undefined): IStyle {
  if (theme === undefined) {
    return null;
  }
  if (isDefaultDark(theme)) {
    return {
      ".ms-MessageBar--error": {
        backgroundColor: '#5f2725'
      },
      ".ms-MessageBar--blocked": {
        backgroundColor: '#5f2725'
      },
      ".ms-MessageBar--success": {
        backgroundColor: '#2a3c1b'
      },
      ".ms-MessageBar--warning": {
        backgroundColor: '#5f4519'
      },
      ".ms-MessageBar--severeWarning": {
        backgroundColor: '#673612'
      }
    };
  }
  return {};
}

function isDefaultDark(theme: Theme): boolean {
  return theme.palette.themePrimary === darkThemeColor.primary &&
    theme.semanticColors.bodyText === darkThemeColor.text &&
    theme.semanticColors.bodyBackground === darkThemeColor.background;
}

function customizeDarkTheme(theme: Theme) {
  //theme.semanticColors.inputBackground = "#2c3136";
}