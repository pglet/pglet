import React from 'react';
import { VerticalBarChart, IVerticalBarChartProps, IVerticalBarChartDataPoint } from '@fluentui/react-charting';
import { DefaultPalette/*, SharedColors, NeutralColors*/ } from '@fluentui/react';
import { DefaultButton } from '@fluentui/react/lib/Button';

export const ChartSample: React.FunctionComponent = () => {

  const [dynamicData, setDynamicData] = React.useState<IVerticalBarChartDataPoint[]>();
  const [colors, setColors] = React.useState<string[]>();
  const [x, setX] = React.useState<number>();

  const [dimensions, setDimensions] = React.useState({ 
    height: window.innerHeight,
    width: window.innerWidth
  })

  const _colors = [
    [DefaultPalette.purpleLight, DefaultPalette.blueMid],
    [DefaultPalette.orangeLighter, DefaultPalette.orangeLight, DefaultPalette.orange],
    [DefaultPalette.greenLight, DefaultPalette.green, DefaultPalette.greenDark],
    [DefaultPalette.magentaLight, DefaultPalette.magenta, DefaultPalette.magentaDark],
  ];
  let _colorIndex = 0;

  const _changeColors = () => {
    _colorIndex = (_colorIndex + 1) % _colors.length;
    setColors(_colors[_colorIndex]);
  }

  const _changeData = () => {
    let arr = dynamicData!;
    arr.shift();
    arr.push({ x: x!.toString(), y: _randomY() })
    setDynamicData(arr);
    setX(x! + 1);
    // this.setState({
    //   x: this.state.x + 1,
    //   dynamicData: arr
    // });
    //setTimeout(() => _changeData(), 2000);
  }

  const _randomY = () => {
    return Math.random() * 45 + 5;
  }

  const xRef = React.useRef(x);
  xRef.current = x;

  let scolor = 'neutralQuaternaryAlt';
  const color = Object.getOwnPropertyNames(DefaultPalette).filter(p => p.toLowerCase() === scolor.toLowerCase())
  if (color.length > 0) {
    scolor = (DefaultPalette as any)[color[0]]
  }

  React.useEffect(() => {
    let arr = [
      { x: '23', y: 10, legend: 'Bar 1', color: 'salmon', },
      { x: '24', y: 36, legend: 'Bar 2', color: scolor, },
      { x: '25', y: 20, xAxisCalloutData: 'X value', yAxisCalloutData: 'Y value' },
      { x: '26', y: 46 },
      { x: '27', y: 13 },
      { x: '28', y: 43 },
      { x: '29', y: 30 },
      { x: '30', y: 45 },
      { x: '31', y: 50 },
      { x: '32', y: 43 },
      { x: '33', y: 19 },
    ];

    setDynamicData(arr);
    setX(34);

    // https://github.com/facebook/react/issues/14010
    // setInterval(() => {
    //   console.log('addPoint')
    //   arr.shift();
    //   arr.push({ x: xRef.current!.toString(), y: _randomY() })
    //   setDynamicData(arr);
    //   setX(x => x! + 1);
    // }, 2000);

    //setColors(_colors[0]);

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
  }, [scolor]);

  let props: IVerticalBarChartProps = {
    data: dynamicData,
    colors: colors,
    chartLabel: 'Chart with Dynamic Data',
    hideLegend: true,
    hideTooltip: true,
    yMaxValue: 50,
    yAxisTickCount: 1,
    height: dimensions.height,
    width: dimensions.width
  }

  return (
    <div style={{ width: '100%', height: '400px' }}>
      <VerticalBarChart {...props} />
      <DefaultButton text="Change data" onClick={_changeData} />
      <DefaultButton text="Change colors" onClick={_changeColors} />
    </div>
  );
}