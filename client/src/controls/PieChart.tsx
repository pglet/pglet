import React from 'react';
import { shallowEqual, useSelector } from 'react-redux'
import { DonutChart, IDonutChartProps } from '@fluentui/react-charting';
import { IControlProps } from './Control.types'
import { getThemeColor, defaultPixels } from './Utils'
import { useTheme, mergeStyles } from '@fluentui/react';

export const MyPieChart = React.memo<IControlProps>(({control, parentDisabled}) => {

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

  const containerClass = mergeStyles({
    width: control.width !== undefined ? defaultPixels(control.width) : undefined,
    height: control.height !== undefined ? defaultPixels(control.height) : undefined,
    padding: control.padding !== undefined ? defaultPixels(control.padding) : undefined,
    margin: control.margin !== undefined ? defaultPixels(control.margin) : undefined       
  });

  const chartProps: IDonutChartProps = {
    valueInsideDonut: control.innervalue !== undefined ? control.innervalue : undefined,
    innerRadius: control.innerradius !== undefined ? parseInt(control.innerradius) : undefined,
    hideLegend: control.legend !== 'true',
    hideTooltip: control.tooltips !== 'true',
    width: dimensions.width ? dimensions.width : containerRef.current?.offsetWidth,
    height: dimensions.height ? dimensions.height : containerRef.current?.offsetHeight,
  };

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
            legend: p.legend,
            data: p.value,
            color: p.color !== undefined ? p.color : colors[colorIdx++ % colors.length],
            xAxisCalloutData: p.tooltip,
          }
        })
      })
    );
  }, shallowEqual);

  if (data.length > 0) {
    chartProps.data = {
      chartData: data[0].points
    };
  }

  return (
      <div ref={containerRef} className={containerClass}><DonutChart {...chartProps} /></div>
  );
})