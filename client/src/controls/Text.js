import React from 'react'
import { Text } from '@fluentui/react';

const MyText = React.memo(({ control }) => {

  //console.log(`render Text: ${control.i}`);

  // https://developer.microsoft.com/en-us/fluentui#/controls/web/references/ifontstyles#IFontStyles

  const textProps = {
    variant: control.size ? control.size : null
  };

  return <Text {...textProps}>{control.value}</Text>;
})

export default MyText