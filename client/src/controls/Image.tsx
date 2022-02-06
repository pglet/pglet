import React from 'react'
import { Image, IImageProps, ImageFit, mergeStyles, useTheme } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { defaultPixels, getThemeColor, isTrue } from './Utils'

export const MyImage = React.memo<IControlProps>(({ control }) => {

  const theme = useTheme();

  let fit: ImageFit | undefined = undefined;
  switch (control.fit ? control.fit.toLowerCase() : '') {
    case 'none': fit = ImageFit.none; break;
    case 'contain': fit = ImageFit.contain; break;
    case 'cover': fit = ImageFit.cover; break;
    case 'center': fit = ImageFit.center; break;
    case 'centercontain': fit = ImageFit.centerContain; break;
    case 'centercover': fit = ImageFit.centerCover; break;
  }

  const imgProps: IImageProps = {
    src: control.src ? control.src : undefined,
    alt: control.alt ? control.alt : undefined,
    title: control.title ? control.title : undefined,
    width: control.width !== undefined ? defaultPixels(control.width) : undefined,
    height: control.height !== undefined ? defaultPixels(control.height) : undefined,
    imageFit: fit,
    maximizeFrame: isTrue(control.maximizeframe),
    className: mergeStyles({
      borderWidth: control.borderwidth ? defaultPixels(control.borderwidth) : undefined,
      borderColor: control.bordercolor ? getThemeColor(theme, control.bordercolor) : undefined,
      borderStyle: control.borderstyle ? control.borderstyle : undefined,
      borderRadius: control.borderradius ? defaultPixels(control.borderradius) : undefined
    })
  };

  return <Image {...imgProps} />;
})