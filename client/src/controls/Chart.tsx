import React from 'react';
import { shallowEqual, useSelector } from 'react-redux'
import { VerticalBarChart, IVerticalBarChartProps } from '@fluentui/react-charting';
import { IControlProps, defaultPixels } from './IControlProps'

export const Chart = React.memo<IControlProps>(({control, parentDisabled}) => {

  const chartProps: IVerticalBarChartProps = {
    styles: {
      root: {
        width: control.width !== undefined ? defaultPixels(control.width) : undefined,
        height: control.height !== undefined ? defaultPixels(control.height) : undefined,
        padding: control.padding !== undefined ? defaultPixels(control.padding) : undefined,
        margin: control.margin !== undefined ? defaultPixels(control.margin) : undefined   
      }
    }
  };

  const data = useSelector<any, any>((state: any) => {
    return control.c.map((childId: any) => state.page.controls[childId])
    .filter((c: any) => c.t === 'data').map((data: any) => 
      ({
        ...data,
        points: data.c.map((childId: any) => state.page.controls[childId])
      })
    );
  }, shallowEqual);  

  //console.log(data);

  if (data.length > 0) {
    chartProps.data = data[0].points;
  }

  return (
      <VerticalBarChart {...chartProps} />
  );
})