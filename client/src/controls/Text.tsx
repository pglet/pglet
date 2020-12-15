import React from 'react'
import { Text, ITextProps } from '@fluentui/react';
import { IControlProps } from './IControlProps'

export const MyText = React.memo<IControlProps>(({control}) => {

  //console.log(`render Text: ${control.i}`);

  // https://developer.microsoft.com/en-us/fluentui#/controls/web/references/ifontstyles#IFontStyles

  const textProps: ITextProps = {
    variant: control.size ? control.size : null,
    nowrap: control.nowrap !== undefined ? control.nowrap : undefined,
    block: control.block !== undefined ? control.block : undefined,
    styles: {
      root: {
        fontWeight: control.bold === 'true' ? 'bold' : undefined,
        fontStyle: control.italic === 'true' ? 'italic' : undefined,
        textAlign: control.align !== undefined ? control.align : undefined,
        width: control.width !== undefined ? control.width : undefined,
        height: control.height !== undefined ? control.height : undefined,
        padding: control.padding !== undefined ? control.padding : undefined,
        margin: control.margin !== undefined ? control.margin : undefined   
      }
    }
  }; 

  return <Text {...textProps}>{control.value}</Text>;
})