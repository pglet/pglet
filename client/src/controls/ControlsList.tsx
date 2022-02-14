import React from 'react'
import { MessageBar, MessageBarType } from '@fluentui/react'
import { IControlsListProps } from './Control.types'
import { MyStack } from './Stack'
import { Textbox } from './Textbox'
import { Searchbox } from './Searchbox'
import { MySpinButton } from './SpinButton'
import { MyDropdown } from './Dropdown'
import { MyChoiceGroup } from './ChoiceGroup'
import { MyCheckbox } from './Checkbox'
import { MyToggle } from './Toggle'
import { Progress } from './Progress'
import { MySpinner } from './Spinner'
import { MySlider } from './Slider'
import { Button } from './Button'
import { MyText } from './Text'
import { Html } from './Html'
import { MyImage } from './Image'
import { MyLink } from './Link'
import { MyDatePicker } from './DatePicker'
import { Tabs } from './Tabs'
import { Toolbar } from './Toolbar'
import { MyNav } from './Nav'
import { Grid } from './Grid'
import { Icon } from './Icon'
import { Message } from './Message'
import { MyDialog } from './Dialog'
import { MyPanel } from './Panel'
import { MyCallout } from './Callout'
import { IFrame } from './IFrame'
import { MyVerticalBarChart } from './VerticalBarChart'
import { MyHorizontalBarChart } from './HorizontalBarChart'
import { MyPieChart } from './PieChart'
import { MyLineChart } from './LineChart'
import { MyPersona } from './Persona'
import { isFalse } from './Utils'
import { MyComboBox } from './ComboBox'
import { MySplit } from './Split'

export const ControlsList: React.FunctionComponent<IControlsListProps> = ({ controls, parentDisabled }) => {

    const controlTypes: any = {
        'textbox': Textbox,
        'searchbox': Searchbox,
        'icon': Icon,
        'checkbox': MyCheckbox,
        'toggle': MyToggle,
        'dropdown': MyDropdown,
        'combobox': MyComboBox,
        'choicegroup': MyChoiceGroup,
        'progress': Progress,
        'spinner': MySpinner,
        'slider': MySlider,
        'text': MyText,
        'html': Html,
        'spinbutton': MySpinButton,
        'link': MyLink,
        'image': MyImage,
        'button': Button,
        'datepicker': MyDatePicker,
        'stack': MyStack,
        'tabs': Tabs,
        'toolbar': Toolbar,
        'nav': MyNav,
        'grid': Grid,
        'message': Message,
        'dialog': MyDialog,
        'panel': MyPanel,
        'callout': MyCallout,
        'iframe': IFrame,
        'verticalbarchart': MyVerticalBarChart,
        'barchart': MyHorizontalBarChart,
        'piechart': MyPieChart,
        'linechart': MyLineChart,
        'persona': MyPersona,
        'split': MySplit
    }

    const renderChild = (control: any) => {
        if (isFalse(control.visible)) {
            return null;
        }
        const ControlType = controlTypes[control.t];
        if (!ControlType) {
            const props = Object.getOwnPropertyNames(control)
                .filter(p => p.length > 1)
                .map(p => `${p}="${control[p]}"`).join(' ');
            return <MessageBar key={control.i} messageBarType={MessageBarType.error} messageBarIconProps={{ iconName: 'WebComponents' }}
                isMultiline><b>Unknown control:</b> {`${control.t} ${props}`}</MessageBar>
        }
        return <ControlType key={control.f ? control.f : control.i} control={control} parentDisabled={parentDisabled} />
    }

    return controls.map(renderChild);
}