import React from 'react'
import { useTheme, IPersonaProps, Persona, PersonaSize, PersonaPresence, PersonaInitialsColor } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { getThemeColor, isTrue } from './Utils'

export const MyPersona = React.memo<IControlProps>(({ control }) => {

  const theme = useTheme();

  let size: PersonaSize | undefined = undefined;
  switch (control.size ? `${control.size}` : "") {
    case '8': size = PersonaSize.size8; break;
    case '24': size = PersonaSize.size24; break;
    case '32': size = PersonaSize.size32; break;
    case '40': size = PersonaSize.size40; break;
    case '48': size = PersonaSize.size48; break;
    case '56': size = PersonaSize.size56; break;
    case '72': size = PersonaSize.size72; break;
    case '100': size = PersonaSize.size100; break;
    case '120': size = PersonaSize.size120; break;
  }

  let presence: PersonaPresence | undefined = undefined;
  switch (control.presence ? control.presence.toLowerCase() : '') {
    case 'none': presence = PersonaPresence.none; break;
    case 'offline': presence = PersonaPresence.offline; break;
    case 'online': presence = PersonaPresence.online; break;
    case 'away': presence = PersonaPresence.away; break;
    case 'blocked': presence = PersonaPresence.blocked; break;
    case 'busy': presence = PersonaPresence.busy; break;
    case 'dnd': presence = PersonaPresence.dnd; break;
  }

  let color: PersonaInitialsColor | undefined = undefined;
  switch (control.initialscolor ? control.initialscolor.toLowerCase() : '') {
    case 'blue': color = PersonaInitialsColor.blue; break;
    case 'burgundy': color = PersonaInitialsColor.burgundy; break;
    case 'coolgray': color = PersonaInitialsColor.coolGray; break;
    case 'cyan': color = PersonaInitialsColor.cyan; break;
    case 'darkblue': color = PersonaInitialsColor.darkBlue; break;
    case 'darkgreen': color = PersonaInitialsColor.darkGreen; break;
    case 'darkred': color = PersonaInitialsColor.darkRed; break;
    case 'gold': color = PersonaInitialsColor.gold; break;
    case 'green': color = PersonaInitialsColor.green; break;
    case 'lightblue': color = PersonaInitialsColor.lightBlue; break;
    case 'lightgreen': color = PersonaInitialsColor.lightGreen; break;
    case 'lightpink': color = PersonaInitialsColor.lightPink; break;
    case 'lightred': color = PersonaInitialsColor.lightRed; break;
    case 'magenta': color = PersonaInitialsColor.magenta; break;
    case 'orange': color = PersonaInitialsColor.orange; break;
    case 'pink': color = PersonaInitialsColor.pink; break;
    case 'purple': color = PersonaInitialsColor.purple; break;
    case 'rust': color = PersonaInitialsColor.rust; break;
    case 'teal': color = PersonaInitialsColor.teal; break;
    case 'transparent': color = PersonaInitialsColor.transparent; break;
    case 'violet': color = PersonaInitialsColor.violet; break;
    case 'warmgray': color = PersonaInitialsColor.warmGray; break;
  }

  const personaProps: IPersonaProps = {
    imageUrl: control.imageurl ? control.imageurl : undefined,
    imageAlt: control.imagealt ? control.imagealt : undefined,
    initialsColor: color,
    initialsTextColor: control.initialstextcolor ? getThemeColor(theme, control.initialstextcolor) : undefined,
    text: control.text ? control.text : undefined,
    secondaryText: control.secondarytext ? control.secondarytext : undefined,
    tertiaryText: control.tertiarytext ? control.tertiarytext : undefined,
    optionalText: control.optionaltext ? control.optionaltext : undefined,
    size: size,
    presence: presence,
    hidePersonaDetails: isTrue(control.hidedetails),
  };

  return <Persona {...personaProps} />;
})