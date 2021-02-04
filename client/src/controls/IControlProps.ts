export interface IControlProps {
    control: any;
    parentDisabled: boolean;
}

export function defaultPixels(size:any) {
    if (!size) {
        return size
    }

    const m = size.toString().match(/^\d*(\.\d+)?$/)
    if (m) {
        // just number
        return `${size}px`;
    }
    return size;
}