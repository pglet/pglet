import React from 'react'
import { Text, ITextProps, IFontStyles, mergeStyles, useTheme } from '@fluentui/react';
import { IControlProps, defaultPixels } from './IControlProps'

// Markdown support
import ReactMarkdown from 'react-markdown'
import gfm from 'remark-gfm'

export const MyText = React.memo<IControlProps>(({ control }) => {

  //console.log(`render Text: ${control.i}`);

  // https://developer.microsoft.com/en-us/fluentui#/controls/web/references/ifontstyles#IFontStyles

  const preFont = 'SFMono-Regular, Menlo, Monaco, Consolas, "Liberation Mono", "Courier New", monospace';

  const textAlign = control.align !== undefined ? control.align : undefined;
  const verticalAlign = control.verticalalign !== undefined ? control.verticalalign : undefined;
  let display = undefined;
  let alignItems = undefined;
  let justifyContent = undefined;

  if (verticalAlign !== undefined) {
    // enable flex mode
    display = 'inline-flex';

    if (verticalAlign === 'top') {
      alignItems = "flex-start";
    } else if (verticalAlign === 'bottom') {
      alignItems = "flex-end";
    } else if (verticalAlign === 'center' || verticalAlign === 'middle') {
      alignItems = "center";
    }

    // adjust horizontal align
    if (textAlign === 'left') {
      justifyContent = "flex-start";
    } else if (textAlign === 'right') {
      justifyContent = "flex-end";
    } else if (textAlign === 'center' || textAlign === 'middle') {
      justifyContent = "center";
    }
  }

  let variant: keyof IFontStyles | undefined = undefined;
  switch (control.size ? control.size.toLowerCase() : '') {
    case 'tiny': variant = 'tiny'; break;
    case 'xsmall': variant = 'xSmall'; break;
    case 'small': variant = 'small'; break;
    case 'smallplus': variant = 'smallPlus'; break;
    case 'medium': variant = 'medium'; break;
    case 'mediumplus': variant = 'mediumPlus'; break;
    case 'large': variant = 'large'; break;
    case 'xlarge': variant = 'xLarge'; break;
    case 'xxlarge': variant = 'xxLarge'; break;
    case 'superlarge': variant = 'superLarge'; break;
    case 'mega': variant = 'mega'; break;
  }

  // https://github.com/microsoft/fluentui/blob/master/packages/merge-styles/README.md
  const theme = useTheme();
  const className = mergeStyles({
    selectors: {
      '& pre': {
        backgroundColor: theme.palette.neutralLighter,
        borderRadius: "2px"
      },
      '& a': {
        color: theme.palette.themePrimary,
      }
    }
  });

  const textProps: ITextProps = {
    variant: variant,
    className: className,
    nowrap: control.nowrap !== undefined ? control.nowrap : undefined,
    block: control.block !== undefined ? control.block : undefined,
    styles: {
      root: {
        display: display,
        alignItems: alignItems,
        justifyContent: justifyContent,
        textAlign: textAlign,
        color: control.color ? control.color : undefined,
        backgroundColor: control.bgcolor ? control.bgcolor : undefined,
        border: control.border ? control.border : undefined,
        borderRadius: control.borderradius ? defaultPixels(control.borderradius) : undefined,
        borderLeft: control.borderleft ? control.borderleft : undefined,
        borderRight: control.borderright ? control.borderright : undefined,
        borderTop: control.bordertop ? control.bordertop : undefined,
        borderBottom: control.borderbottom ? control.borderbottom : undefined,
        fontWeight: control.bold === 'true' ? 'bold' : undefined,
        fontStyle: control.italic === 'true' ? 'italic' : undefined,
        whiteSpace: control.pre === 'true' ? 'pre' : undefined,
        fontFamily: control.pre === 'true' ? preFont : undefined,
        width: control.width !== undefined ? defaultPixels(control.width) : undefined,
        height: control.height !== undefined ? defaultPixels(control.height) : undefined,
        padding: control.padding !== undefined ? defaultPixels(control.padding) : undefined,
        margin: control.margin !== undefined ? defaultPixels(control.margin) : undefined,
      }
    }
  };

  if (control.markdown === 'true') {
    return <Text className={className}><ReactMarkdown plugins={[gfm]} children={control.value} /></Text>;
  } else {
    return <Text {...textProps}>{control.pre === "true" ? <pre>{control.value}</pre> : control.value}</Text>;
  }
})