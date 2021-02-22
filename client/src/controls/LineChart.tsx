import React from 'react';
import { shallowEqual, useSelector } from 'react-redux'
import { LineChart, ILineChartProps } from '@fluentui/react-charting';
import { IControlProps } from './Control.types'
import { parseNumber, getThemeColor, defaultPixels, isTrue } from './Utils'
import { useTheme } from '@fluentui/react';

export const MyLineChart = React.memo<IControlProps>(({control, parentDisabled}) => {

  const theme = useTheme();
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

  const xtype = control.xtype ? control.xtype : "number"

  const colors = [
    getThemeColor(theme, "themeDark"),
    getThemeColor(theme, "themePrimary"),
    getThemeColor(theme, "themeTertiary"),
    getThemeColor(theme, "themeLighter")]  

  const data = useSelector<any, any>((state: any) => {
    let colorIdx = 0;
    return control.c.map((childId: any) => state.page.controls[childId])
    .filter((c: any) => c.t === 'data').map((data: any) => 
      ({
        ...data,
        color: data.color !== undefined ? data.color : colors[colorIdx++ % colors.length],
        data: data.c.map((childId: any) => {
          const p = state.page.controls[childId];
          const y = parseNumber(p.y)
          return {
            x: xtype === "date" ? new Date(p.x) : parseNumber(p.x),
            y: y,
            tick: isTrue(p.tick),
            legend: p.legend,
            xAxisCalloutData: p.xtooltip ? p.xtooltip : p.x,
            yAxisCalloutData: p.ytooltip ? p.ytooltip : control.yformat !== undefined ? control.yformat.replace('{y}', y) : y
          }
        })
      })
    );
  }, shallowEqual);

  const ticks = data.map((d:any) => d.data.filter((p:any) => p.tick).map((p:any) => p.x))
    .reduce((acc: any, items: any) => ([...acc, ...items]));

  const chartProps: ILineChartProps = {
    data: {
      lineChartData: data
    },
    hideLegend: control.legend !== 'true',
    hideTooltip: control.tooltips !== 'true',
    strokeWidth: control.strokewidth !== undefined ? parseInt(control.strokewidth) : 2,
    yMinValue: control.ymin !== undefined ? parseFloat(control.ymin) : undefined,
    yMaxValue: control.ymax !== undefined ? parseFloat(control.ymax) : undefined,
    yAxisTickCount: control.yticks !== undefined ? parseInt(control.yticks) : 1,
    yAxisTickFormat: control.yformat !== undefined ? (y:any) => control.yformat.replace('{y}', y) : undefined,
    tickValues: ticks.length > 0 ? ticks : undefined,
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

  if (data.length > 0) {
    const ys = data.map((d:any) => d.data.map((p:any) => p.y))
      .reduce((acc: any, items: any) => ([...acc, ...items]));
    if (chartProps.yMinValue === undefined) {
      chartProps.yMinValue = Math.min(...ys)
    }
    if (chartProps.yMaxValue === undefined) {
      chartProps.yMaxValue = Math.max(...ys)
    }    
  }

  return (
      <LineChart {...chartProps} />
  );
})