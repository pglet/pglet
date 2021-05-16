import { Theme, isDark, createTheme } from '@fluentui/react';
import { BaseSlots, ThemeGenerator, themeRulesStandardCreator } from '@fluentui/react/lib/ThemeGenerator';

export const lightThemeColor = {
    primary: '#8e16c9',
    text: '#020203',
    background: '#ffffff'
}

export function buildTheme(themePrimaryColor:any, themeTextColor:any, themeBackgroundColor:any) : Theme {
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

    return createTheme({
      ...{ palette: themeAsJson },
      isInverted: isDark(themeRules[BaseSlots[BaseSlots.backgroundColor]].color!),
    });
}