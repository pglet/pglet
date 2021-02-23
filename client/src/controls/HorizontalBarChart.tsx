import React from 'react';
import { shallowEqual, useSelector } from 'react-redux'
import { HorizontalBarChart, IHorizontalBarChartProps } from '@fluentui/react-charting';
import { IControlProps } from './Control.types'
import { getThemeColor, defaultPixels, parseNumber, isTrue } from './Utils'
import { useTheme, mergeStyles } from '@fluentui/react';

export const MyHorizontalBarChart = React.memo<IControlProps>(({control, parentDisabled}) => {

  const theme = useTheme();

  const containerRef = React.useRef<HTMLDivElement>(null);

  const [dimensions, setDimensions] = React.useState({ 
    height: containerRef.current?.offsetHeight,
    width: containerRef.current?.offsetWidth
  })

  React.useEffect(() => {

    let resizeTimeout:any = null;
    function handleResize() {
      clearTimeout(resizeTimeout);
      resizeTimeout = setTimeout(() => {
        setDimensions({
          height: containerRef.current?.offsetHeight,
          width: containerRef.current?.offsetWidth
        })
      }, 500)
    }

    window.addEventListener('resize', handleResize)
    return () => window.removeEventListener("resize", handleResize);
  }, []);

  const colors = [
    getThemeColor(theme, "themeDarker"),
    getThemeColor(theme, "themeDarkAlt"),
    getThemeColor(theme, "themeTertiary"),
    getThemeColor(theme, "themeLighter")]

  const data = useSelector<any, any>((state: any) => {
    let colorIdx = 0;
    return control.c.map((childId: any) => state.page.controls[childId])
    .filter((c: any) => c.t === 'data').map((data: any) => 
      ({
        ...data,
        points: data.c.map((childId: any) => {
          const p = state.page.controls[childId];
          return {
            chartTitle: p.legend,
            chartData: [
              {
                legend: p.legend,
                horizontalBarChartdata: {
                  x: parseNumber(p.x),
                  y: parseNumber(p.y),
                },
                color: p.color !== undefined ? getThemeColor(theme, p.color) : colors[colorIdx++ % colors.length],
                xAxisCalloutData: p.xtooltip,
                yAxisCalloutData: p.ytooltip,
                ratio: isTrue(p.ratio)
              }
            ]
          }
        })
      })
    );
  }, shallowEqual);

  const containerClass = mergeStyles({
    width: control.width !== undefined ? defaultPixels(control.width) : undefined,
    height: control.height !== undefined ? defaultPixels(control.height) : undefined,
    padding: control.padding !== undefined ? defaultPixels(control.padding) : undefined,
    margin: control.margin !== undefined ? defaultPixels(control.margin) : undefined       
  });

  const chartProps: IHorizontalBarChartProps = {
    hideTooltip: control.tooltips !== 'true',
    //barHeight: control.barheight !== undefined ? parseNumber(control.barheight) : undefined,
    //hideRatio: control.ratio !== 'true',
    chartDataMode: control.datamode !== undefined ? control.datamode : 'default',
    width: dimensions.width ? dimensions.width : containerRef.current?.offsetWidth,
  };

  if (data.length > 0) {
    chartProps.data = data[0].points;

    // const ratios = data[0].points.map((d:any) => d.chartData.map((p:any) => !p.ratio))
    //   .reduce((acc: any, items: any) => ([...acc, ...items]));
    //chartProps.hideRatio = [true, false];
  }

  return (
      <div ref={containerRef} className={containerClass}><HorizontalBarChart {...chartProps} /></div>
  );
})