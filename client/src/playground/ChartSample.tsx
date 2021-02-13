import React, { useState, useEffect, useRef } from 'react';
import { VerticalBarChart, IVerticalBarChartProps, IDataPoint } from '@fluentui/react-charting';
import { DefaultPalette } from '@fluentui/react/lib/Styling';
import { DefaultButton } from '@fluentui/react/lib/Button';

export const ChartSample: React.FunctionComponent = () => {

  const [dynamicData, setDynamicData] = useState<IDataPoint[]>();
  const [colors, setColors] = useState<string[]>();
  const [x, setX] = useState<number>();

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

  const xRef = useRef(x);
  xRef.current = x;

  useEffect(() => {
    let arr = [
      { x: '23', y: 10 },
      { x: '24', y: 36 },
      { x: '25', y: 20 },
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

    setColors(_colors[0]);

    function handleResize() {
      setDimensions({
        height: window.innerHeight,
        width: window.innerWidth
      })

    }

    window.addEventListener('resize', handleResize)

  }, []);

  return (
    <div style={{ width: '100%', height: '400px' }}>
      <VerticalBarChart
        data={dynamicData}
        colors={colors}
        chartLabel={'Chart with Dynamic Data'}
        hideLegend={true}
        hideTooltip={true}
        yMaxValue={50}
        yAxisTickCount={5}
        height={dimensions.height}
        width={dimensions.width}
      />
      <DefaultButton text="Change data" onClick={_changeData} />
      <DefaultButton text="Change colors" onClick={_changeColors} />
    </div>
  );
}