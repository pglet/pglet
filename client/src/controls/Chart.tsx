import React from 'react';
import { shallowEqual, useSelector } from 'react-redux'
import { VerticalBarChart, IVerticalBarChartProps } from '@fluentui/react-charting';
import * as d3 from 'd3-format';
import { IControlProps, defaultPixels } from './IControlProps'

export const Chart = React.memo<IControlProps>(({control, parentDisabled}) => {

  const [dimensions, setDimensions] = React.useState({ 
    height: window.innerHeight,
    width: window.innerWidth
  })

  React.useEffect(() => {

    let resizeTimeout:any = null;
    function handleResize() {
      clearTimeout(resizeTimeout);
      resizeTimeout = setTimeout(() => {
        setDimensions({
          height: window.innerHeight,
          width: window.innerWidth
        })
      }, 500)
    }

    window.addEventListener('resize', handleResize)
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  const chartProps: IVerticalBarChartProps = {
    hideLegend: control.legend !== 'true',
    hideTooltip: control.tooltips !== 'true',
    barWidth: control.barwidth !== undefined ? parseInt(control.barwidth) : undefined,
    colors: control.colors !== undefined ? control.colors.split(/[ ,]+/g) : undefined,
    yMinValue: control.ymin !== undefined ? parseFloat(control.ymin) : undefined,
    yMaxValue: control.ymax !== undefined ? parseFloat(control.ymax) : undefined,
    yAxisTickCount: control.yticks !== undefined ? parseInt(control.yticks) : 1,
    yAxisTickFormat: control.yformat !== undefined ? d3.format(control.yformat) : undefined,
    height: dimensions.height,
    width: dimensions.width,
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
        points: data.c.map((childId: any) => {
          const p = state.page.controls[childId];
          return {
            x: p.x,
            y: parseFloat(p.y),
            legend: p.legend,
            color: p.color,
            xAxisCalloutData: p.xtooltip,
            yAxisCalloutData: p.ytooltip
          }
        })
      })
    );
  }, shallowEqual);

  //console.log(data);

  if (data.length > 0) {
    chartProps.data = data[0].points;
    const yvals = chartProps.data!.map(p => p.y);
    if (chartProps.yMaxValue === undefined) {
      chartProps.yMaxValue = Math.max(...yvals)
    }
  }

  return (
      <VerticalBarChart {...chartProps} />
  );
})