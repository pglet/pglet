import React from 'react'
import { Image, IImageProps, ImageFit } from '@fluentui/react';
import { IControlProps } from './Control.types'
import { defaultPixels } from './Utils'

export const MyImage = React.memo<IControlProps>(({control}) => {

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
    maximizeFrame: control.maximizeframe === "true"
  };

  return <Image {...imgProps} />;
})